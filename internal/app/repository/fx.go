package repository

import (
	"go.uber.org/fx"
	"prodigo/internal/app/repository/categories"
	"prodigo/internal/app/repository/products"
)

var Module = fx.Module("repository",
	fx.Provide(
		categories.New,
		products.New,
	),
)
