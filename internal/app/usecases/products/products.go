package products

import (
	"context"
	"errors"
	"fmt"
	"prodigo/internal/app/models"
	"prodigo/internal/app/repository/products"
)

type Service struct {
	repository products.Repository
}

func New(repository products.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateProduct(ctx context.Context, p *models.Product) error {
	if err := s.repository.CreateProduct(ctx, p); err != nil {
		return errors.New("failed to create product")
	}
	return nil
}

func (s *Service) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	if _, err := s.repository.GetAllProducts(ctx); err != nil {
		return nil, errors.New("failed to get all products")
	}
	return nil, nil
}

func (s *Service) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	product, err := s.repository.GetProductByID(ctx, id)
	if err != nil {
		return nil, errors.New("failed to get product")
	}
	return product, nil

}

func (s *Service) UpdateProduct(ctx context.Context, p *models.Product) error {
	if err := s.repository.UpdateProduct(ctx, p); err != nil {
		return errors.New("failed to update product")
	}
	return nil
}

func (s *Service) DeleteProduct(ctx context.Context, id int64) error {
	if err := s.repository.DeleteProduct(ctx, id); err != nil {
		return errors.New("failed to delete product")
	}
	return nil
}

func (s *Service) UpdateProductStatus(ctx context.Context, id int64, status string) error {
	p, err := s.repository.GetProductByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}
	p.Status = status
	if err := s.repository.UpdateProduct(ctx, p); err != nil {
		return errors.New("failed to update product status")
	}
	return nil
}
