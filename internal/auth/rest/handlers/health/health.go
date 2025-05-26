package health

import (
	"net/http"

	"prodigo/internal/auth/dto"
	"prodigo/internal/auth/usecases/health"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service health.Service
}

func New(service health.Service) *Handler {
	return &Handler{service: service}
}

// Check godoc
//
//	@Summary		Check service health
//	@Description	Ping the database and cache to verify service health
//	@Tags			health
//	@Produce		json
//	@Failure		500	{object}	dto.Error
//	@Success		200	{object}	dto.Response
//	@Router			/health [get]
func (h *Handler) Check(c *gin.Context) {
	if err := h.service.Check(c.Request.Context()); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{Message: "OK"})
}
