package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-go-golems/bot-signup/internal/database"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type profileResponse struct {
	User           *database.User           `json:"user"`
	BotCredentials *database.BotCredentials `json:"bot_credentials"`
	Message        string                   `json:"message,omitempty"`
}

type updateProfileRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

func (s *Server) handleGetProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := s.currentUser(w, r)
	if !ok {
		return
	}
	resp := profileResponse{User: user}
	creds, err := s.db.GetCredentialsByUserID(r.Context(), user.ID)
	if err == nil {
		resp.BotCredentials = creds
	} else if !errors.Is(err, database.ErrNotFound) {
		respondError(w, http.StatusInternalServerError, "database error")
		return
	}
	if user.Status == database.UserStatusWaiting {
		resp.Message = "Your account is pending approval."
	}
	respondJSON(w, http.StatusOK, resp)
}

func (s *Server) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := s.currentUser(w, r)
	if !ok {
		return
	}
	var req updateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if req.Email != "" && !emailPattern.MatchString(req.Email) {
		respondError(w, http.StatusBadRequest, "invalid email")
		return
	}
	if len(req.DisplayName) < 2 || len(req.DisplayName) > 50 {
		respondError(w, http.StatusBadRequest, "invalid display name")
		return
	}
	updated, err := s.db.UpdateUserProfile(r.Context(), user.ID, req.Email, req.DisplayName)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"user": updated})
}

func (s *Server) currentUser(w http.ResponseWriter, r *http.Request) (*database.User, bool) {
	userID, ok := currentUserID(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "missing user context")
		return nil, false
	}
	user, err := s.db.GetUserByID(r.Context(), userID)
	if errors.Is(err, database.ErrNotFound) {
		respondError(w, http.StatusUnauthorized, "user not found")
		return nil, false
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "database error")
		return nil, false
	}
	return user, true
}
