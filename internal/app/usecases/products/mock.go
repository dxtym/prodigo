package products

import (
	"context"
	"github.com/stretchr/testify/mock"
	"prodigo/internal/app/models"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateProduct(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockService) GetAllProducts(ctx context.Context, fs *models.ProductFilterSearch) ([]*models.Product, error) {
	args := m.Called(ctx, fs)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockService) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockService) UpdateProduct(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockService) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockService) UpdateProductStatus(ctx context.Context, id int64, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockService) RestoreProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
