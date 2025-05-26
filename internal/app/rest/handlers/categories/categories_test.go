package categories

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"prodigo/internal/app/models"
	"prodigo/internal/app/usecases/categories"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	service := new(categories.MockService)
	defer service.AssertExpectations(t)
	handler := New(service)
	assert.NotNil(t, handler)
}

func TestHandler_CreateCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)
		defer service.AssertExpectations(t)

		service.On("CreateCategory", mock.Anything, mock.Anything).Return(nil)

		body := `{"name":"Electronics"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))

		handler.CreateCategory(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"name":"Electronics"`)
	})
	t.Run("invalid json", func(t *testing.T) {
		handler := New(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader("not-json"))

		handler.CreateCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error"`)
	})
	t.Run("service error", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)
		defer service.AssertExpectations(t)

		service.On("CreateCategory", mock.Anything, mock.Anything).Return(errors.New("db error"))

		body := `{"name":"FailCat"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateCategory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error"`)
	})
}

func TestHandler_GetAllCategories(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)
		defer service.AssertExpectations(t)

		expected := []*models.Category{
			{ID: 1, Name: "Electronics"},
			{ID: 2, Name: "Books"},
		}
		service.On("GetAllCategories", mock.Anything).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/categories", nil)

		handler.GetAllCategories(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Electronics")
		assert.Contains(t, w.Body.String(), "Books")
	})
	t.Run("service error", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)
		defer service.AssertExpectations(t)

		service.On("GetAllCategories", mock.Anything).Return([]*models.Category{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/categories", nil)

		handler.GetAllCategories(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error"`)
	})
}

func TestHandler_UpdateCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)
		defer service.AssertExpectations(t)

		category := models.Category{ID: 1, Name: "Updated Name"}
		service.On("UpdateCategory", mock.Anything, &category).Return(nil)

		body := `{"name": "Updated Name"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/categories/1", strings.NewReader(body))

		handler.UpdateCategory(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Updated Name")
	})
	t.Run("invalid id", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/categories/abc", nil)

		handler.UpdateCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid id")
	})
	t.Run("service error", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)
		defer service.AssertExpectations(t)

		category := models.Category{ID: 1, Name: "New Name"}
		service.On("UpdateCategory", mock.Anything, &category).Return(errors.New("update failed"))

		body := `{"name": "New Name"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/categories/1", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateCategory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error"`)
	})
	t.Run("bad json", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/categories/1", strings.NewReader("{bad json"))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error"`)
	})
}

func TestHandler_DeleteCategory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)
		defer service.AssertExpectations(t)

		service.On("DeleteCategory", mock.Anything, int64(1)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/categories/1", nil)

		handler.DeleteCategory(c)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
	})
	t.Run("invalid id", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/categories/abc", nil)

		handler.DeleteCategory(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid id")
	})

	t.Run("service error", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)
		defer service.AssertExpectations(t)

		service.On("DeleteCategory", mock.Anything, int64(1)).Return(errors.New("delete failed"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/categories/1", nil)

		handler.DeleteCategory(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "delete failed")
	})
}

func TestHandler_CategoryStatistics(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)
		defer service.AssertExpectations(t)

		service.On("CategoryStatistics", mock.Anything).Return([]*models.CategoryStats{
			{CategoryID: 1, CategoryName: "SUV", ProductCount: 5},
		}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/categories/stats", nil)

		handler.GetCategoryStatistics(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "SUV")
	})
	t.Run("internal error", func(t *testing.T) {
		service := new(categories.MockService)
		handler := New(service)
		defer service.AssertExpectations(t)

		service.On("CategoryStatistics", mock.Anything).Return([]*models.CategoryStats{}, errors.New("db error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/categories/stats", nil)

		handler.GetCategoryStatistics(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "db error")
	})
}
