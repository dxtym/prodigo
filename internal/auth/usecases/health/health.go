package health

import (
	"context"
	"fmt"

	"prodigo/internal/auth/repository/health"
)

type Service struct {
	repository health.Repository
}

func New(repository health.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Check(ctx context.Context) error {
	if err := s.repository.Check(ctx); err != nil {
		return fmt.Errorf("failed to check health: %w", err)
	}
	return nil
}
