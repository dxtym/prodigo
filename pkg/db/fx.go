package db

import (
	"context"
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
)
