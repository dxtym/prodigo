package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"prodigo/internal/app/repository"
	"prodigo/internal/app/rest"
	"prodigo/internal/app/rest/handlers"
	"prodigo/internal/app/usecases"
	"prodigo/pkg/config"
	"prodigo/pkg/db"
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
		fx.Invoke(func(srv *rest.Server, cfg *config.Config) error { return srv.Start(cfg.AuthHost, cfg.AuthPort) }),
	).Run()
}
