package iterator

type Item struct {
	Timestamp int64
	Data      []byte
}

type Iterator interface {
	Curr() Item
	Next() bool
	Err() error
}
