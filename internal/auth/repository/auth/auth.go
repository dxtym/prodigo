package auth

import (
	"context"
	"errors"
	"fmt"
	"prodigo/internal/auth/models"
	db "prodigo/pkg/db/postgres"
	rdb "prodigo/pkg/db/redis"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	CreateUser(context.Context, *models.User) error
	GetByUsername(context.Context, string) (*models.User, error)
	SaveToken(context.Context, int64, string, time.Duration) error
	GetToken(context.Context, int64) (string, error)
}

type repository struct {
	pool   db.Pool
	client rdb.Client
}

func New(pool db.Pool, client rdb.Client) (Repository, error) {
	if _, err := pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		username VARCHAR(20) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
		deleted_at TIMESTAMP
	);`); err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	return &repository{
		pool:   pool,
		client: client,
	}, nil
}

func (r *repository) CreateUser(ctx context.Context, user *models.User) error {
	cmd, err := r.pool.Exec(ctx, `	
	INSERT INTO users (
		username, 
		password
	) VALUES ($1, $2);`, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrUserNotCreated
	}

	return nil
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.pool.QueryRow(ctx, `
	SELECT
		id,
		username,
		password
	FROM users
	WHERE username = $1 AND deleted_at IS NULL;`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

func (r *repository) SaveToken(ctx context.Context, userID int64, token string, duration time.Duration) error {
	key := fmt.Sprintf("user:token:%d", userID)
	if err := r.client.Set(ctx, key, token, duration).Err(); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	return nil
}

func (r *repository) GetToken(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf("user:token:%d", userID)
	token, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrTokenNotFound
		}
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	return token, nil
}
