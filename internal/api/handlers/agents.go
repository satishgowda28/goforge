package handlers

import (
	"net/http"
)

// ListAgents returns all agents loaded from YAML configs.
// Implemented in Phase 2.
func ListAgents(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}

// StartRun starts a new agent run.
// Implemented in Phase 2.
func StartRun(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}
