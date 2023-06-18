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
		// what origin that allowed to call
		AllowedOrigins: []string{"https://*", "http://*"},
		// method allowed
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPSTIONS"},
		// allowed header request
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "x-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		// allow cridentials ? cookies, HTTP authentication or client side SSL certificates.
		AllowCredentials: true,
		// MaxAge indicates how long (in seconds) the results of a preflight request
		MaxAge: 300,
	}))
	// prebuild middle ware to check / ping services
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/", app.Broker)
	return mux
}
