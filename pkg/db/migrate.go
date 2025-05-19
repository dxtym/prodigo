package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func migrateDB(source, destination string) error {
	mg, err := migrate.New(source, destination)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := mg.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if srcErr, dstErr := mg.Close(); srcErr != nil || dstErr != nil {
		return fmt.Errorf("failed to close migration instance: srcErr=%v, dstErr=%v", srcErr, dstErr)
	}

	return nil
}
