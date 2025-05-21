package usecases

import (
	"go.uber.org/fx"
	"prodigo/internal/app/usecases/categories"
	"prodigo/internal/app/usecases/products"
)

var Module = fx.Module("usecases",
	fx.Provide(
		categories.New,
		products.New,
	),
)
