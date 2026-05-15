package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/satishgowda28/goforge/internal/api/handlers"
)

func runsRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handlers.ListRuns)
	r.Get("/{id}", handlers.GetRun)
	r.Get("/{id}/stream", handlers.StreamRun)
	r.Delete("/{id}", handlers.CancelRun)

	return r
}
