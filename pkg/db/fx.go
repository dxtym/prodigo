package db

import (
	"context"
	"fmt"
	"prodigo/pkg/config"
	"prodigo/pkg/db/migration"
	"prodigo/pkg/db/postgres"
	"prodigo/pkg/db/redis"

	"go.uber.org/fx"
)

var Module = fx.Module("db",
	fx.Provide(
		func(conf *config.Config) (postgres.Pool, error) {
			return postgres.New(context.Background(), conf.PostgresDSN)
		},
		func(conf *config.Config) (redis.Client, error) {
			return redis.New(context.Background(), conf.RedisAddr, conf.RedisPass)
		},
		func(conf *config.Config) (*migration.Migration, error) {
			return migration.New(conf.MigrateURL, conf.PostgresDSN)
		},
	),
	fx.Invoke(func(lc fx.Lifecycle, mg *migration.Migration) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := mg.Up(); err != nil {
					return fmt.Errorf("failed to apply migrations: %w", err)
				}
				return nil
			},
			OnStop: func(ctx context.Context) error {
				if err := mg.Down(); err != nil {
					return fmt.Errorf("failed to revert migrations: %w", err)
				}
				return nil
			},
		})
	}),
)
