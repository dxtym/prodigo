package db

import (
	"context"
	"fmt"
	"prodigo/pkg/config"
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
	),
	fx.Invoke(func(lc fx.Lifecycle, conf *config.Config) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := migrateDB(conf.MigrateURL, conf.PostgresDSN); err != nil {
					return fmt.Errorf("failed to migrate database: %w", err)
				}
				return nil
			},
		})
	}),
)
