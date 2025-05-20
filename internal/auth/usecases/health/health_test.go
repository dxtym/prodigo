package health_test

import (
	"context"
	"prodigo/internal/auth/repository/health"
	healthService "prodigo/internal/auth/usecases/health"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	repository := new(health.MockRepository)

	service := healthService.New(repository)

	assert.NotNil(t, service)
}

func TestService_Check(t *testing.T) {
	repository := new(health.MockRepository)
	repository.On("Check", mock.Anything).Return(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	service := healthService.New(repository)
	err := service.Check(ctx)

	assert.NoError(t, err)
}
