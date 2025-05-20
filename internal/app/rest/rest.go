package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"prodigo/internal/app/rest/handlers/categories"
	"prodigo/internal/app/rest/handlers/products"
	"time"
)

type Server struct {
	mux             *gin.Engine
	srv             *http.Server
	categoryHandler *categories.Handler
	productHandler  *products.Handler
}

func New(mux *gin.Engine, productHandler *products.Handler, categoryHandler *categories.Handler) *Server {
	return &Server{mux: mux, productHandler: productHandler, categoryHandler: categoryHandler}
}

func (s *Server) SetupRoutes() {
	s.mux.Use(gin.Logger())
	s.mux.Use(gin.Recovery())

	v1 := s.mux.Group("/api/v1")
	{
		v1.POST("/products", s.productHandler.CreateProduct)
		v1.GET("/products", s.productHandler.GetAllProducts)
		v1.GET("/products/:id", s.productHandler.GetProductByID)
		v1.PUT("/products/:id", s.productHandler.UpdateProduct)
		v1.DELETE("/products/:id", s.productHandler.DeleteProduct)
		v1.PUT("/products/:id/status", s.productHandler.UpdateProductStatus)
		v1.POST("/products/:id/image", s.productHandler.UploadProductImage)
		v1.GET("/products/:id/image", s.productHandler.GetProductImage)

		v1.POST("/categories", s.categoryHandler.CreateCategory)
		v1.GET("/categories", s.categoryHandler.GetAllCategories)
		v1.PUT("/categories/:id", s.categoryHandler.UpdateCategory)
		v1.DELETE("/categories/:id", s.categoryHandler.DeleteCategory)
	}
}

func (s *Server) Start(host, port string) error {
	s.SetupRoutes()

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
