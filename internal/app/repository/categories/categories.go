package categories

import (
	"context"
	"errors"
	"prodigo/internal/app/models"
	"prodigo/pkg/db/postgres"

	"go.uber.org/fx"
)

type Repository interface {
	CreateCategory(ctx context.Context, c *models.Category) error
	UpdateCategory(ctx context.Context, c *models.Category) error
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	CategoryStatistics(ctx context.Context) ([]*models.CategoryStats, error)
}

type Params struct {
	fx.In

	Pool postgres.Pool `name:"app_postgres"`
}

type repository struct {
	pool postgres.Pool `name:"app_postgres"`
}

func New(p Params) Repository {
	return &repository{pool: p.Pool}
}

func (r *repository) CreateCategory(ctx context.Context, c *models.Category) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO categories (name) 
		 VALUES ($1) 
		 RETURNING id, created_at, updated_at`,
		c.Name,
	)
	if err != nil {
		return errors.New("failed to create category")
	}
	return nil
}

func (r *repository) UpdateCategory(ctx context.Context, c *models.Category) error {
	cmd, err := r.pool.Exec(ctx,
		`UPDATE categories 
		 SET name = $1, updated_at = NOW() 
		 WHERE id = $2 AND deleted_at IS NULL`,
		c.Name, c.ID,
	)
	if err != nil {
		return errors.New("failed to update category")
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("category not found or deleted")
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
		return nil, errors.New("failed to get all categories")
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var c models.Category
		if err = rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, errors.New("failed to scan category")
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
		return errors.New("failed to delete/archive category")
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("category not found or deleted")
	}
	return nil
}

func (r *repository) CategoryStatistics(ctx context.Context) ([]*models.CategoryStats, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT 
			c.id, c.name,
			COUNT(p.id),
			SUM(p.quantity),
			SUM(p.price * p.quantity)
		FROM categories AS c
		LEFT JOIN products AS p ON p.category_id = c.id AND p.deleted_at IS NULL
		WHERE c.deleted_at IS NULL
		GROUP BY c.id, c.name
		ORDER BY c.name`)
	if err != nil {
		return nil, errors.New("failed to get category statistics" + err.Error() + "")
	}
	defer rows.Close()

	var stats []*models.CategoryStats
	for rows.Next() {
		var s models.CategoryStats
		if err := rows.Scan(&s.CategoryID, &s.CategoryName, &s.ProductCount, &s.TotalQuantity, &s.TotalValue); err != nil {
			return nil, errors.New("failed to scan category statistics")
		}
		stats = append(stats, &s)
	}
	return stats, nil
}
