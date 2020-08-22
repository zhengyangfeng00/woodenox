package stream

type SubscribeOpt func(*subscribeOpts)

func WithBlock() SubscribeOpt {
	return func(opts *subscribeOpts) {
		opts.block = true
	}
}

type subscribeOpts struct {
	block bool
}

func newSubOpts() *subscribeOpts {
	return &subscribeOpts{}
}

func newSubscribeOpts() *subscribeOpts {
	return &subscribeOpts{}
}

type ProduceOpt func(*produceOpts)

type produceOpts struct{}
