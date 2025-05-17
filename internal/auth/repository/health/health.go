package health

import (
	"context"
	"fmt"

	"prodigo/pkg/db/postgres"
	"prodigo/pkg/db/redis"
)

type Repository interface {
	Check(context.Context) error
}

type repository struct {
	pool   postgres.Pool
	client redis.Client
}

func New(pool postgres.Pool, client redis.Client) Repository {
	return &repository{
		pool:   pool,
		client: client,
	}
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
