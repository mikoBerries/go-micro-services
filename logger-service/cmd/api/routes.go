package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()
	// add middleware before routing to specific microservices
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPSTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "x-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	// prebuild middle ware to check / ping services
	mux.Use(middleware.Heartbeat("/ping"))

	// mux.Post("/", app.Broker)
	// mux.Post("/handle", app.HandleSubmission)
	return mux
}
