package categories

import (
	"net/http"
	"prodigo/internal/app/models"
	"prodigo/internal/app/usecases/categories"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service categories.ServiceInterface
}

func New(service categories.ServiceInterface) *Handler {
	return &Handler{service: service}
}

// CreateCategory godoc
//
//	@Summary		Create a new category
//	@Description	Create a new category with the provided details
//	@Tags			categories
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.Category	true	"Category details"
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Success		201		{object}	models.Category
//	@Router			/categories/ [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateCategory(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

// GetAllCategories godoc
//
//	@Summary		Get all categories
//	@Description	Get a list of all categories
//	@Tags			categories
//
// @Security	ApiKeyAuth
// @Accept			json
//
//	@Produce		json
//	@Failure		500	{object}	map[string]string
//	@Success		200	{object}	[]models.Category
//	@Router			/categories/ [get]
func (h *Handler) GetAllCategories(c *gin.Context) {
	cats, err := h.service.GetAllCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// UpdateCategory godoc
//
//	@Summary		Update an existing category
//	@Description	Update the details of an existing category by ID
//	@Tags			categories
//
// @Security	ApiKeyAuth
// @Accept			json
//
//	@Produce		json
//	@Param			id	path		int64	true	"Category ID"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Success		200	{object}	models.Category
//	@Router			/categories/{id} [put]
func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var cat models.Category
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat.ID = id

	if err := h.service.UpdateCategory(c.Request.Context(), &cat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// DeleteCategory godoc
//
//	@Summary		Delete a category
//	@Description	Delete an existing category by ID
//	@Tags			categories
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Category ID"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Success		204	{object}	map[string]string
//	@Router			/categories/{id} [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteCategory(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "category deleted"})
}

// CategoryStatistics godoc
//
//	@Summary		Get category statistics
//	@Description	Get statistics for categories (product count, total quantity and value)
//	@Tags			categories
//
// @Security	ApiKeyAuth
//
//	@Accept			json
//	@Produce		json
//	@Failure		500	{object}	map[string]string
//	@Success		204	{object}	[]models.CategoryStats
//	@Router			/categories/stats [get]
func (h *Handler) CategoryStatistics(c *gin.Context) {
	stats, err := h.service.CategoryStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
