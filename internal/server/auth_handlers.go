package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-go-golems/bot-signup/internal/auth"
	"github.com/go-go-golems/bot-signup/internal/database"
)

var (
	discordIDPattern = regexp.MustCompile(`^[0-9]{3,20}$`)
	emailPattern     = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
)

type signupRequest struct {
	DiscordID   string `json:"discord_id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string         `json:"token"`
	User  *database.User `json:"user"`
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.normalize()
	if errorsByField := validateSignup(req); len(errorsByField) > 0 {
		respondJSON(w, http.StatusBadRequest, map[string]any{"errors": errorsByField})
		return
	}

	if _, err := s.db.GetUserByDiscordID(r.Context(), req.DiscordID); err == nil {
		respondError(w, http.StatusConflict, "Discord ID already registered")
		return
	} else if !errors.Is(err, database.ErrNotFound) {
		respondError(w, http.StatusInternalServerError, "database error")
		return
	}
	if _, err := s.db.GetUserByEmail(r.Context(), req.Email); err == nil {
		respondError(w, http.StatusConflict, "Email already registered")
		return
	} else if !errors.Is(err, database.ErrNotFound) {
		respondError(w, http.StatusInternalServerError, "database error")
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}
	user, err := s.db.CreateUser(r.Context(), req.DiscordID, req.Email, req.DisplayName, passwordHash)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	token, err := auth.GenerateToken(user.ID, string(user.Role), s.jwtSecret)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	respondJSON(w, http.StatusCreated, authResponse{Token: token, User: user})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	user, err := s.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil || !auth.CheckPassword(user.PasswordHash, req.Password) {
		respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	token, err := auth.GenerateToken(user.ID, string(user.Role), s.jwtSecret)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	respondJSON(w, http.StatusOK, authResponse{Token: token, User: user})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
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

func (req *signupRequest) normalize() {
	req.DiscordID = strings.TrimSpace(req.DiscordID)
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.DisplayName = strings.TrimSpace(req.DisplayName)
}

func validateSignup(req signupRequest) map[string]string {
	errs := map[string]string{}
	if !discordIDPattern.MatchString(req.DiscordID) {
		errs["discord_id"] = "Discord ID must be a numeric string"
	}
	if !emailPattern.MatchString(req.Email) {
		errs["email"] = "Email must be valid"
	}
	if len(req.DisplayName) < 2 || len(req.DisplayName) > 50 {
		errs["display_name"] = "Display name must be between 2 and 50 characters"
	}
	if len(req.Password) < 8 {
		errs["password"] = "Password must be at least 8 characters"
	}
	return errs
}
