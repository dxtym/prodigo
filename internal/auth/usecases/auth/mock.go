package auth

import (
	"context"
	"prodigo/internal/auth/dto"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Register(ctx context.Context, req dto.RegisterRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockService) Login(ctx context.Context, req dto.LoginRequest) (string, string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockService) Refresh(ctx context.Context, req dto.RefreshRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

var _ Service = (*MockService)(nil)
