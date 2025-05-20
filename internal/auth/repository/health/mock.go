package health

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Check(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}