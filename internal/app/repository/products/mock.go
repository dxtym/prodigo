package products

import (
	"context"
	"prodigo/internal/app/models"

	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateProduct(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockRepo) GetAllProducts(ctx context.Context, fs *models.ProductFilterSearch) ([]*models.Product, error) {
	args := m.Called(ctx, fs)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockRepo) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockRepo) UpdateProduct(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockRepo) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepo) UpdateProductStatus(ctx context.Context, id int64, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockRepo) RestoreProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
