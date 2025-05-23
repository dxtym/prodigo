package categories

import (
	"context"
	"prodigo/internal/app/models"

	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateCategory(ctx context.Context, c *models.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockRepo) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *MockRepo) UpdateCategory(ctx context.Context, c *models.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockRepo) DeleteCategory(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepo) CategoryStatistics(ctx context.Context) ([]*models.CategoryStats, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.CategoryStats), args.Error(1)
}
