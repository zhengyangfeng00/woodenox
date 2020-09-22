package messagebus

import (
	"sync"

	"github.com/zhengyangfeng00/woodenox/iterator"
	"github.com/zhengyangfeng00/woodenox/stream"
)

type MessageBus interface {
	Subscribe(
		streamName string,
		opts ...stream.SubscribeOpt,
	) (iterator.Iterator, stream.UnsubFunc)

	ProduceTo(stream string, opts ...stream.ProduceOpt) Output
}

type messageBus struct {
	streamLock sync.Mutex
	streams    map[string]stream.Stream

	newStreamFn func(name string, opts ...stream.Option) stream.Stream
}

func New() MessageBus {
	mb := &messageBus{
		streams: make(map[string]stream.Stream),
	}
	mb.newStreamFn = stream.New
	return mb
}

func (m *messageBus) Subscribe(
	s string,
	opts ...stream.SubscribeOpt,
) (iterator.Iterator, stream.UnsubFunc) {
	return m.getOrCreateStream(s).NewSubscriber(opts...)
}

func (m *messageBus) ProduceTo(stream string, opts ...stream.ProduceOpt) Output {
	return newOutput(m.getOrCreateStream(stream))
}

func (m *messageBus) Stop() {
	m.streamLock.Lock()
	defer m.streamLock.Unlock()

	for _, s := range m.streams {
		s.Stop()
	}
}

func (m *messageBus) getOrCreateStream(
	name string,
	opts ...stream.Option,
) stream.Stream {
	m.streamLock.Lock()
	defer m.streamLock.Unlock()

	s, ok := m.streams[name]
	if !ok {
		s = m.newStreamFn(name, opts...)
		go s.Run()
		m.streams[name] = s
	}
	return s
}
