package stream

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zhengyangfeng00/woodenox/iterator"
)

func TestStream(t *testing.T) {
	t.Run("accept 1 item", func(t *testing.T) {
		s := New("test").(*streamImpl)
		done := make(chan struct{})
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					// N.B: ignore errors.
					_ = s.Accept(iterator.Item{})
				}
			}
		}()
		it, unsub := s.NewSubscriber(WithBlock())
		require.True(t, it.Next())
		close(done)
		unsub()
	})

	t.Run("stream full", func(t *testing.T) {
		s := New("test").(*streamImpl)
		for {
			if s.Accept(iterator.Item{}) == errStreamFull {
				return
			}
		}
	})
}

func TestStreamPurge(t *testing.T) {
	t.Run("no subscriber", func(t *testing.T) {
		s := New("test").(*streamImpl)
		s.Accept(iterator.Item{})
		s.purge()
		require.Equal(t, 0, len(s.buf))
	})

	t.Run("purge 1 item", func(t *testing.T) {
		s := New("test").(*streamImpl)
		it, unsub := s.NewSubscriber(WithBlock())
		s.Accept(iterator.Item{})
		s.Accept(iterator.Item{})

		require.True(t, it.Next())
		// it now points to the second item and the first item can be
		// purged.
		require.True(t, it.Next())

		s.purge()
		require.Equal(t, 1, len(s.buf))
		require.Equal(t, int64(1), s.startOffset)
		unsub()
	})
}
