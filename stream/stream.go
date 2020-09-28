package stream

import (
	"errors"
	"sync"
	"time"

	"github.com/zhengyangfeng00/woodenox/iterator"
	"go.uber.org/zap"
)

var (
	errNoLongerExist = errors.New("item no longer exist")
	errNotAvailable  = errors.New("item not available")
	errNotSubscribed = errors.New("not subscribed")
	errStreamFull    = errors.New("stream is full")
	errClosed        = errors.New("stream closed")
)

type UnsubFunc func()

// Stream defines the public interface for a stream.
type Stream interface {
	NewSubscriber(opts ...SubscribeOpt) (iterator.Iterator, UnsubFunc)
	Accept(item iterator.Item) error
	Run()
	Stop()
}

type subscriberWithOffset struct {
	it     *streamIter
	offset int64
}

type streamImpl struct {
	name           string
	capacity       int
	purgeInterval  time.Duration
	notifyInterval time.Duration
	logger         *zap.Logger

	// A coarse-grained lock for the whole structure.
	mu sync.Mutex

	// offset is a monotonically increasing int to represent the item
	// order, while the underlying buffer is continously being cleaned
	// to make the memory usage bounded.
	startOffset int64
	buf         []iterator.Item
	subscribers map[*streamIter]*subscriberWithOffset
	purgeTicker *time.Ticker
	// notifyTicker is used to make spurious notify call to
	// subscribers in case they miss the notify message when a new
	// item is added.
	notifyTicker *time.Ticker
	done         chan struct{}
}

func New(name string, opts ...Option) Stream {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}
	return &streamImpl{
		name:         name,
		capacity:     options.capacity,
		subscribers:  make(map[*streamIter]*subscriberWithOffset),
		done:         make(chan struct{}),
		purgeTicker:  time.NewTicker(options.purgeInterval),
		notifyTicker: time.NewTicker(options.notifyInterval),
		logger:       options.logger,
	}
}

func (s *streamImpl) Run() {
	for {
		select {
		case <-s.purgeTicker.C:
			s.purge()
		case <-s.notifyTicker.C:
			s.mu.Lock()
			s.notifyWithLockHeld()
			s.mu.Unlock()
		case <-s.done:
			return
		}
	}
}

func (s *streamImpl) Stop() {
	close(s.done)
}

func (s *streamImpl) NewSubscriber(
	opts ...SubscribeOpt) (iterator.Iterator, UnsubFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()

	options := newSubOpts()
	for _, opt := range opts {
		opt(options)
	}

	it := newStreamIter(s, options.block)
	offset := s.idxToOffset(int64(len(s.buf) - 1))
	s.subscribers[it] = &subscriberWithOffset{
		it:     it,
		offset: offset,
	}
	return it, it.close
}

func (s *streamImpl) removeIter(it *streamIter) {
	// N.B: closing iterator doesn't require synchronization.
	it.close()

	s.mu.Lock()
	delete(s.subscribers, it)
	s.mu.Unlock()
}

func (s *streamImpl) next(it *streamIter) (iterator.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	itOffset, ok := s.subscribers[it]
	if !ok {
		return iterator.Item{}, errNotSubscribed
	}

	if itOffset.offset >= s.startOffset+int64(len(s.buf))-1 {
		return iterator.Item{}, errNotAvailable
	}

	itOffset.offset = itOffset.offset + 1
	return s.buf[s.offsetToIdx(itOffset.offset)], nil
}

func (s *streamImpl) Accept(item iterator.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.buf) >= s.capacity {
		return errStreamFull
	}
	s.buf = append(s.buf, item)

	s.logger.Debug(
		"notifying subscribers",
		zap.String("stream", s.name),
		zap.Int("numOfSubscribers", len(s.subscribers)),
	)
	s.notifyWithLockHeld()
	return nil
}

func (s *streamImpl) notifyWithLockHeld() {
	for _, sub := range s.subscribers {
		sub.it.notify()
	}
}

// purge removes old items from buffer. It is called inside the
// eventloop to keep it concurrent safe.
func (s *streamImpl) purge() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.buf) == 0 {
		return
	}

	if len(s.subscribers) == 0 {
		s.buf = nil
		s.startOffset = 0
		return
	}

	var lowestOffset = s.idxToOffset(int64(len(s.buf)) - 1)
	for _, sub := range s.subscribers {
		o := sub.offset
		if lowestOffset > o {
			lowestOffset = o
		}
	}

	// Some subscribers haven't started consuming.
	if lowestOffset == -1 {
		return
	}

	// delete items before lowestOffset
	idx := s.offsetToIdx(lowestOffset)
	s.buf = s.buf[idx:]
	s.startOffset = lowestOffset
}

func (s *streamImpl) idxToOffset(idx int64) int64 {
	return idx + s.startOffset
}

func (s *streamImpl) offsetToIdx(offset int64) int64 {
	return offset - s.startOffset
}
