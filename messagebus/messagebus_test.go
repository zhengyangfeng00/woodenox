package messagebus

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/zhengyangfeng00/woodenox/iterator"
	"github.com/zhengyangfeng00/woodenox/stream"
)

func TestMessageBus(t *testing.T) {
	t.Run("subscribe", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mb := New().(*messageBus)
		s := stream.NewMockStream(ctrl)
		s.EXPECT().Run()
		s.EXPECT().Stop()
		s.EXPECT().NewSubscriber(gomock.Any())
		mb.newStreamFn = func(
			name string,
			opts ...stream.Option,
		) stream.Stream {
			return s
		}

		mb.Subscribe("test")
		mb.Stop()
	})

	t.Run("produce", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mb := New().(*messageBus)
		s := stream.NewMockStream(ctrl)
		s.EXPECT().Run()
		s.EXPECT().Stop()
		s.EXPECT().Accept(gomock.Any())
		mb.newStreamFn = func(
			name string,
			opts ...stream.Option,
		) stream.Stream {
			return s
		}

		out := mb.ProduceTo("test")
		out.Write(iterator.Item{})
		mb.Stop()
	})
}
