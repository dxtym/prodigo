package health_test

import (
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
	pool := new(db.MockPool)
	defer pool.AssertExpectations(t)

	client := new(rdb.MockClient)
	defer client.AssertExpectations(t)

	pool.On("Ping", mock.Anything).Return(nil).Once()

	client.On("Ping", mock.Anything).Return(&redis.StatusCmd{}).Once()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repository := health.New(pool, client)

	err := repository.Check(ctx)
	assert.NoError(t, err)
}
