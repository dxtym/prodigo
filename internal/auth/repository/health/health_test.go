package health_test

import (
	"errors"
	"prodigo/internal/auth/repository/health"
	db "prodigo/pkg/db/postgres"
	rdb "prodigo/pkg/db/redis"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestRepository_Check(t *testing.T) {
	tests := []struct {
		name  string
		build func(*db.MockPool, *rdb.MockClient)
		check func(error)
	}{
		{
			name: "success",
			build: func(pool *db.MockPool, client *rdb.MockClient) {
				pool.On("Ping", mock.Anything).Return(nil).Once()

				cmd := redis.NewStatusCmd(context.Background())

				client.On("Ping", mock.Anything).Return(cmd).Once()
			},
			check: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "invalid pool",
			build: func(pool *db.MockPool, client *rdb.MockClient) {
				pool.On("Ping", mock.Anything).Return(errors.New("pool error")).Once()
			},
			check: func(err error) {
				assert.ErrorContains(t, err, errors.New("pool error").Error())
			},
		},
		{
			name: "invalid client",
			build: func(pool *db.MockPool, client *rdb.MockClient) {
				pool.On("Ping", mock.Anything).Return(nil).Once()

				cmd := redis.NewStatusCmd(context.Background())
				cmd.SetErr(errors.New("client error"))

				client.On("Ping", mock.Anything).Return(cmd).Once()
			},
			check: func(err error) {
				assert.ErrorContains(t, err, errors.New("client error").Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := new(db.MockPool)
			defer pool.AssertExpectations(t)

			client := new(rdb.MockClient)
			defer client.AssertExpectations(t)

			tt.build(pool, client)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			repository := health.New(health.Params{
				Pool:   pool,
				Client: client,
			})

			err := repository.Check(ctx)
			tt.check(err)
		})
	}
}
