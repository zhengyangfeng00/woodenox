package stream

import (
	"github.com/zhengyangfeng00/woodenox/iterator"
)

var _ iterator.Iterator = (*streamIter)(nil)

// stream defines the interface that the streamIter depends on. This
// is an internal interface.
type stream interface {
	Accept(item iterator.Item) error
	next(it *streamIter) (iterator.Item, error)
}

type streamIter struct {
	block    bool
	stream   stream
	curr     iterator.Item
	err      error
	done     chan struct{}
	notifyCh chan struct{}
}

func newStreamIter(s stream, block bool) *streamIter {
	return &streamIter{
		block:    block,
		stream:   s,
		done:     make(chan struct{}),
		notifyCh: make(chan struct{}),
	}
}

// Curr returns the current item pointed by the iterator. Curr should
// be called only after Next returns true, otherwise the behavior is
// undefined.
func (it *streamIter) Curr() iterator.Item {
	return it.curr
}

func (it *streamIter) Next() bool {
	if it.next() {
		return true
	}
	if !it.block {
		return false
	}
	for {
		select {
		case <-it.notifyCh:
			if it.next() {
				return true
			}
		case <-it.done:
			return false
		}
	}
}

func (it *streamIter) next() bool {
	item, err := it.stream.next(it)
	if err != nil {
		return false
	}
	it.curr = item
	it.err = nil
	return true
}

func (it *streamIter) Err() error {
	return it.err
}

// notify should be called when an new item is added to the
// stream. This method never blocks. The caller is expected to make
// superfluous calls to make sure the update is not missed.
func (it *streamIter) notify() {
	select {
	case it.notifyCh <- struct{}{}:
	default:
	}
}

func (it *streamIter) close() {
	close(it.done)
}
