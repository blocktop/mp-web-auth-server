package main

import (
	"github.com/blocktop/mp-common/server/middleware"
	"github.com/go-chi/chi"
)

func setRoutes(r *chi.Mux) {
	r.Get("/health", middleware.HealthHandler)

	r.Get("/token", handleGetToken)
	r.Post("/token", handlePostToken)

}
