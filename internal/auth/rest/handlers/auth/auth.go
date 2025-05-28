package auth

import (
	"errors"
	"net/http"
	"prodigo/internal/auth/dto"
	"prodigo/internal/auth/usecases/auth"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service auth.Service
}

func New(service auth.Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with username and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RegisterRequest	true	"User registration details"
//	@Failure		400		{object}	dto.Error
//	@Failure		500		{object}	dto.Error
//	@Success		201		{object}	dto.Response
//	@Router			/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: err.Error()})
		return
	}

	if err := h.service.Register(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{Message: "OK"})
}

// Login godoc
//
//	@Summary		Login a user
//	@Description	Login a user with username and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"User login details"
//	@Failure		400		{object}	dto.Error
//	@Failure		401		{object}	dto.Error
//	@Failure		404		{object}	dto.Error
//	@Failure		500		{object}	dto.Error
//	@Success		200		{object}	dto.LoginResponse
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, dto.Error{Error: err.Error()})
			return
		}
		if errors.Is(err, auth.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, dto.Error{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Refresh godoc
//
//	@Summary		Refresh access token
//	@Description	Refresh the access token using a valid refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.RefreshRequest	true	"Refresh token request details"
//	@Failure		400		{object}	dto.Error
//	@Failure		401		{object}	dto.Error
//	@Failure		404		{object}	dto.Error
//	@Failure		500		{object}	dto.Error
//	@Success		200		{object}	dto.RefreshResponse
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Error{Error: err.Error()})
		return
	}

	accessToken, err := h.service.Refresh(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, auth.ErrTokenNotFound) {
			c.JSON(http.StatusNotFound, dto.Error{Error: err.Error()})
			return
		}
		if errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrExpiredToken) {
			c.JSON(http.StatusUnauthorized, dto.Error{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.RefreshResponse{AccessToken: accessToken})
}
