package main

import (
	"prodigo/internal/auth/repository"
	"prodigo/internal/auth/rest"
	"prodigo/internal/auth/rest/handlers"
	"prodigo/internal/auth/usecases"
	"prodigo/pkg/config"
	"prodigo/pkg/db"
	"prodigo/pkg/jwt"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		db.AuthModule,
		repository.Module,
		usecases.Module,
		handlers.Module,
		rest.Module,
		jwt.Module,
	).Run()
}
