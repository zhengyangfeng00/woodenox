package woodenoxfx

import (
	"github.com/zhengyangfeng00/woodenox/messagebus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(newModule),
)

type Params struct {
	fx.In

	Logger *zap.Logger
}

type Result struct {
	fx.Out

	MessageBus messagebus.MessageBus
}

func newModule(p Params) Result {
	return Result{
		MessageBus: messagebus.New(p.Logger),
	}
}
