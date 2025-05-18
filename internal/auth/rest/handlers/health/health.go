package health

import (
	"net/http"

	"prodigo/internal/auth/usecases/health"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *health.Service
}

func New(service *health.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Check(c *gin.Context) {
	if err := h.service.Check(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
