package main

import (
	"prodigo/internal/app/repository"
	"prodigo/internal/app/rest"
	"prodigo/internal/app/rest/casbin"
	"prodigo/internal/app/rest/handlers"
	"prodigo/internal/app/rest/middleware"
	"prodigo/internal/app/usecases"
	"prodigo/pkg/config"
	"prodigo/pkg/db"
	"prodigo/pkg/jwt"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		db.AppModule,
		repository.Module,
		usecases.Module,
		handlers.Module,
		rest.Module,
		middleware.Module,
		jwt.Module,
		casbin.Module,
	).Run()
}
