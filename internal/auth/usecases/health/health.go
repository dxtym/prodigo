package health

import (
	"context"
	"fmt"

	"prodigo/internal/auth/repository/health"
)

type Service interface {
	Check(context.Context) error
}

type service struct {
	repository health.Repository
}

func New(repository health.Repository) Service {
	return &service{repository: repository}
}

func (s *service) Check(ctx context.Context) error {
	if err := s.repository.Check(ctx); err != nil {
		return fmt.Errorf("failed to check health: %w", err)
	}
	return nil
}
