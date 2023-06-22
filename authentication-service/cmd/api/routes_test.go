package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_routes_exist(t *testing.T) {
	testApp := Config{}

	testRoutes := testApp.routes()
	chiRoutes := testRoutes.(chi.Router)

	// populate whats routes insde?
	routes := []string{"/authenticate"}
	for _, route := range routes {
		routeExists(t, chiRoutes, route)
	}
}

func routeExists(t *testing.T, routes chi.Router, route string) {
	found := false
	// walk func to search route inside chi.router
	walkFunc := func(method string, foundRoute string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == foundRoute {
			found = true
		}
		return nil
	}

	_ = chi.Walk(routes, walkFunc)
	// do check bool
	if !found {
		t.Errorf("did not find %s in registered routes", route)
	}
}
