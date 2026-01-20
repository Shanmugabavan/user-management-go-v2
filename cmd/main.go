package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"user-management/api/route"
	"user-management/bootstrap"
	"user-management/internal/validator"
)

func main() {
	app := bootstrap.Application{}
	app.App()

	validator.Init()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	r.Route("/user/", func(r chi.Router) {})

	route.Setup(app.Env, app.ConnectionPool, r)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
