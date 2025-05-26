package products

import (
	"context"
	"errors"
	"fmt"
	"prodigo/internal/app/models"
	"prodigo/internal/app/repository/products"
)

type ServiceInterface interface {
	CreateProduct(ctx context.Context, p *models.Product) error
	GetProduct(ctx context.Context, id int64) (*models.Product, error)
	GetAllProducts(ctx context.Context, fs *models.ProductFilterSearch) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, p *models.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	RestoreProduct(ctx context.Context, id int64) error
	UpdateProductStatus(ctx context.Context, id int64, status string) error
}

var ErrNotFound = errors.New("product not found")

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

func (s *Service) GetAllProducts(ctx context.Context, fs *models.ProductFilterSearch) ([]*models.Product, error) {
	prods, err := s.repository.GetAllProducts(ctx, fs)
	if err != nil {
		return nil, errors.New("failed to get all products")
	}
	return prods, nil
}

func (s *Service) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	product, err := s.repository.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, products.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, errors.New("failed to get product")
	}
	return product, nil

}

func (s *Service) UpdateProduct(ctx context.Context, p *models.Product) error {
	update, err := s.repository.GetProductByID(ctx, p.ID)

	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if p.Title != "" {
		update.Title = p.Title
	}
	if p.CategoryID != 0 {
		update.CategoryID = p.CategoryID
	}
	if p.Price != 0 {
		update.Price = p.Price
	}
	if p.Quantity != 0 {
		update.Quantity = p.Quantity
	}
	if p.Image != "" {
		update.Image = p.Image
	}

	if p.Status != "" {
		update.Status = p.Status
	}

	if err = s.repository.UpdateProduct(ctx, update); err != nil {
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

func (s *Service) RestoreProduct(ctx context.Context, id int64) error {
	if err := s.repository.RestoreProduct(ctx, id); err != nil {
		return errors.New("failed to restore product")
	}
	return nil
}
