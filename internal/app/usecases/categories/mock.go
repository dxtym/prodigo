package categories

import (
	"context"
	"github.com/stretchr/testify/mock"
	"prodigo/internal/app/models"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateCategory(ctx context.Context, c *models.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockService) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *MockService) UpdateCategory(ctx context.Context, c *models.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockService) DeleteCategory(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockService) CategoryStatistics(ctx context.Context) ([]*models.CategoryStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.CategoryStats), args.Error(1)
}
