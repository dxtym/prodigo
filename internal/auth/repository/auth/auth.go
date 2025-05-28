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
	"go.uber.org/fx"
)

type Repository interface {
	CreateUser(context.Context, *models.User) error
	GetByUsername(context.Context, string) (*models.User, error)
	SaveToken(context.Context, int64, string, time.Duration) error
	GetToken(context.Context, int64) (string, error)
}

type Params struct {
	fx.In

	Pool   db.Pool    `name:"auth_postgres"`
	Client rdb.Client `name:"auth_redis"`
}

type repository struct {
	pool   db.Pool
	client rdb.Client
}

func New(p Params) Repository {
	return &repository{pool: p.Pool, client: p.Client}
}

func (r *repository) CreateUser(ctx context.Context, user *models.User) error {
	if _, err := r.pool.Exec(ctx, `	
	INSERT INTO users (
		username, 
		password
	) VALUES ($1, $2);`, user.Username, user.Password); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.pool.QueryRow(ctx, `
	SELECT
		id,
		username,
		password,
		role
	FROM users
	WHERE username = $1 AND deleted_at IS NULL;`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
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
