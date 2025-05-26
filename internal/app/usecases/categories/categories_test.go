package categories

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"prodigo/internal/app/models"
	"prodigo/internal/app/repository/categories"
	"testing"
)

func TestService_CreateCategory(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		mockRepo := new(categories.MockRepo)
		service := &Service{repository: mockRepo}

		category := &models.Category{
			Name: "test",
		}
		mockRepo.On("CreateCategory", mock.Anything, category).Return(nil)

		err := service.CreateCategory(nil, category)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)

	})
	t.Run("error from repository", func(t *testing.T) {
		mockRepo := new(categories.MockRepo)
		service := &Service{repository: mockRepo}
		category := &models.Category{
			Name: "test",
		}
		mockRepo.On("CreateCategory", mock.Anything, category).Return(errors.New("db error"))
		err := service.CreateCategory(nil, category)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create category")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_UpdateCategory(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		mockRepo := new(categories.MockRepo)
		service := &Service{repository: mockRepo}
		category := &models.Category{
			Name: "test",
		}
		mockRepo.On("UpdateCategory", mock.Anything, category).Return(nil)
		err := service.UpdateCategory(nil, category)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error from repository", func(t *testing.T) {
		mockRepo := new(categories.MockRepo)
		service := &Service{repository: mockRepo}
		category := &models.Category{
			Name: "test",
		}
		mockRepo.On("UpdateCategory", mock.Anything, category).Return(errors.New("db error"))
		err := service.UpdateCategory(nil, category)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to update category")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_DeleteCategory(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		mockRepo := new(categories.MockRepo)
		service := &Service{repository: mockRepo}
		id := int64(1)
		mockRepo.On("DeleteCategory", mock.Anything, id).Return(nil)
		err := service.DeleteCategory(nil, id)
		assert.NoError(t, err)
	})
	t.Run("error from repository", func(t *testing.T) {
		mockRepo := new(categories.MockRepo)
		service := &Service{repository: mockRepo}
		id := int64(1)
		mockRepo.On("DeleteCategory", mock.Anything, id).Return(errors.New("db error"))
		err := service.DeleteCategory(nil, id)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to delete category")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_GetAllCategories(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		mockRepo := new(categories.MockRepo)
		service := &Service{repository: mockRepo}
		mockRepo.On("GetAllCategories", mock.Anything).Return([]*models.Category{}, nil)
		categs, err := service.GetAllCategories(nil)
		assert.NoError(t, err)
		assert.Len(t, categs, 0)
	})
	t.Run("error from repository", func(t *testing.T) {
		mockRepo := new(categories.MockRepo)
		service := &Service{repository: mockRepo}
		mockRepo.On("GetAllCategories", mock.Anything).Return([]*models.Category{}, errors.New("db error"))
		categs, err := service.GetAllCategories(nil)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to get all categories")
		assert.Nil(t, categs)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_CategoryStatistics(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(categories.MockRepo)
		service := New(mockRepo)

		expectedStats := []*models.CategoryStats{
			{
				CategoryID:    1,
				CategoryName:  "Test Category",
				ProductCount:  1,
				TotalQuantity: 10,
				TotalValue:    1000,
			},
		}

		mockRepo.On("CategoryStatistics", mock.Anything).Return(expectedStats, nil).Once()
		stats, err := service.CategoryStatistics(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedStats, stats)
		mockRepo.AssertExpectations(t)
	})
	t.Run("error from repository", func(t *testing.T) {
		mockRepo := new(categories.MockRepo)
		service := New(mockRepo)
		mockRepo.On("CategoryStatistics", mock.Anything).Return([]*models.CategoryStats{}, errors.New("db error")).Once()
		stats, err := service.CategoryStatistics(context.Background())
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to get category statistics")
		assert.Nil(t, stats)
		mockRepo.AssertExpectations(t)
	})

}
