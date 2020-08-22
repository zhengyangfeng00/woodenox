package messagebus

import (
	"github.com/zhengyangfeng00/woodenox/iterator"
	"github.com/zhengyangfeng00/woodenox/stream"
)

// Output provides an interface for writing to a stream. It handles
// retry, buffering, backpressure.
type Output interface {
	Write(item iterator.Item)
}

type outputImpl struct {
	s stream.Stream
}

func newOutput(s stream.Stream) Output {
	return &outputImpl{
		s: s,
	}
}

func (o *outputImpl) Write(item iterator.Item) {
	// TODO: retry, buffer, backpressure etc to be added.
	o.s.Accept(item)
}
