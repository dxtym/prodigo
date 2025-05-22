package auth

import (
	"context"
	"prodigo/internal/auth/models"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) SaveToken(ctx context.Context, userID int64, token string, duration time.Duration) error {
	args := m.Called(ctx, userID, token, duration)
	return args.Error(0)
}

func (m *MockRepository) GetToken(ctx context.Context, userID int64) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

var _ Repository = (*MockRepository)(nil)
