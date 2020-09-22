package woodenoxfx

import (
	"github.com/zhengyangfeng00/woodenox/messagebus"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(newModule),
)

type Params struct {
	fx.In
}

type Result struct {
	fx.Out

	MessageBus messagebus.MessageBus
}

func newModule(p Params) Result {
	return Result{
		MessageBus: messagebus.New(),
	}
}
