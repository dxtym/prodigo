package casbin

import (
	"prodigo/pkg/config"

	"go.uber.org/fx"
)

var Module = fx.Module("casbin",
	fx.Provide(
		func(conf *config.Config) (Enforcer, error) {
			return New(conf.AppCasbin, conf.AppPolicy)
		},
	),
)
