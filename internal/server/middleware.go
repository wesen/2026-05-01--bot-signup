package server

import (
	"context"
	"net/http"
)

type contextKey string

const (
	contextKeyUserID contextKey = "user_id"
	contextKeyRole   contextKey = "role"
)

func (s *Server) SessionMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := s.sessions.ReadSession(r)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "not authenticated")
			return
		}
		user, err := s.db.GetUserByID(r.Context(), userID)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "not authenticated")
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUserID, user.ID)
		ctx = context.WithValue(ctx, contextKeyRole, string(user.Role))
		next(w, r.WithContext(ctx))
	}
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		role, _ := r.Context().Value(contextKeyRole).(string)
		if role != "admin" {
			respondError(w, http.StatusForbidden, "forbidden")
			return
		}
		next(w, r)
	}
}

func currentUserID(r *http.Request) (int64, bool) {
	id, ok := r.Context().Value(contextKeyUserID).(int64)
	return id, ok
}
