package handlers

import (
	"go.uber.org/fx"
	"prodigo/internal/app/rest/handlers/categories"
	"prodigo/internal/app/rest/handlers/products"
)

var Module = fx.Module("handlers",
	fx.Provide(
		categories.New,
		products.New,
	),
)
