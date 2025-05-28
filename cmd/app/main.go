package main

import (
	"context"
	"fmt"
	"prodigo/internal/app/repository"
	"prodigo/internal/app/rest"
	"prodigo/internal/app/rest/casbin"
	"prodigo/internal/app/rest/handlers"
	"prodigo/internal/app/rest/middleware"
	"prodigo/internal/app/usecases"
	"prodigo/pkg/config"
	"prodigo/pkg/db"
	"prodigo/pkg/jwt"
	"prodigo/pkg/migration"

	"go.uber.org/fx"
)

// запуск

func main() {
	fx.New(
		config.Module,
		repository.Module,
		db.Module,
		usecases.Module,
		handlers.Module,
		rest.Module,
		middleware.Module,
		jwt.Module,
		casbin.Module,
		fx.Invoke(func(lc fx.Lifecycle, conf *config.Config) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					if err := migration.Migrate(conf.AppMigrate, conf.AppPostgres); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
			})
		}),
	).Run()
}
