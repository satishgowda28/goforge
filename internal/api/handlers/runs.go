package handlers

import (
	"net/http"
)

// ListRuns returns recent agent runs.
// Implemented in Phase 4 (persistence).
func ListRuns(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}

// GetRun returns a single run with all its steps.
// Implemented in Phase 4 (persistence).
func GetRun(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}

// StreamRun streams live agent events via SSE.
// Implemented in Phase 3 (streaming).
func StreamRun(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}

// CancelRun cancels an in-progress agent run.
// Implemented in Phase 3 (concurrency).
func CancelRun(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("not implemented"))
}
