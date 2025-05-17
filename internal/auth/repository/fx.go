package repository

import (
	"prodigo/internal/auth/repository/health"

	"go.uber.org/fx"
)

var Module = fx.Module("repository",
	fx.Provide(
		health.New,
	),
)
