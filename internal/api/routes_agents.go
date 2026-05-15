package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/satishgowda28/goforge/internal/api/handlers"
)

func agentsRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handlers.ListAgents)
	r.Post("/run", handlers.StartRun)

	return r
}
