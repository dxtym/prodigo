package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"prodigo/internal/app/rest/handlers/categories"
	"prodigo/internal/app/rest/handlers/products"
	"prodigo/internal/app/rest/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	mux             *gin.Engine
	mw              *middleware.Middleware
	srv             *http.Server
	categoryHandler *categories.Handler
	productHandler  *products.Handler
}

func New(
	mw *middleware.Middleware,
	productHandler *products.Handler,
	categoryHandler *categories.Handler,
) *Server {
	return &Server{
		mux:             gin.New(),
		mw:              mw,
		productHandler:  productHandler,
		categoryHandler: categoryHandler,
	}
}

func (s *Server) SetupRoutes() {
	s.mux.Use(gin.Logger())
	s.mux.Use(gin.Recovery())

	v1 := s.mux.Group("/api/v1")
	{
		v1.Use(s.mw.Auth())

		prods := v1.Group("/products")
		{
			prods.POST("/", s.productHandler.CreateProduct)
			prods.GET("/", s.productHandler.GetAllProducts)
			prods.GET("/:id", s.productHandler.GetProductByID)
			prods.PUT("/:id", s.productHandler.UpdateProduct)
			prods.DELETE("/:id", s.productHandler.DeleteProduct)
			prods.PUT("/:id/status", s.productHandler.UpdateProductStatus)
			prods.POST("/:id/image", s.productHandler.UploadProductImage)
			prods.GET("/:id/image", s.productHandler.GetProductImage)
		}

		cats := v1.Group("/categories")
		{
			cats.POST("/", s.categoryHandler.CreateCategory)
			cats.GET("/", s.categoryHandler.GetAllCategories)
			cats.PUT("/:id", s.categoryHandler.UpdateCategory)
			cats.DELETE("/:id", s.categoryHandler.DeleteCategory)
		}
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
