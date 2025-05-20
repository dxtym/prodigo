package products

import (
	"context"
	"errors"
	"fmt"
	"prodigo/internal/app/models"
	"prodigo/internal/app/repository/categories"
	"prodigo/pkg/db/postgres"
)

type Repository interface {
	CreateProduct(ctx context.Context, p *models.Product) error
	GetProductByID(ctx context.Context, id int64) (*models.Product, error)
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, p *models.Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

type repository struct {
	pool postgres.Pool
}

func New(pool postgres.Pool, categoriesRepo categories.Repository) (Repository, error) {
	_, err := pool.Exec(context.Background(), `
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    category_id INTEGER references categories(id),
    price INTEGER NOT NULL check (price > 0),
    quantity INTEGER NOT NULL check ( price > 0 ),
    image TEXT,
    status TEXT NOT NULL DEFAULT 'available',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table products: %w", err)
	}
	return &repository{
		pool: pool,
	}, nil
}

func (r *repository) CreateProduct(ctx context.Context, p *models.Product) error {
	err := r.pool.QueryRow(ctx, `
		INSERT INTO products (title, category_id, price, quantity, image, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
`, p.Title, p.CategoryID, p.Price, p.Quantity, p.Image, p.Status).Scan(&p.ID)
	if err != nil {
		return errors.New("failed to create product" + err.Error() + "")
	}
	return nil
}

func (r *repository) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	var p models.Product
	err := r.pool.QueryRow(ctx, `
		SELECT id, title, category_id, price, quantity, image, status, created_at, updated_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL

`, id).Scan(&p.ID, &p.Title, &p.CategoryID, &p.Price, &p.Quantity, &p.Image, &p.Status, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, errors.New("failed to get product: " + err.Error() + "")
	}
	return &p, nil
}

func (r *repository) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, category_id, price, quantity, image, status, created_at, updated_at
		FROM products
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, errors.New("failed to get all products: " + err.Error() + "")
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID, &p.Title, &p.CategoryID, &p.Price,
			&p.Quantity, &p.Image, &p.Status,
			&p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, errors.New("failed to scan product: " + err.Error() + "")
		}
		products = append(products, &p)
	}
	return products, nil
}

func (r *repository) UpdateProduct(ctx context.Context, p *models.Product) error {
	upd, err := r.pool.Exec(ctx, `
	UPDATE products
	SET title = $1, category_id = $2, price = $3, quantity = $4, image = $5, status = $6, updated_at = NOW()
	WHERE id = $7 AND deleted_at IS NULL
`, p.Title, p.CategoryID, p.Price, p.Quantity, p.Image, p.Status, p.ID)
	if err != nil {
		return errors.New("failed to update product: " + err.Error() + "")
	}
	if upd.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *repository) DeleteProduct(ctx context.Context, id int64) error {
	dlt, err := r.pool.Exec(ctx, `
	UPDATE products SET deleted_at = NOW() WHERE id = $1
`, id)
	if err != nil {
		return errors.New("failed to delete product: " + err.Error() + "")
	}
	if dlt.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}
