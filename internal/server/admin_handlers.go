package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-go-golems/bot-signup/internal/database"
)

type approveRequest struct {
	ApplicationID string `json:"application_id"`
	BotToken      string `json:"bot_token"`
	GuildID       string `json:"guild_id"`
	PublicKey     string `json:"public_key"`
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.db.GetStats(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "database error")
		return
	}
	respondJSON(w, http.StatusOK, stats)
}

func (s *Server) handleWaitlist(w http.ResponseWriter, r *http.Request) {
	users, total, err := s.db.ListUsersByStatus(r.Context(), database.UserStatusWaiting, 1, 100)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "database error")
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"users": users, "total": total})
}

func (s *Server) handleListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	status := database.UserStatus(r.URL.Query().Get("status"))
	users, total, err := s.db.ListUsers(r.Context(), database.ListUsersOptions{Status: status, Page: page, PerPage: perPage})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "database error")
		return
	}
	respondJSON(w, http.StatusOK, map[string]any{"users": users, "total": total, "page": max(page, 1), "per_page": max(perPage, 20)})
}

func (s *Server) handleApproveUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := pathID(w, r)
	if !ok {
		return
	}
	adminID, ok := currentUserID(r)
	if !ok {
		respondError(w, http.StatusUnauthorized, "missing admin context")
		return
	}
	var req approveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ApplicationID == "" || req.BotToken == "" || req.GuildID == "" || req.PublicKey == "" {
		respondError(w, http.StatusBadRequest, "all credential fields are required")
		return
	}
	creds, err := s.db.ApproveUser(r.Context(), userID, adminID, &database.BotCredentials{
		ApplicationID: req.ApplicationID,
		BotToken:      req.BotToken,
		GuildID:       req.GuildID,
		PublicKey:     req.PublicKey,
	})
	if errors.Is(err, database.ErrNotFound) {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusConflict, err.Error())
		return
	}
	user, _ := s.db.GetUserByID(r.Context(), userID)
	respondJSON(w, http.StatusOK, map[string]any{"message": "User approved successfully", "user": user, "bot_credentials": creds})
}

func (s *Server) handleRejectUser(w http.ResponseWriter, r *http.Request) {
	s.setUserStatus(w, r, database.UserStatusRejected, "User rejected")
}

func (s *Server) handleDisableUser(w http.ResponseWriter, r *http.Request) {
	s.setUserStatus(w, r, database.UserStatusSuspended, "User disabled")
}

func (s *Server) handleSuspendUser(w http.ResponseWriter, r *http.Request) {
	s.setUserStatus(w, r, database.UserStatusSuspended, "User disabled")
}

func (s *Server) handleEnableUser(w http.ResponseWriter, r *http.Request) {
	s.setUserStatus(w, r, database.UserStatusApproved, "User enabled")
}

func (s *Server) handleUpdateCredentials(w http.ResponseWriter, r *http.Request) {
	userID, ok := pathID(w, r)
	if !ok {
		return
	}
	var req approveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := s.db.UpdateBotCredentials(r.Context(), &database.BotCredentials{UserID: userID, ApplicationID: req.ApplicationID, BotToken: req.BotToken, GuildID: req.GuildID, PublicKey: req.PublicKey}); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			respondError(w, http.StatusNotFound, "credentials not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update credentials")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "Credentials updated"})
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := s.db.DeleteUser(r.Context(), userID); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": "User deleted"})
}

func (s *Server) setUserStatus(w http.ResponseWriter, r *http.Request, status database.UserStatus, message string) {
	userID, ok := pathID(w, r)
	if !ok {
		return
	}
	if err := s.db.UpdateUserStatus(r.Context(), userID, status); err != nil {
		if errors.Is(err, database.ErrNotFound) {
			respondError(w, http.StatusNotFound, "user not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update user")
		return
	}
	respondJSON(w, http.StatusOK, map[string]string{"message": message})
}

func pathID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return 0, false
	}
	return id, true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
