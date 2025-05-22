package health_test

import (
	"context"
	"errors"
	healthRepository "prodigo/internal/auth/repository/health"
	healthService "prodigo/internal/auth/usecases/health"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Check(t *testing.T) {
	tests := []struct {
		name    string
		wantErr error
	}{
		{name: "success", wantErr: nil},
		{name: "invalid check", wantErr: errors.New("some error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := new(healthRepository.MockRepository)
			defer repository.AssertExpectations(t)

			repository.On("Check", mock.Anything).Return(tt.wantErr).Once()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			service := healthService.New(repository)

			err := service.Check(ctx)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
