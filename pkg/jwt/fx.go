package jwt

import (
	"prodigo/pkg/config"

	"go.uber.org/fx"
)

var Module = fx.Module("jwt",
	fx.Provide(
		func(conf *config.Config) (TokenMaker, error) {
			return New(conf.AuthSecretKey)
		},
	),
)
