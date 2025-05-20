package categories

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"prodigo/internal/app/models"
	"prodigo/pkg/db/postgres"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("can't run ddl", func(t *testing.T) {
		repo := new(postgres.MockPool)
		defer repo.AssertExpectations(t)

		repo.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("cannot create table categories"))
		pool, err := New(repo)
		assert.NotNil(t, err)
		assert.Nil(t, pool)
	})
	t.Run("can run ddl", func(t *testing.T) {
		repo := new(postgres.MockPool)
		defer repo.AssertExpectations(t)

		repo.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("INSERT 1"), nil)
		pool, err := New(repo)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
	})
}

func Test_repository_CreateCategory(t *testing.T) {
	t.Run("error on insert", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRow := new(postgres.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		category := &models.Category{}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("insert error"))

		repo := &repository{pool: mockPool}
		err := repo.CreateCategory(context.Background(), category)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insert error")

	})
	t.Run("success on insert", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRow := new(postgres.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		category := &models.Category{
			Name: "Valid Category",
		}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		repo := &repository{pool: mockPool}
		err := repo.CreateCategory(context.Background(), category)

		assert.NoError(t, err)
	})
}

func TestRepository_GetAllCategories(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRows := new(postgres.MockRow)

		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		repo := &repository{pool: mockPool}

		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return(nil).Once()

		categories, err := repo.GetAllCategories(context.Background())
		assert.NoError(t, err)
		assert.Len(t, categories, 1)
	})
	t.Run("query_error", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRows := new(postgres.MockRow)
		defer mockPool.AssertExpectations(t)

		repo := &repository{pool: mockPool}

		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, errors.New("query failed"))

		categories, err := repo.GetAllCategories(context.Background())
		assert.Error(t, err)
		assert.Nil(t, categories)
	})
	t.Run("scan_error", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRows := new(postgres.MockRow)

		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		repo := &repository{pool: mockPool}

		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("scan failed")).Once()
		mockRows.On("Close").Return(nil).Once()

		categories, err := repo.GetAllCategories(context.Background())
		assert.Error(t, err)
		assert.Nil(t, categories)
	})
}

func Test_repository_UpdateCategory(t *testing.T) {
	t.Run("not found cat", func(t *testing.T) {
		mockPool := new(postgres.MockPool)

		repo := &repository{pool: mockPool}
		defer mockPool.AssertExpectations(t)

		ctx := context.Background()
		tag := pgconn.NewCommandTag("UPDATE 0")
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(tag, nil)
		err := repo.UpdateCategory(ctx, &models.Category{ID: 1})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category not found or delete")
	})
	t.Run("success update of cat", func(t *testing.T) {
		mockPool := new(postgres.MockPool)

		repo := &repository{pool: mockPool}

		defer mockPool.AssertExpectations(t)

		ctx := context.Background()
		tag := pgconn.NewCommandTag("UPDATE 1")
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(tag, nil)

		err := repo.UpdateCategory(ctx, &models.Category{ID: 1})
		assert.NoError(t, err)
	})
	t.Run("exec error", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRow := new(postgres.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)
		repo := &repository{pool: mockPool}
		ctx := context.Background()
		tag := pgconn.NewCommandTag("UPDATE 0")
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(tag, errors.New("exec error"))
		err := repo.UpdateCategory(ctx, &models.Category{ID: 1})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exec error")
	})
}

func Test_repository_DeleteCategory(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		mockPool := new(postgres.MockPool)

		defer mockPool.AssertExpectations(t)

		pool := &repository{pool: mockPool}

		ctx := context.Background()
		tag := pgconn.NewCommandTag("DELETE 1")
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(tag, nil)
		err := pool.DeleteCategory(ctx, 1)
		assert.NoError(t, err)
	})
	t.Run("exec error", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRow := new(postgres.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)
		repo := &repository{pool: mockPool}
		ctx := context.Background()
		tag := pgconn.NewCommandTag("DELETE 0")
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(tag, errors.New("exec error"))

		err := repo.DeleteCategory(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exec error")
	})
	t.Run("not found cat or deleted", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRow := new(postgres.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)
		repo := &repository{pool: mockPool}
		ctx := context.Background()
		tag := pgconn.NewCommandTag("DELETE 0")
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(tag, nil)
		err := repo.DeleteCategory(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "category not found or delete")
	})
}
