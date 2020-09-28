package stream

import (
	"time"

	"go.uber.org/zap"
)

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

func WithLogger(l *zap.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}

type options struct {
	capacity       int
	purgeInterval  time.Duration
	notifyInterval time.Duration
	logger         *zap.Logger
}

func newOptions() *options {
	return &options{
		capacity:       1000,
		purgeInterval:  10 * time.Second,
		notifyInterval: time.Second,
		logger:         zap.NewNop(),
	}
}
