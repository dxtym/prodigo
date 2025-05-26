package categories

import (
	"context"
	"errors"
	"prodigo/internal/app/models"
	"prodigo/internal/app/repository/categories"
)

type Service struct {
	repository categories.Repository
}

type ServiceInterface interface {
	CreateCategory(ctx context.Context, c *models.Category) error
	UpdateCategory(ctx context.Context, c *models.Category) error
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	CategoryStatistics(ctx context.Context) ([]*models.CategoryStats, error)
}

func New(repository categories.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateCategory(ctx context.Context, c *models.Category) error {
	if err := s.repository.CreateCategory(ctx, c); err != nil {
		return errors.New("failed to create category")
	}
	return nil
}

func (s *Service) UpdateCategory(ctx context.Context, c *models.Category) error {
	if err := s.repository.UpdateCategory(ctx, c); err != nil {
		return errors.New("failed to update category")
	}
	return nil
}

func (s *Service) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	cats, err := s.repository.GetAllCategories(ctx)
	if err != nil {
		return nil, errors.New("failed to get all categories")
	}
	return cats, nil
}

func (s *Service) DeleteCategory(ctx context.Context, id int64) error {
	if err := s.repository.DeleteCategory(ctx, id); err != nil {
		return errors.New("failed to delete category")
	}
	return nil
}

func (s *Service) CategoryStatistics(ctx context.Context) ([]*models.CategoryStats, error) {
	stats, err := s.repository.CategoryStatistics(ctx)
	if err != nil {
		return nil, errors.New("failed to get category statistics")
	}
	return stats, nil
}
