package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"prodigo/internal/auth/dto"
	authHandler "prodigo/internal/auth/rest/handlers/auth"
	authService "prodigo/internal/auth/usecases/auth"
	"prodigo/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Register(t *testing.T) {
	tests := []struct {
		name     string
		arg      dto.RegisterRequest
		wantCode int
		wantErr  error
	}{
		{
			name: "success",
			arg: dto.RegisterRequest{
				Username: utils.GenerateRandomString(10),
				Password: utils.GenerateRandomString(10),
			},
			wantCode: http.StatusCreated,
			wantErr:  nil,
		},
		{
			name: "bad request",
			arg: dto.RegisterRequest{
				Username: "",
				Password: "",
			},
			wantCode: http.StatusBadRequest,
			wantErr:  nil,
		},
		{
			name: "internal server error",
			arg: dto.RegisterRequest{
				Username: utils.GenerateRandomString(10),
				Password: utils.GenerateRandomString(10),
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(authService.MockService)
			require.NotNil(t, service)
			defer service.AssertExpectations(t)

			service.On("Register",
				mock.Anything,
				mock.Anything,
			).Return(tt.wantErr).Maybe()

			body, err := json.Marshal(tt.arg)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))

			handler := authHandler.New(service)
			handler.Register(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestHandler_Login(t *testing.T) {
	tests := []struct {
		name     string
		arg      dto.LoginRequest
		wantCode int
		wantErr  error
	}{
		{
			name: "success",
			arg: dto.LoginRequest{
				Username: utils.GenerateRandomString(10),
				Password: utils.GenerateRandomString(10),
			},
			wantCode: http.StatusOK,
			wantErr:  nil,
		},
		{
			name: "bad request",
			arg: dto.LoginRequest{
				Username: "",
				Password: "",
			},
			wantCode: http.StatusBadRequest,
			wantErr:  nil,
		},
		{
			name: "user not found",
			arg: dto.LoginRequest{
				Username: utils.GenerateRandomString(10),
				Password: utils.GenerateRandomString(10),
			},
			wantCode: http.StatusNotFound,
			wantErr:  authService.ErrUserNotFound,
		},
		{
			name: "invalid credentials",
			arg: dto.LoginRequest{
				Username: utils.GenerateRandomString(10),
				Password: utils.GenerateRandomString(10),
			},
			wantCode: http.StatusUnauthorized,
			wantErr:  authService.ErrInvalidCredentials,
		},
		{
			name: "internal server error",
			arg: dto.LoginRequest{
				Username: utils.GenerateRandomString(10),
				Password: utils.GenerateRandomString(10),
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(authService.MockService)
			require.NotNil(t, service)
			defer service.AssertExpectations(t)

			service.On("Login",
				mock.Anything,
				mock.Anything,
			).Return(
				utils.GenerateRandomString(10),
				utils.GenerateRandomString(10),
				tt.wantErr,
			).Maybe()

			body, err := json.Marshal(tt.arg)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))

			handler := authHandler.New(service)
			handler.Login(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestHandler_Refresh(t *testing.T) {
	tests := []struct {
		name     string
		arg      dto.RefreshRequest
		wantCode int
		wantErr  error
	}{
		{
			name: "success",
			arg: dto.RefreshRequest{
				RefreshToken: utils.GenerateRandomString(10),
			},
			wantCode: http.StatusOK,
			wantErr:  nil,
		},
		{
			name: "bad request",
			arg: dto.RefreshRequest{
				RefreshToken: "",
			},
			wantCode: http.StatusBadRequest,
			wantErr:  nil,
		},
		{
			name: "token not found",
			arg: dto.RefreshRequest{
				RefreshToken: utils.GenerateRandomString(10),
			},
			wantCode: http.StatusNotFound,
			wantErr:  authService.ErrTokenNotFound,
		},
		{
			name: "invalid token",
			arg: dto.RefreshRequest{
				RefreshToken: utils.GenerateRandomString(10),
			},
			wantCode: http.StatusUnauthorized,
			wantErr:  authService.ErrInvalidToken,
		},
		{
			name: "internal server error",
			arg: dto.RefreshRequest{
				RefreshToken: utils.GenerateRandomString(10),
			},
			wantCode: http.StatusInternalServerError,
			wantErr:  errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(authService.MockService)
			require.NotNil(t, service)
			defer service.AssertExpectations(t)

			service.On("Refresh",
				mock.Anything,
				mock.Anything,
			).Return(
				utils.GenerateRandomString(10),
				tt.wantErr,
			).Maybe()

			body, err := json.Marshal(tt.arg)
			require.NoError(t, err)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewReader(body))

			handler := authHandler.New(service)
			handler.Refresh(ctx)

			require.NoError(t, err)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
