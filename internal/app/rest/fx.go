package rest

import "go.uber.org/fx"

var Module = fx.Module("rest",
	fx.Provide(
		New,
	),
)
