package health_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	healthHandler "prodigo/internal/auth/rest/handlers/health"
	healthService "prodigo/internal/auth/usecases/health"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Check(t *testing.T) {
	tests := []struct {
		name     string
		wantCode int
		wantBody string
		wantErr  error
	}{
		{
			name:     "success",
			wantCode: http.StatusOK,
			wantBody: `{"message":"OK"}`,
			wantErr:  nil,
		},
		{
			name:     "internal server error",
			wantCode: http.StatusInternalServerError,
			wantBody: `{"error":"some error"}`,
			wantErr:  errors.New("some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(healthService.MockService)
			require.NotNil(t, service)
			defer service.AssertExpectations(t)

			service.On("Check", mock.Anything).Return(tt.wantErr)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

			handler := healthHandler.New(service)
			handler.Check(ctx)

			assert.Equal(t, tt.wantCode, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}
