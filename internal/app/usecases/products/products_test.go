package products

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"prodigo/internal/app/models"
	"prodigo/internal/app/repository/products"
	"testing"
)

func TestService_CreateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}

		product := &models.Product{
			Title: "Test Product",
		}

		mockRepo.On("CreateProduct", mock.Anything, product).Return(nil).Once()

		err := service.CreateProduct(context.Background(), product)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})
	t.Run("error from repository", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}

		product := &models.Product{
			Title: "Test Product",
		}

		mockRepo.On("CreateProduct", mock.Anything, product).Return(errors.New("db error")).Once()

		err := service.CreateProduct(context.Background(), product)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create product")

		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetAllProducts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}
		fs := &models.ProductFilterSearch{}

		mockRepo.On("GetAllProducts", mock.Anything, fs).Return([]*models.Product{}, nil).Once()
		prods, err := service.GetAllProducts(context.Background(), fs)
		assert.NoError(t, err)
		assert.Len(t, prods, 0)
	})
	t.Run("error from repository", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}
		fs := &models.ProductFilterSearch{}
		mockRepo.On("GetAllProducts", mock.Anything, fs).Return([]*models.Product{}, errors.New("db error")).Once()
		prods, err := service.GetAllProducts(context.Background(), fs)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to get all products")
		assert.Nil(t, prods)
	})
}

func TestService_GetProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}
		id := int64(1)
		mockRepo.On("GetProductByID", mock.Anything, id).Return(&models.Product{}, nil).Once()
		prod, err := service.GetProduct(context.Background(), id)
		assert.NoError(t, err)
		assert.NotNil(t, prod)
	})
	t.Run("error from repository", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}
		id := int64(1)
		mockRepo.On("GetProductByID", mock.Anything, id).Return(&models.Product{}, errors.New("db error")).Once()
		prod, err := service.GetProduct(context.Background(), id)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to get product")
		assert.Nil(t, prod)
	})
}

func TestService_UpdateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}
		existProduct := &models.Product{
			ID:    1,
			Title: "exist Product",
		}
		updatedProduct := &models.Product{
			ID:    1,
			Title: "updated Product",
		}
		mockRepo.On("GetProductByID", mock.Anything, int64(1)).Return(existProduct, nil).Once()

		mockRepo.On("UpdateProduct", mock.Anything, updatedProduct).Return(nil).Once()
		err := service.UpdateProduct(context.Background(), updatedProduct)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNumberOfCalls(t, "UpdateProduct", 1)
		mockRepo.AssertCalled(t, "UpdateProduct", mock.Anything, updatedProduct)
	})
	t.Run("error from repository", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}
		existProduct := &models.Product{
			ID:    1,
			Title: "exist Product",
		}
		updatedProduct := &models.Product{
			ID:    1,
			Title: "updated Product",
		}
		mockRepo.On("GetProductByID", mock.Anything, int64(1)).Return(existProduct, nil).Once()
		mockRepo.On("UpdateProduct", mock.Anything, updatedProduct).Return(errors.New("db error")).Once()
		err := service.UpdateProduct(context.Background(), updatedProduct)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to update product")
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNumberOfCalls(t, "UpdateProduct", 1)
		mockRepo.AssertCalled(t, "UpdateProduct", mock.Anything, updatedProduct)
	})
}

func TestService_DeleteProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}
		id := int64(1)
		mockRepo.On("DeleteProduct", mock.Anything, id).Return(nil).Once()
		err := service.DeleteProduct(context.Background(), id)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNumberOfCalls(t, "DeleteProduct", 1)
		mockRepo.AssertCalled(t, "DeleteProduct", mock.Anything, id)
	})
	t.Run("error from repository", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}
		id := int64(1)
		mockRepo.On("DeleteProduct", mock.Anything, id).Return(errors.New("db error")).Once()
		err := service.DeleteProduct(context.Background(), id)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to delete product")
	})
}

func TestService_UpdateProductStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}
		existProduct := &models.Product{
			ID:     1,
			Status: "available",
		}

		updatedProduct := &models.Product{
			ID:     1,
			Status: "not available",
		}

		mockRepo.On("GetProductByID", mock.Anything, int64(1)).Return(existProduct, nil).Once()
		mockRepo.On("UpdateProduct", mock.Anything, mock.MatchedBy(func(p *models.Product) bool {
			return updatedProduct.ID == 1 && updatedProduct.Status == "not available"
		})).Return(nil).Once()

		err := service.UpdateProductStatus(context.Background(), 1, "not available")
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)

	})
	t.Run("product_not_found", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := &Service{repository: mockRepo}

		mockRepo.On("GetProductByID", mock.Anything, int64(1)).Return(&models.Product{}, products.ErrNotFound).Once()

		err := service.UpdateProductStatus(context.Background(), 1, "not available")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product not found")

		mockRepo.AssertExpectations(t)
	})

}

func TestService_RestoreProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := Service{repository: mockRepo}
		mockRepo.On("RestoreProduct", mock.Anything, int64(1)).Return(nil).Once()
		err := service.RestoreProduct(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("restore failure", func(t *testing.T) {
		mockRepo := new(products.MockRepo)
		service := Service{repository: mockRepo}
		mockRepo.On("RestoreProduct", mock.Anything, int64(1)).Return(products.ErrNotFound).Once()
		err := service.RestoreProduct(context.Background(), 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to restore product")
	})
}
