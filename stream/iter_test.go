package stream

import (
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	iterator "github.com/zhengyangfeng00/woodenox/iterator"
)

func TestStreamIter(t *testing.T) {
	t.Run("empty stream", func(t *testing.T) {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		s := NewMockstream(ctl)
		s.EXPECT().
			next(gomock.Any()).
			Return(iterator.Item{}, errNotAvailable).AnyTimes()

		it := newStreamIter(s, false)
		require.False(t, it.Next())
	})

	t.Run("unblock by new item", func(t *testing.T) {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		var (
			s                 = NewMockstream(ctl)
			firstCallReturned = make(chan struct{})
			first             = true
		)
		s.EXPECT().next(gomock.Any()).DoAndReturn(func(it *streamIter) (iterator.Item, error) {
			if first {
				close(firstCallReturned)
				first = false
				return iterator.Item{}, errNotAvailable
			}
			return iterator.Item{}, nil
		}).AnyTimes()

		// N.B: block when no new item is available.
		it := newStreamIter(s, true)

		// Keep calling notify on the iterator even if the call is
		// superfluous.
		done := make(chan struct{})
		go func() {
			// N.B: wait for the first call to Next to return false
			// before notifying to verify the unblock behavior.
			<-firstCallReturned
			for {
				select {
				case <-done:
					return
				default:
					it.notify()
				}
			}
		}()

		// Assert that the call to Next eventually unblocks.
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			it.Next()
		}()

		wg.Wait()
		close(done)
	})

	t.Run("unblock by close", func(t *testing.T) {
		ctl := gomock.NewController(t)
		defer ctl.Finish()

		s := NewMockstream(ctl)
		s.EXPECT().next(gomock.Any()).
			Return(iterator.Item{}, errNotAvailable).AnyTimes()

		it := newStreamIter(s, true)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			require.False(t, it.Next())
		}()
		it.close()
		wg.Wait()
	})
}
