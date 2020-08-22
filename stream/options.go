package stream

import "time"

type Option func(*options)

func WithCapacity(capacity int) Option {
	return func(o *options) {
		o.capacity = capacity
	}
}

func WithPurgeInterval(interval time.Duration) Option {
	return func(o *options) {
		o.purgeInterval = interval
	}
}

type options struct {
	capacity      int
	purgeInterval time.Duration
}

func newOptions() *options {
	return &options{
		capacity:      1000,
		purgeInterval: 10 * time.Second,
	}
}
