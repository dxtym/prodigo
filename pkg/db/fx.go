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
		fx.Annotate(
			func(conf *config.Config) (postgres.Pool, error) {
				return postgres.New(context.Background(), conf.AuthPostgres)
			},
			fx.ResultTags(`name:"auth_postgres"`),
		),
		fx.Annotate(
			func(conf *config.Config) (redis.Client, error) {
				return redis.New(context.Background(), conf.AuthRedis)
			},
			fx.ResultTags(`name:"auth_redis"`),
		),
		fx.Annotate(
			func(conf *config.Config) (postgres.Pool, error) {
				return postgres.New(context.Background(), conf.AppPostgres)
			},
			fx.ResultTags(`name:"app_postgres"`),
		),
	),
)
