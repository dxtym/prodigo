package products

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
	t.Run("success on create product", func(t *testing.T) {
		mockRepo := new(MockRepo)

		p := &models.Product{
			Title:      "Test Product",
			CategoryID: 1,
			Price:      100,
			Quantity:   10,
			Image:      "",
			Status:     "available",
		}

		mockRepo.On("CreateProduct", mock.Anything, p).Return(nil)

		err := mockRepo.CreateProduct(context.Background(), p)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error on create product", func(t *testing.T) {
		mockRepo := new(MockRepo)
		p := &models.Product{}

		mockRepo.On("CreateProduct", mock.Anything, p).Return(errors.New("db error"))

		err := mockRepo.CreateProduct(context.Background(), p)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})
}

func Test_repository_GetAllProducts(t *testing.T) {
	t.Run("success on get all products", func(t *testing.T) {
		mockRepo := new(MockRepo)

		expected := []*models.Product{
			{ID: 1, Title: "Product 1"},
			{ID: 2, Title: "Product 2"},
		}
		mockRepo.On("GetAllProducts", mock.Anything).Return(expected, nil)

		products, err := mockRepo.GetAllProducts(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, products)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error on get all products", func(t *testing.T) {
		mockRepo := new(MockRepo)

		mockRepo.On("GetAllProducts", mock.Anything).Return([]*models.Product(nil), errors.New("db error"))

		products, err := mockRepo.GetAllProducts(context.Background())

		assert.Error(t, err)
		assert.Nil(t, products)
		mockRepo.AssertExpectations(t)
	})
	t.Run("empty products", func(t *testing.T) {
		mockRepo := new(MockRepo)
		mockRepo.On("GetAllProducts", mock.Anything).Return([]*models.Product(nil), nil)
		products, err := mockRepo.GetAllProducts(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, []*models.Product(nil), products)
		mockRepo.AssertExpectations(t)
	})
}

func Test_repository_GetProductByID(t *testing.T) {
	t.Run("success on get product by id", func(t *testing.T) {
		mockRepo := new(MockRepo)

		expected := &models.Product{ID: 1, Title: "Product 1"}
		mockRepo.On("GetProductByID", mock.Anything, int64(1)).Return(expected, nil)

		product, err := mockRepo.GetProductByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, product)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error on get product by id", func(t *testing.T) {
		mockRepo := new(MockRepo)
		mockRepo.On("GetProductByID", mock.Anything, int64(1)).Return((*models.Product)(nil), errors.New("not found"))
		product, err := mockRepo.GetProductByID(context.Background(), int64(1))
		assert.Error(t, err)
		assert.Nil(t, product)
		mockRepo.AssertExpectations(t)
	})
	t.Run("not found product", func(t *testing.T) {
		mockRepo := new(MockRepo)
		mockRepo.On("GetProductByID", mock.Anything, int64(1)).Return((*models.Product)(nil), nil)
		product, err := mockRepo.GetProductByID(context.Background(), int64(1))
		assert.NoError(t, err)
		assert.Nil(t, product)
		mockRepo.AssertExpectations(t)
	})
}

func Test_repository_UpdateProduct(t *testing.T) {
	t.Run("success on update product", func(t *testing.T) {
		mockRepo := new(MockRepo)
		p := &models.Product{ID: 1, Title: "Product 1"}
		mockRepo.On("UpdateProduct", mock.Anything, p).Return(nil)
		err := mockRepo.UpdateProduct(context.Background(), p)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error on update product", func(t *testing.T) {
		mockRepo := new(MockRepo)
		p := &models.Product{}
		mockRepo.On("UpdateProduct", mock.Anything, p).Return(errors.New("update error"))
		err := mockRepo.UpdateProduct(context.Background(), p)
		assert.Error(t, err)
		assert.EqualError(t, err, "update error")
	})
}

func Test_repository_DeleteProduct(t *testing.T) {
	t.Run("success on delete product", func(t *testing.T) {
		mockRepo := new(MockRepo)
		mockRepo.On("DeleteProduct", mock.Anything, int64(1)).Return(nil)
		err := mockRepo.DeleteProduct(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error on delete product", func(t *testing.T) {
		mockRepo := new(MockRepo)
		mockRepo.On("DeleteProduct", mock.Anything, int64(1)).Return(errors.New("delete error"))
		err := mockRepo.DeleteProduct(context.Background(), 1)
		assert.Error(t, err)
		assert.EqualError(t, err, "delete error")
	})
}
