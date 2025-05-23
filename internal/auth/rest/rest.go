package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"prodigo/internal/auth/rest/handlers/auth"
	"prodigo/internal/auth/rest/handlers/health"

	"github.com/gin-gonic/gin"
)

type Server struct {
	mux           *gin.Engine
	srv           *http.Server
	healthHandler *health.Handler
	authHandler   *auth.Handler
}

func New(healthHandler *health.Handler, authHandler *auth.Handler) *Server {
	return &Server{
		mux:           gin.New(),
		healthHandler: healthHandler,
		authHandler:   authHandler,
	}
}

func (s *Server) setupRoutes() {
	s.mux.Use(gin.Logger())
	s.mux.Use(gin.Recovery())

	v1 := s.mux.Group("/api/v1")
	{
		v1.GET("/health", s.healthHandler.Check)

		auths := v1.Group("/auth")
		{
			auths.POST("/register", s.authHandler.Register)
			auths.POST("/login", s.authHandler.Login)
			auths.POST("/refresh", s.authHandler.Refresh)
		}
	}
}

func (s *Server) Start(host, port string) error {
	s.setupRoutes()

	s.srv = &http.Server{
		Addr:              net.JoinHostPort(host, port),
		Handler:           s.mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := s.srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.srv == nil {
		return nil
	}

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	return nil
}
