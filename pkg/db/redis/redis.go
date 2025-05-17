package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client interface {
	Set(context.Context, string, any, time.Duration) *redis.StatusCmd
	Get(context.Context, string) *redis.StringCmd
	Ping(context.Context) *redis.StatusCmd
}

type conn struct {
	*redis.Client
}

func New(ctx context.Context, addr, pass string) (Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &conn{Client: client}, nil
}
