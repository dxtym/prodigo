package products

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"prodigo/internal/app/models"
	"prodigo/internal/app/repository/categories"
	"prodigo/pkg/db/postgres"
	"strings"
)

type Repository interface {
	CreateProduct(ctx context.Context, p *models.Product) error
	GetProductByID(ctx context.Context, id int64) (*models.Product, error)
	GetAllProducts(ctx context.Context, fs *models.ProductFilterSearch) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, p *models.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	RestoreProduct(ctx context.Context, id int64) error
}

var ErrNotFound = errors.New("product not found")

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
    quantity INTEGER NOT NULL check ( quantity > 0 ),
    image TEXT,
    status TEXT NOT NULL,
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
	_, err := r.pool.Exec(ctx, `
		INSERT INTO products (title, category_id, price, quantity, image, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
`, p.Title, p.CategoryID, p.Price, p.Quantity, p.Image, p.Status)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.New("failed to get product: " + err.Error() + "")
	}
	return &p, nil
}

func (r *repository) GetAllProducts(ctx context.Context, fs *models.ProductFilterSearch) ([]*models.Product, error) {
	var (
		args  []interface{}
		where []string
		i     = 1
	)

	if fs.CategoryName != "" {
		where = append(where, fmt.Sprintf("c.name ILIKE $%d", i))
		args = append(args, "%"+fs.CategoryName+"%")
		i++
	}
	if fs.Status != "" {
		where = append(where, fmt.Sprintf("p.status = $%d", i))
		args = append(args, fs.Status)
		i++
	}
	if fs.PriceMin > 0 {
		where = append(where, fmt.Sprintf("p.price >= $%d", i))
		args = append(args, fs.PriceMin)
		i++
	}
	if fs.PriceMax > 0 {
		where = append(where, fmt.Sprintf("p.price <= $%d", i))
		args = append(args, fs.PriceMax)
		i++
	}
	if fs.Search != "" {
		where = append(where, fmt.Sprintf("p.title ILIKE $%d", i))
		args = append(args, "%"+fs.Search+"%")
		i++
	}

	fullQuery := ""
	if len(where) > 0 {
		fullQuery = " AND " + strings.Join(where, " AND ")
	}

	rows, err := r.pool.Query(ctx, fmt.Sprintf(`
		SELECT p.id, p.title, p.category_id, p.price, p.quantity, p.image, p.status, p.created_at, p.updated_at
		FROM products as p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.deleted_at IS NULL %s
	`, fullQuery), args...)
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

func (r *repository) RestoreProduct(ctx context.Context, id int64) error {
	restore, err := r.pool.Exec(ctx, `
	UPDATE products SET deleted_at = NULL WHERE id = $1
`, id)
	if err != nil {
		return errors.New("failed to restore product: " + err.Error() + "")
	}
	if restore.RowsAffected() == 0 {
		return errors.New("product not found")
	}
	return nil
}
