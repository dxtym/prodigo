package usecases

import (
	"prodigo/internal/auth/usecases/auth"
	"prodigo/internal/auth/usecases/health"

	"go.uber.org/fx"
)

var Module = fx.Module("usecases",
	fx.Provide(
		health.New,
		auth.New,
	),
)
