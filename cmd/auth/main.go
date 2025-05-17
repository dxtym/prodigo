package main

import (
	"prodigo/internal/auth/repository"
	"prodigo/internal/auth/rest"
	"prodigo/internal/auth/rest/handlers"
	"prodigo/internal/auth/usecases"
	"prodigo/pkg/config"
	"prodigo/pkg/db"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		db.Module,
		repository.Module,
		usecases.Module,
		handlers.Module,
		rest.Module,
		fx.Provide(gin.New),
		fx.Invoke(func(srv *rest.Server, cfg *config.AuthConfig) error {
			return srv.Start(cfg.Host, cfg.Port)
		}),
	).Run()
}
