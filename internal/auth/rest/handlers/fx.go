package handlers

import (
	"prodigo/internal/auth/rest/handlers/health"

	"go.uber.org/fx"
)

var Module = fx.Module("handlers",
	fx.Provide(
		health.New,
	),
)
