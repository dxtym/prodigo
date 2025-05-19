package migration

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

type Migration struct {
	mg *migrate.Migrate
}

func New(src, dest string) (*Migration, error) {
	mg, err := migrate.New(src, dest)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return &Migration{mg: mg}, nil
}

func (m *Migration) Up() error {
	if err := m.mg.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func (m *Migration) Down() error {
	if err := m.mg.Down(); err != nil {
		return fmt.Errorf("failed to revert migrations: %w", err)
	}

	return nil
}