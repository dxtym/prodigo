package main

import (
	"context"
	"fmt"
	"prodigo/internal/auth/repository"
	"prodigo/internal/auth/rest"
	"prodigo/internal/auth/rest/handlers"
	"prodigo/internal/auth/usecases"
	"prodigo/pkg/config"
	"prodigo/pkg/db"
	"prodigo/pkg/jwt"
	"prodigo/pkg/migration"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		db.Module,
		repository.Module,
		usecases.Module,
		handlers.Module,
		rest.Module,
		jwt.Module,
		fx.Invoke(func(lc fx.Lifecycle, conf *config.Config) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					if err := migration.Migrate(conf.AuthMigrate, conf.AuthPostgres); err != nil {
						return fmt.Errorf("failed to run auth migrations: %w", err)
					}
					return nil
				},
			})
		}),
	).Run()
}
