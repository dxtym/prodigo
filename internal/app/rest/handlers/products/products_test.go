package products

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"prodigo/internal/app/models"
	"prodigo/internal/app/usecases/products"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	service := new(products.MockService)
	handler := New(service)
	assert.NotNil(t, handler)
	assert.Equal(t, service, handler.service)
}

func TestHandler_CreateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		service.On("CreateProduct", mock.Anything, mock.Anything).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"title":"Test Product","price":100,"quantity":10,"status":"available"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProduct(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), `"product created"`)

	})
	t.Run("error from service", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		service.On("CreateProduct", mock.Anything, mock.Anything).Return(errors.New("fail"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"title":"Test Product","price":100,"quantity":10,"status":"available"}`
		c.Request = httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProduct(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error"`)
	})
	t.Run("invalid body", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(""))
		c.Request.Header.Set("Content-Type", "application/json")
		handler.CreateProduct(c)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error"`)
	})
}

func TestHandler_UpdateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		updated := &models.Product{
			ID:       1,
			Title:    "Updated",
			Price:    100,
			Quantity: 5,
			Status:   "available",
		}

		service.On("UpdateProduct", mock.Anything, mock.Anything).Return(nil)
		service.On("GetProduct", mock.Anything, int64(1)).Return(updated, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"title":"Updated","price":100,"quantity":5,"status":"available"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.UpdateProduct(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"title":"Updated"`)
	})
	t.Run("invalid id", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest(http.MethodPut, "/products/abc", nil)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.UpdateProduct(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"invalid id"`)
	})
	t.Run("get product not found", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}

		service.On("UpdateProduct", mock.Anything, mock.Anything).Return(nil)
		service.On("GetProduct", mock.Anything, int64(1)).Return((*models.Product)(nil), errors.New("not found"))
		defer service.AssertExpectations(t)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := `{"title":"Updated","price":100,"quantity":5,"status":"available"}`
		c.Request = httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.UpdateProduct(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandler_GetAllProducts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		expected := []*models.Product{
			{ID: 1, Title: "Phone", Price: 1000, Quantity: 10, Status: "available"},
		}

		service.On("GetAllProducts", mock.Anything, mock.Anything).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/products", nil)

		handler.GetAllProducts(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"title":"Phone"`)
	})
	t.Run("service error", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		service.On("GetAllProducts", mock.Anything, mock.Anything).Return([]*models.Product{}, errors.New("fail"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/products", nil)

		handler.GetAllProducts(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"fail"`)
	})
}

func TestHandler_GetProductByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		expected := &models.Product{
			ID:       1,
			Title:    "Laptop",
			Price:    1500,
			Quantity: 5,
			Status:   "available",
		}

		service.On("GetProduct", mock.Anything, int64(1)).Return(expected, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodGet, "/products/1", nil)

		handler.GetProductByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"title":"Laptop"`)
	})
	t.Run("invalid id", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest(http.MethodGet, "/products/abc", nil)

		handler.GetProductByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"invalid id"`)
	})
	t.Run("product not found", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		service.On("GetProduct", mock.Anything, int64(99)).Return((*models.Product)(nil), products.ErrNotFound)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "99"}}
		c.Request = httptest.NewRequest(http.MethodGet, "/products/99", nil)

		handler.GetProductByID(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), products.ErrNotFound.Error())
	})
}

func TestHandler_DeleteProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		c.Request = req

		c.Params = gin.Params{{Key: "id", Value: "1"}}
		service.On("DeleteProduct", req.Context(), int64(1)).Return(nil)

		handler.DeleteProduct(c)
		assert.Equal(t, http.StatusNoContent, w.Code)
		service.AssertExpectations(t)
	})
	t.Run("invalid id", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/products/abc", nil)

		handler.DeleteProduct(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"invalid id"`)
	})
	t.Run("delete error", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		service.On("DeleteProduct", mock.Anything, int64(2)).Return(errors.New("delete error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest(http.MethodDelete, "/products/2", nil)

		handler.DeleteProduct(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"delete error"`)
	})
}

func TestHandler_UpdateProductStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		service.On("UpdateProductStatus", mock.Anything, int64(1), "available").Return(nil)

		body := `{"status": "available"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/products/1/status", strings.NewReader(body))

		handler.UpdateProductStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"status updated"}`, w.Body.String())
	})
	t.Run("invalid id", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/products/abc/status", nil)

		handler.UpdateProductStatus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"invalid id"`)
	})
	t.Run("not found", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		service.On("UpdateProductStatus", mock.Anything, int64(2), "archived").Return(products.ErrNotFound)

		body := `{"status": "archived"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/products/2/status", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateProductStatus(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"`+products.ErrNotFound.Error()+`"`)
	})
	t.Run("invalid status", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		body := `{"status": ""}`
		c.Request = httptest.NewRequest(http.MethodPatch, "/products/1/status", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateProductStatus(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"invalid status"`)
	})
	t.Run("internal error", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		service.On("UpdateProductStatus", mock.Anything, int64(3), "sold").Return(errors.New("update error"))

		body := `{"status": "sold"}`
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "3"}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/products/3/status", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.UpdateProductStatus(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"update error"`)
	})
}

func TestHandler_RestoreProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		service.On("RestoreProduct", mock.Anything, int64(1)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "1"}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/products/1/restore", nil)

		handler.RestoreProduct(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `{"message":"product restored"}`)
	})
	t.Run("invalid id", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/products/abc/restore", nil)

		handler.RestoreProduct(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"invalid id"`)
	})
	t.Run("restore error", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}
		defer service.AssertExpectations(t)

		service.On("RestoreProduct", mock.Anything, int64(2)).Return(errors.New("restore failed"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "2"}}
		c.Request = httptest.NewRequest(http.MethodPatch, "/products/2/restore", nil)

		handler.RestoreProduct(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error":"restore failed"`)
	})
}

func TestHandler_UploadProductImage(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		service := new(products.MockService)
		handler := &Handler{service: service}

		req := httptest.NewRequest(http.MethodPost, "/products/abc/upload", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request = req

		handler.UploadProductImage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), `"invalid product id"`)
	})
}

func TestHandler_GetProductImage(t *testing.T) {
	t.Run("invalid id", func(t *testing.T) {
		handler := &Handler{}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest(http.MethodGet, "/products/abc/image", nil)

		handler.GetProductImage(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid product id")
	})

	t.Run("image not found", func(t *testing.T) {
		handler := &Handler{}

		id := "9999"
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: id}}
		c.Request = httptest.NewRequest(http.MethodGet, "/products/9999/image", nil)

		handler.GetProductImage(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "image not found")
	})
}
