package health

import (
	"context"
	"fmt"

	"prodigo/pkg/db/postgres"
	"prodigo/pkg/db/redis"

	"go.uber.org/fx"
)

type Repository interface {
	Check(context.Context) error
}

type RepositoryParams struct {
	fx.In

	Pool   postgres.Pool `name:"auth_postgres"`
	Client redis.Client  `name:"auth_redis"`
}

type repository struct {
	pool   postgres.Pool `name:"auth_postgres"`
	client redis.Client  `name:"auth_redis"`
}

func New(p RepositoryParams) Repository {
	return &repository{pool: p.Pool, client: p.Client}
}

func (r *repository) Check(ctx context.Context) error {
	if err := r.pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	return nil
}
