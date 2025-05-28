package products

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"prodigo/internal/app/models"
	"prodigo/internal/app/usecases/products"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service products.ServiceInterface
}

func New(service products.ServiceInterface) *Handler {
	return &Handler{service: service}
}

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	Create a new product with the provided details
//	@Tags			products
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.Product	true	"Product details"
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Success		201		{object}	map[string]string
//	@Router			/products/ [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateProduct(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message":  "product created",
		"title":    p.Title,
		"price":    p.Price,
		"quantity": p.Quantity,
		"status":   p.Status,
	})
}

// GetAllProducts godoc
//
//	@Summary		Get all products
//	@Description	Get a list of all products with optional filters
//	@Tags			products
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			category	query		string	false	"Filter by category name"
//	@Param			status	query		string	false	"Filter by product status"
//	@Param			price_min	query		int	false	"Minimum price filter"
//	@Param			price_max	query		int	false	"Maximum price filter"
//	@Param			search	query		string	false	"Search term for product title"
//	@Failure		500		{object}	map[string]string
//	@Success		200		{object}	[]models.Product
//	@Router			/products/ [get]
func (h *Handler) GetAllProducts(c *gin.Context) {
	var fs models.ProductFilterSearch
	if v := c.Query("category"); v != "" {
		fs.CategoryName = v
	}
	if v := c.Query("status"); v != "" {
		fs.Status = v
	}
	if v := c.Query("price_min"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			fs.PriceMin = n
		}
	}
	if v := c.Query("price_max"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			fs.PriceMax = n
		}
	}
	if v := c.Query("search"); v != "" {
		fs.Search = v
	}
	prods, err := h.service.GetAllProducts(c.Request.Context(), &fs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prods)
}

// GetProductByID godoc
//
//	@Summary		Get a product by ID
//	@Description	Get the details of a product by its ID
//	@Tags			products
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Product ID"
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Success		200	{object}	models.Product
//	@Router			/products/{id} [get]
func (h *Handler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	product, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, products.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
//
//	@Summary		Update an existing product
//	@Description	Update the details of an existing product by ID
//	@Tags			products
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int64			true	"Product ID"
//	@Param			request	body		models.Product	true	"Product details"
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Success		200		{object}	models.Product
//	@Router			/products/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var p models.Product
	if err = c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = id

	if err = h.service.UpdateProduct(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updatedProduct, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, updatedProduct)

}

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Delete an existing product by ID
//	@Tags			products
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Product ID"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Success		204	{object}	map[string]string
//	@Router			/products/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "product deleted"})
}

type UpdateStatus struct {
	Status string `json:"status"`
}

// UpdateProductStatus godoc
//
//	@Summary		Update product status
//	@Description	Update the status of a product by ID
//	@Tags			products
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int64			true	"Product ID"
//	@Param			request	body		UpdateStatus	true	"Product status"
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Success		200		{object}	map[string]string
//	@Router			/products/{id}/status [put]
func (h *Handler) UpdateProductStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var payload UpdateStatus
	if err = c.ShouldBindJSON(&payload); err != nil || payload.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	err = h.service.UpdateProductStatus(c.Request.Context(), id, payload.Status)
	if err != nil {
		if errors.Is(err, products.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

const uploadDir = "./uploads/products"

// UploadProductImage godoc
//
//	@Summary		Upload product image
//	@Description	Upload an image for a product by ID
//	@Tags			products
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		int64	true	"Product ID"
//	@Param			image	formData	file	true	"Product image file"
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Success		200		{object}	map[string]string
//	@Router			/products/{id}/image [post]
func (h *Handler) UploadProductImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}

	defer func() {
		if err = file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file header"})
		return
	}

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type: " + contentType})
		return
	}
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 10MB"})
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rewind file"})
		return
	}

	productDir := filepath.Join(uploadDir, idStr)
	const perm = 0o750
	if err = os.MkdirAll(productDir, perm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create directory"})
		return
	}

	filePath := filepath.Join(productDir, "image.jpg")
	// #nosec G304
	dst, errf := os.Create(filePath)
	if errf != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create file"})
		return
	}

	defer func() {
		if err = dst.Close(); err != nil {
			fmt.Println("Error closing file dst:", err)
		}
	}()

	if _, err = io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	product.Image = filePath
	if err := h.service.UpdateProduct(c.Request.Context(), product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update image in DB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "image uploaded", "filename": header.Filename})
}

// GetProductImage godoc
//
//	@Summary		Get product image
//	@Description	Get the image of a product by ID
//	@Tags			products
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		image/jpeg
//	@Param			id	path		int64	true	"Product ID"
//	@Success		200	{string}	binary	"Image binary data"
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/products/{id}/image [get]
func (h *Handler) GetProductImage(c *gin.Context) {
	id := c.Param("id")
	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	filePath := filepath.Join(uploadDir, id, "image.jpg")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
		return
	}

	c.File(filePath)
}

// RestoreProduct godoc
//
//	@Summary		Restore a deleted product
//	@Description	Restore a soft-deleted product by ID
//	@Tags			products
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Product ID"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Success		200	{object}	map[string]string
//	@Router			/products/{id}/restore [put]
func (h *Handler) RestoreProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.RestoreProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusAccepted)
	c.JSON(http.StatusOK, gin.H{"message": "product restored"})
}
