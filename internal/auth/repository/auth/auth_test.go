package auth_test

import (
	"context"
	"errors"
	"prodigo/internal/auth/models"
	"prodigo/internal/auth/repository/auth"
	db "prodigo/pkg/db/postgres"

	rdb "prodigo/pkg/db/redis"
	"prodigo/pkg/utils"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		arg     *models.User
		want    pgconn.CommandTag
		wantErr error
	}{
		{
			name: "success",
			arg: &models.User{
				Username: utils.GenerateRandomString(10),
				Password: utils.GenerateRandomString(10),
			},
			want:    pgconn.NewCommandTag("INSERT 1"),
			wantErr: nil,
		},
		{
			name: "invalid username",
			arg: &models.User{
				Username: "",
				Password: utils.GenerateRandomString(10),
			},
			want:    pgconn.CommandTag{},
			wantErr: errors.New("failed to create user"),
		},
		{
			name: "invalid password",
			arg: &models.User{
				Username: utils.GenerateRandomString(10),
				Password: "",
			},
			want:    pgconn.CommandTag{},
			wantErr: errors.New("failed to create user"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := new(db.MockPool)
			require.NotNil(t, pool)
			defer pool.AssertExpectations(t)

			client := new(rdb.MockClient)
			require.NotNil(t, client)
			defer client.AssertExpectations(t)

			pool.On("Exec",
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(pgconn.CommandTag{}, nil).Once()
			pool.On("Exec",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(tt.want, tt.wantErr).Once()

			repository := auth.New(auth.RepositoryParams{
				Pool:   pool,
				Client: client,
			})
			require.NotNil(t, repository)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := repository.CreateUser(ctx, tt.arg)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRepository_GetByUsername(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    *models.User
		wantErr error
	}{
		{
			name:    "success",
			arg:     utils.GenerateRandomString(10),
			want:    &models.User{},
			wantErr: nil,
		},
		{
			name:    "user not found",
			arg:     utils.GenerateRandomString(10),
			want:    nil,
			wantErr: auth.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := new(db.MockPool)
			require.NotNil(t, pool)
			defer pool.AssertExpectations(t)

			row := new(db.MockRow)
			require.NotNil(t, row)
			defer row.AssertExpectations(t)

			client := new(rdb.MockClient)
			require.NotNil(t, client)
			defer client.AssertExpectations(t)

			pool.On("Exec",
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(pgconn.CommandTag{}, nil).Once()
			pool.On("QueryRow",
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(row).Once()

			row.On("Scan",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(tt.wantErr).Once()

			repository := auth.New(auth.RepositoryParams{
				Pool:   pool,
				Client: client,
			})
			require.NotNil(t, repository)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			user, err := repository.GetByUsername(ctx, tt.arg)
			assert.Equal(t, user, tt.want)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRepository_SaveToken(t *testing.T) {
	type arg struct {
		key      int64
		value    string
		duration time.Duration
	}

	tests := []struct {
		name    string
		arg     arg
		wantErr error
	}{
		{
			name: "success",
			arg: arg{
				key:      utils.GenerateRandomInt(10),
				value:    utils.GenerateRandomString(10),
				duration: time.Minute,
			},
			wantErr: nil,
		},
		{
			name: "invalid value",
			arg: arg{
				key:      utils.GenerateRandomInt(10),
				value:    "",
				duration: time.Minute,
			},
			wantErr: errors.New("failed to save token"),
		},
		{
			name: "invalid duration",
			arg: arg{
				key:      1,
				value:    utils.GenerateRandomString(10),
				duration: -time.Minute,
			},
			wantErr: errors.New("failed to save token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := new(db.MockPool)
			require.NotNil(t, pool)
			defer pool.AssertExpectations(t)

			client := new(rdb.MockClient)
			require.NotNil(t, client)
			defer client.AssertExpectations(t)

			pool.On("Exec",
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(pgconn.CommandTag{}, nil).Once()

			cmd := redis.NewStatusCmd(context.Background())
			cmd.SetErr(tt.wantErr)

			client.On("Set",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(cmd).Once()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			repository := auth.New(auth.RepositoryParams{
				Pool:   pool,
				Client: client,
			})
			require.NotNil(t, repository)

			err := repository.SaveToken(ctx, tt.arg.key, tt.arg.value, tt.arg.duration)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRepository_GetToken(t *testing.T) {
	tests := []struct {
		name    string
		arg     int64
		want    string
		wantErr error
	}{
		{
			name:    "success",
			arg:     utils.GenerateRandomInt(10),
			want:    utils.GenerateRandomString(10),
			wantErr: nil,
		},
		{
			name:    "token not found",
			arg:     utils.GenerateRandomInt(10),
			want:    "",
			wantErr: auth.ErrTokenNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := new(db.MockPool)
			require.NotNil(t, pool)
			defer pool.AssertExpectations(t)

			client := new(rdb.MockClient)
			require.NotNil(t, client)
			defer client.AssertExpectations(t)

			pool.On("Exec",
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(pgconn.CommandTag{}, nil).Once()

			cmd := redis.NewStringCmd(context.Background())
			cmd.SetVal(tt.want)
			cmd.SetErr(tt.wantErr)

			client.On("Get", mock.Anything, mock.Anything).Return(cmd).Once()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			repository := auth.New(auth.RepositoryParams{
				Pool:   pool,
				Client: client,
			})
			require.NotNil(t, repository)

			token, err := repository.GetToken(ctx, tt.arg)
			assert.Equal(t, token, tt.want)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
