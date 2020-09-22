package node

import (
	"github.com/zhengyangfeng00/woodenox/iterator"
	"github.com/zhengyangfeng00/woodenox/messagebus"
	"github.com/zhengyangfeng00/woodenox/stream"
)

var _ ProcessorUtils = (*Node)(nil)

type Node struct {
	mb        messagebus.MessageBus
	processor Processor
	done      chan struct{}
	iters     []streamIter
	outputs   map[string]messagebus.Output
}

func New(mb messagebus.MessageBus) *Node {
	return &Node{
		mb:      mb,
		outputs: make(map[string]messagebus.Output),
		done:    make(chan struct{}),
	}
}

type streamIter struct {
	stream    string
	iter      iterator.Iterator
	unsubFunc stream.UnsubFunc
}

type streamItem struct {
	stream string
	item   iterator.Item
}

func (n *Node) SetProcessor(p Processor) {
	n.processor = p
}

func (n *Node) Subscribe(s string) {
	iter, unsubFunc := n.mb.Subscribe(s, stream.WithBlock())
	n.iters = append(n.iters, streamIter{
		stream:    s,
		iter:      iter,
		unsubFunc: unsubFunc,
	})
}

func (n *Node) Run() {
	items := make(chan streamItem)
	for _, tp := range n.iters {
		go func(si streamIter) {
			// N.B: we expect the iterator to block when there is no new items.
			for tp.iter.Next() {
				select {
				case items <- streamItem{stream: si.stream, item: tp.iter.Curr()}:
				case <-n.done:
					return
				}
			}
		}(tp)
	}
	for {
		select {
		case <-n.done:
			return
		case si := <-items:
			n.processor.ProcessItem(si.stream, si.item)
		}
	}
}

func (n *Node) Stop() {
	close(n.done)
	for _, tp := range n.iters {
		tp.unsubFunc()
	}
}

func (n *Node) ProduceItem(stream string, item iterator.Item) {
	output, ok := n.outputs[stream]
	if !ok {
		output = n.mb.ProduceTo(stream)
		n.outputs[stream] = output
	}
	output.Write(item)
}
