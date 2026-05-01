package server

import (
	"net/http"

	"github.com/go-go-golems/bot-signup/internal/database"
)

// Server owns the HTTP handlers for the bot signup application.
type Server struct {
	db      *database.DB
	version string
}

// New constructs a Server with production-safe defaults.
func New(db *database.DB, version string) *Server {
	if version == "" {
		version = "dev"
	}
	return &Server{db: db, version: version}
}

// RegisterRoutes registers public API routes. The SPA fallback will be added
// in a later phase after the frontend is built.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", s.handleHealth)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": s.version,
	})
}
