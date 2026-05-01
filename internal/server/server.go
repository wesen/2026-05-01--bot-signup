package server

import (
	"net/http"

	"github.com/go-go-golems/bot-signup/internal/database"
)

// Server owns the HTTP handlers for the bot signup application.
type Server struct {
	db        *database.DB
	jwtSecret []byte
	version   string
}

// New constructs a Server with production-safe defaults.
func New(db *database.DB, jwtSecret []byte, version string) *Server {
	if version == "" {
		version = "dev"
	}
	return &Server{db: db, jwtSecret: jwtSecret, version: version}
}

// RegisterRoutes registers public API routes. The SPA fallback will be added
// in a later phase after the frontend is built.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/stats", s.handleStats)

	mux.HandleFunc("POST /api/auth/signup", s.handleSignup)
	mux.HandleFunc("POST /api/auth/login", s.handleLogin)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
	mux.HandleFunc("GET /api/auth/me", s.AuthMiddleware(s.handleMe))

	mux.HandleFunc("GET /api/profile", s.AuthMiddleware(s.handleGetProfile))
	mux.HandleFunc("PUT /api/profile", s.AuthMiddleware(s.handleUpdateProfile))
	mux.HandleFunc("PUT /api/profile/password", s.AuthMiddleware(s.handleChangePassword))

	mux.HandleFunc("GET /api/admin/waitlist", s.AuthMiddleware(AdminOnly(s.handleWaitlist)))
	mux.HandleFunc("GET /api/admin/users", s.AuthMiddleware(AdminOnly(s.handleListUsers)))
	mux.HandleFunc("POST /api/admin/users/{id}/approve", s.AuthMiddleware(AdminOnly(s.handleApproveUser)))
	mux.HandleFunc("POST /api/admin/users/{id}/reject", s.AuthMiddleware(AdminOnly(s.handleRejectUser)))
	mux.HandleFunc("POST /api/admin/users/{id}/suspend", s.AuthMiddleware(AdminOnly(s.handleSuspendUser)))
	mux.HandleFunc("PUT /api/admin/users/{id}/credentials", s.AuthMiddleware(AdminOnly(s.handleUpdateCredentials)))
	mux.HandleFunc("DELETE /api/admin/users/{id}", s.AuthMiddleware(AdminOnly(s.handleDeleteUser)))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": s.version,
	})
}
