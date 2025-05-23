package health

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Check(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

var _ Service = (*MockService)(nil)
