package server

import (
	"context"
	"net/http"

	"github.com/go-go-golems/bot-signup/internal/auth"
	"github.com/go-go-golems/bot-signup/internal/database"
)

type discordOAuthClient interface {
	AuthCodeURL(state string) string
	ExchangeAndFetchUser(ctx context.Context, code string) (*auth.DiscordUser, error)
}

type Options struct {
	Version              string
	SessionSecret        []byte
	SecureCookies        bool
	DiscordClientID      string
	DiscordClientSecret  string
	DiscordRedirectURL   string
	DiscordOAuthOverride discordOAuthClient
}

// Server owns the HTTP handlers for the bot signup application.
type Server struct {
	db           *database.DB
	sessions     *auth.SessionManager
	discordOAuth discordOAuthClient
	version      string
}

// New constructs a Server with production-safe defaults.
func New(db *database.DB, opts Options) *Server {
	version := opts.Version
	if version == "" {
		version = "dev"
	}
	secret := opts.SessionSecret
	if len(secret) == 0 {
		secret = []byte("dev-insecure-change-me")
	}
	discordOAuth := opts.DiscordOAuthOverride
	if discordOAuth == nil {
		discordOAuth = auth.NewDiscordOAuth(opts.DiscordClientID, opts.DiscordClientSecret, opts.DiscordRedirectURL)
	}
	return &Server{
		db:           db,
		sessions:     auth.NewSessionManager(secret, opts.SecureCookies),
		discordOAuth: discordOAuth,
		version:      version,
	}
}

// RegisterRoutes registers API and auth routes. The SPA fallback will be added
// in a later phase after the frontend is built.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/stats", s.handleStats)

	mux.HandleFunc("GET /auth/discord/login", s.handleDiscordLogin)
	mux.HandleFunc("GET /auth/discord/callback", s.handleDiscordCallback)
	mux.HandleFunc("POST /api/auth/logout", s.handleLogout)
	mux.HandleFunc("GET /api/auth/session", s.handleSession)
	mux.HandleFunc("GET /api/auth/me", s.SessionMiddleware(s.handleMe))

	mux.HandleFunc("GET /api/profile", s.SessionMiddleware(s.handleGetProfile))
	mux.HandleFunc("PUT /api/profile", s.SessionMiddleware(s.handleUpdateProfile))

	mux.HandleFunc("GET /api/admin/waitlist", s.SessionMiddleware(AdminOnly(s.handleWaitlist)))
	mux.HandleFunc("GET /api/admin/users", s.SessionMiddleware(AdminOnly(s.handleListUsers)))
	mux.HandleFunc("POST /api/admin/users/{id}/approve", s.SessionMiddleware(AdminOnly(s.handleApproveUser)))
	mux.HandleFunc("POST /api/admin/users/{id}/reject", s.SessionMiddleware(AdminOnly(s.handleRejectUser)))
	mux.HandleFunc("POST /api/admin/users/{id}/disable", s.SessionMiddleware(AdminOnly(s.handleDisableUser)))
	mux.HandleFunc("POST /api/admin/users/{id}/suspend", s.SessionMiddleware(AdminOnly(s.handleSuspendUser)))
	mux.HandleFunc("PUT /api/admin/users/{id}/credentials", s.SessionMiddleware(AdminOnly(s.handleUpdateCredentials)))
	mux.HandleFunc("DELETE /api/admin/users/{id}", s.SessionMiddleware(AdminOnly(s.handleDeleteUser)))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": s.version,
	})
}
