package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-go-golems/bot-signup/internal/auth"
)

type contextKey string

const (
	contextKeyUserID contextKey = "user_id"
	contextKeyRole   contextKey = "role"
)

func (s *Server) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}
		tokenString, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok || tokenString == "" {
			respondError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}
		claims, err := auth.ParseToken(tokenString, s.jwtSecret)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "invalid token")
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, contextKeyRole, claims.Role)
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
