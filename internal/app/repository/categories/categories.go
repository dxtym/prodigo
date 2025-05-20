package categories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"prodigo/internal/app/models"
	"prodigo/pkg/db/postgres"
)

type Repository interface {
	CreateCategory(ctx context.Context, c *models.Category) error
	UpdateCategory(ctx context.Context, c *models.Category) error
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
}

type repository struct {
	pool postgres.Pool
}

func New(pool postgres.Pool) (Repository, error) {
	_, err := pool.Exec(context.Background(), `
CREATE TABLE IF NOT EXISTS categories (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMP
)`)
	if err != nil {
		return nil, errors.New("failed to create categories table")
	}
	log.Println("Categories table created successfully")

	return &repository{
		pool: pool,
	}, nil
}

func (r *repository) CreateCategory(ctx context.Context, c *models.Category) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO categories (name) 
		 VALUES ($1) 
		 RETURNING id, created_at, updated_at`,
		c.Name,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *repository) UpdateCategory(ctx context.Context, c *models.Category) error {
	cmd, err := r.pool.Exec(ctx,
		`UPDATE categories 
		 SET name = $1, updated_at = NOW() 
		 WHERE id = $2 AND deleted_at IS NULL`,
		c.Name, c.ID,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("category not found or deleted")
	}
	return nil
}

func (r *repository) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, created_at, updated_at, deleted_at 
		 FROM categories 
		 WHERE deleted_at IS NULL`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r *repository) DeleteCategory(ctx context.Context, id int64) error {
	cmd, err := r.pool.Exec(ctx,
		`UPDATE categories 
		 SET deleted_at = NOW(), updated_at = NOW() 
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("category not found or deleted")
	}
	return nil
}
