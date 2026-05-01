package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/go-go-golems/bot-signup/internal/database"
)

func (s *Server) handleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	state, err := s.sessions.WriteOAuthState(w, r.URL.Query().Get("return_to"))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create oauth state")
		return
	}
	http.Redirect(w, r, s.discordOAuth.AuthCodeURL(state), http.StatusFound)
}

func (s *Server) handleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	returnTo, err := s.sessions.ConsumeOAuthState(w, r, r.URL.Query().Get("state"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid oauth state")
		return
	}
	code := r.URL.Query().Get("code")
	if code == "" {
		respondError(w, http.StatusBadRequest, "missing oauth code")
		return
	}
	discordUser, err := s.discordOAuth.ExchangeAndFetchUser(r.Context(), code)
	if err != nil {
		log.Printf("discord oauth callback failed: %v", err)
		respondError(w, http.StatusBadGateway, "discord oauth failed")
		return
	}
	user, err := s.db.UpsertDiscordUser(r.Context(), discordUser.ID, discordUser.Email, discordUser.DisplayName(), discordUser.AvatarURL())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	if err := s.sessions.WriteSession(w, user.ID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create session")
		return
	}
	if returnTo == "" || returnTo == "/waiting-list" {
		returnTo = routeForUser(user)
	}
	http.Redirect(w, r, returnTo, http.StatusFound)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	s.sessions.ClearSession(w)
	respondJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := currentUserID(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "missing user context")
		return
	}
	user, err := s.db.GetUserByID(r.Context(), userID)
	if errors.Is(err, database.ErrNotFound) {
		respondError(w, http.StatusUnauthorized, "user not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "database error")
		return
	}
	respondJSON(w, http.StatusOK, user)
}

func routeForUser(user *database.User) string {
	if user.Role == database.UserRoleAdmin {
		return "/admin"
	}
	switch user.Status {
	case database.UserStatusApproved:
		return "/profile"
	default:
		return "/waiting-list"
	}
}
