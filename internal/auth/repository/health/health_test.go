package health_test

import (
	"prodigo/internal/auth/repository/health"
	"prodigo/pkg/db/postgres"
	rdb "prodigo/pkg/db/redis"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

func TestNew(t *testing.T) {
	pool := new(postgres.MockPool)
	client := new(rdb.MockClient)

	repository := health.New(pool, client)

	assert.NotNil(t, repository)
}

func TestRepository_Check(t *testing.T) {
	pool := new(postgres.MockPool)
	pool.On("Ping", mock.Anything).Return(nil)
	client := new(rdb.MockClient)
	client.On("Ping", mock.Anything).Return(&redis.StatusCmd{})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repository := health.New(pool, client)
	err := repository.Check(ctx)

	assert.NoError(t, err)
}
