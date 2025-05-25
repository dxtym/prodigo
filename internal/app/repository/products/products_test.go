package products

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
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
			Return(pgconn.CommandTag{}, errors.New("ddd"))
		pool, err := New(repo, nil)
		assert.NotNil(t, err)
		assert.Nil(t, pool)
	})
	t.Run("can run ddl", func(t *testing.T) {
		repo := new(postgres.MockPool)
		defer repo.AssertExpectations(t)

		repo.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("INSERT 1"), nil)
		pool, err := New(repo, nil)
		assert.Nil(t, err)
		assert.NotNil(t, pool)
	})
}

func Test_repository_CreateProduct(t *testing.T) {
	t.Run("error on insert", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		repo := &repository{pool: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.CommandTag{}, errors.New("cannot create product"))

		err := repo.CreateProduct(context.Background(), &models.Product{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create product")

		mockPool.AssertExpectations(t)
	})
	t.Run("success on insert product", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		defer mockPool.AssertExpectations(t)

		repo := &repository{pool: mockPool}

		mockPool.On("Exec", mock.Anything, mock.Anything, mock.Anything).
			Return(pgconn.NewCommandTag("INSERT 1"), nil)

		err := repo.CreateProduct(context.Background(), &models.Product{})

		assert.NoError(t, err)
	})
}

func TestRepository_GetProductByID(t *testing.T) {
	t.Run("error on get product by id", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRow := new(postgres.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repository{pool: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(pgx.ErrNoRows)

		task, err := pool.GetProductByID(context.Background(), 1)
		assert.NotNil(t, err)
		assert.Nil(t, task)
	})
	t.Run("success on get product by id", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRow := new(postgres.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repository{pool: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		task, err := pool.GetProductByID(context.Background(), 1)
		assert.Nil(t, err)
		assert.NotNil(t, task)
	})
	t.Run("error on get task by id", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRow := new(postgres.MockRow)
		defer mockPool.AssertExpectations(t)
		defer mockRow.AssertExpectations(t)

		pool := &repository{pool: mockPool}

		mockPool.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).
			Return(mockRow)
		mockRow.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error"))

		task, err := pool.GetProductByID(context.Background(), 1)
		assert.NotNil(t, err)
		assert.Nil(t, task)
	})
}

func TestRepository_GetAllProducts(t *testing.T) {
	t.Run("success on get all products", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRows := new(postgres.MockRow)

		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		repo := &repository{pool: mockPool}

		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Close").Return()

		products, err := repo.GetAllProducts(context.Background(), &models.ProductFilterSearch{})
		assert.NoError(t, err)
		assert.Len(t, products, 1)
	})
	t.Run("query error", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRows := new(postgres.MockRow)

		defer mockPool.AssertExpectations(t)

		repo := &repository{pool: mockPool}

		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, errors.New("query failed"))

		products, err := repo.GetAllProducts(context.Background(), &models.ProductFilterSearch{})
		assert.Error(t, err)
		assert.Nil(t, products)
	})
	t.Run("scan error", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		mockRows := new(postgres.MockRow)

		defer mockPool.AssertExpectations(t)
		defer mockRows.AssertExpectations(t)

		repo := &repository{pool: mockPool}

		mockPool.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockRows, nil)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("scan failed")).Once()
		mockRows.On("Close").Return()

		products, err := repo.GetAllProducts(context.Background(), &models.ProductFilterSearch{})
		assert.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestRepository_UpdateProduct(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockPool := new(postgres.MockPool)

		repo := &repository{pool: mockPool}

		ctx := context.Background()
		tag := pgconn.NewCommandTag("UPDATE 1")
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(tag, nil)

		err := repo.UpdateProduct(ctx, &models.Product{ID: 1})
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockPool := new(postgres.MockPool)

		repo := &repository{pool: mockPool}

		ctx := context.Background()
		tag := pgconn.NewCommandTag("UPDATE 0")
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(tag, nil)

		err := repo.UpdateProduct(ctx, &models.Product{ID: 2})
		assert.EqualError(t, err, "product not found")
	})

	t.Run("exec error", func(t *testing.T) {
		mockPool := new(postgres.MockPool)

		repo := &repository{pool: mockPool}

		ctx := context.Background()
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, errors.New("some error"))

		err := repo.UpdateProduct(ctx, &models.Product{ID: 3})
		assert.Error(t, err)
	})
}

func TestRepository_DeleteProduct(t *testing.T) {
	t.Run("success delete", func(t *testing.T) {
		mockPool := new(postgres.MockPool)

		repo := &repository{pool: mockPool}
		ctx := context.Background()
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(pgconn.NewCommandTag("DELETE 1"), nil)

		err := repo.DeleteProduct(ctx, 1)
		assert.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		repo := &repository{pool: mockPool}
		ctx := context.Background()
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(pgconn.NewCommandTag("DELETE 0"), nil)
		err := repo.DeleteProduct(ctx, 2)
		assert.EqualError(t, err, "product not found")
	})
	t.Run("exec error", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		repo := &repository{pool: mockPool}
		ctx := context.Background()
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, errors.New("some error"))
		err := repo.DeleteProduct(ctx, 3)
		assert.Error(t, err)
	})
}

func Test_repository_RestoreProduct(t *testing.T) {
	t.Run("success restore", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		repo := &repository{pool: mockPool}
		ctx := context.Background()
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(pgconn.NewCommandTag("UPDATE 1"), nil)
		err := repo.RestoreProduct(ctx, 1)
		assert.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		repo := &repository{pool: mockPool}
		ctx := context.Background()
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(pgconn.NewCommandTag("UPDATE 0"), nil)
		err := repo.RestoreProduct(ctx, 2)
		assert.EqualError(t, err, "product not found")
	})
	t.Run("exec error", func(t *testing.T) {
		mockPool := new(postgres.MockPool)
		repo := &repository{pool: mockPool}
		ctx := context.Background()
		mockPool.On("Exec", ctx, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, errors.New("some error"))
		err := repo.RestoreProduct(ctx, 3)
		assert.Error(t, err)
	})
}
