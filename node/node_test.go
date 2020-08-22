package node

import (
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	iterator "github.com/zhengyangfeng00/woodenox/iterator"
	"github.com/zhengyangfeng00/woodenox/messagebus"
)

func TestNodeSubscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mb := messagebus.NewMockMessageBus(ctrl)
	it := iterator.NewMockIterator(ctrl)
	it.EXPECT().Next().AnyTimes().Return(true)
	it.EXPECT().Curr().AnyTimes().Return(iterator.Item{})
	unsubCalled := false
	unsub := func() {
		unsubCalled = true
	}
	mb.EXPECT().Subscribe(gomock.Any(), gomock.Any()).Return(it, unsub)
	node := New(mb)

	processedItem := make(chan struct{})
	p := NewMockProcessor(ctrl)
	p.EXPECT().ProcessItem(gomock.Any(), gomock.Any()).Do(
		func(stream string, item iterator.Item) {
			select {
			case <-processedItem:
				return
			default:
				close(processedItem)
			}
		},
	).AnyTimes()
	node.SetProcessor(p)
	node.Subscribe("test")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		node.Run()
		wg.Done()
	}()

	<-processedItem
	node.Stop()
	wg.Wait()
	require.True(t, unsubCalled)
}

func TestNodeProduceItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mb := messagebus.NewMockMessageBus(ctrl)
	output := messagebus.NewMockOutput(ctrl)
	output.EXPECT().Write(gomock.Any())
	mb.EXPECT().ProduceTo(gomock.Any(), gomock.Any()).Return(output)
	node := New(mb)

	node.ProduceItem("test", iterator.Item{})
}
