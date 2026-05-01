package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-go-golems/bot-signup/internal/auth"
	"github.com/go-go-golems/bot-signup/internal/database"
)

func createTestUser(t *testing.T, srv *Server, discordID, email, role string) (*database.User, string) {
	t.Helper()
	hash, err := auth.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user, err := srv.db.CreateUser(t.Context(), discordID, email, "Test User", hash)
	if err != nil {
		t.Fatalf("create test user: %v", err)
	}
	if role != "" && role != string(user.Role) {
		if err := srv.db.UpdateUserRole(t.Context(), user.ID, database.UserRole(role)); err != nil {
			t.Fatalf("update role: %v", err)
		}
		user, err = srv.db.GetUserByID(t.Context(), user.ID)
		if err != nil {
			t.Fatalf("reload user: %v", err)
		}
	}
	token, err := auth.GenerateToken(user.ID, string(user.Role), srv.jwtSecret)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	return user, token
}

func TestProfileAndPassword(t *testing.T) {
	srv := newTestServer(t)
	_, token := createTestUser(t, srv, "111", "user@example.com", "user")
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected profile 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	updateBody := []byte(`{"email":"new@example.com","display_name":"New Name"}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateResp := httptest.NewRecorder()
	mux.ServeHTTP(updateResp, updateReq)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected update profile 200, got %d body=%s", updateResp.Code, updateResp.Body.String())
	}

	passwordBody := []byte(`{"current_password":"password123","new_password":"newpassword123"}`)
	passwordReq := httptest.NewRequest(http.MethodPut, "/api/profile/password", bytes.NewReader(passwordBody))
	passwordReq.Header.Set("Authorization", "Bearer "+token)
	passwordResp := httptest.NewRecorder()
	mux.ServeHTTP(passwordResp, passwordReq)
	if passwordResp.Code != http.StatusOK {
		t.Fatalf("expected password 200, got %d body=%s", passwordResp.Code, passwordResp.Body.String())
	}
}

func TestAdminApproveUser(t *testing.T) {
	srv := newTestServer(t)
	waitingUser, _ := createTestUser(t, srv, "222", "waiting@example.com", "user")
	_, adminToken := createTestUser(t, srv, "999", "admin@example.com", "admin")
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	waitlistReq := httptest.NewRequest(http.MethodGet, "/api/admin/waitlist", nil)
	waitlistReq.Header.Set("Authorization", "Bearer "+adminToken)
	waitlistResp := httptest.NewRecorder()
	mux.ServeHTTP(waitlistResp, waitlistReq)
	if waitlistResp.Code != http.StatusOK {
		t.Fatalf("expected waitlist 200, got %d body=%s", waitlistResp.Code, waitlistResp.Body.String())
	}

	approveBody := []byte(`{"application_id":"987","bot_token":"token","guild_id":"111","public_key":"abcdef"}`)
	approveReq := httptest.NewRequest(http.MethodPost, "/api/admin/users/"+jsonNumber(waitingUser.ID)+"/approve", bytes.NewReader(approveBody))
	approveReq.Header.Set("Authorization", "Bearer "+adminToken)
	approveResp := httptest.NewRecorder()
	mux.ServeHTTP(approveResp, approveReq)
	if approveResp.Code != http.StatusOK {
		t.Fatalf("expected approve 200, got %d body=%s", approveResp.Code, approveResp.Body.String())
	}
	updated, err := srv.db.GetUserByID(t.Context(), waitingUser.ID)
	if err != nil {
		t.Fatalf("get updated user: %v", err)
	}
	if updated.Status != database.UserStatusApproved {
		t.Fatalf("expected approved, got %s", updated.Status)
	}
	creds, err := srv.db.GetCredentialsByUserID(t.Context(), waitingUser.ID)
	if err != nil {
		t.Fatalf("get credentials: %v", err)
	}
	if creds.ApplicationID != "987" {
		t.Fatalf("unexpected credentials: %+v", creds)
	}
}

func TestAdminRouteRejectsNonAdmin(t *testing.T) {
	srv := newTestServer(t)
	_, token := createTestUser(t, srv, "333", "user@example.com", "user")
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/waitlist", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", resp.Code, resp.Body.String())
	}
}

func jsonNumber(n int64) string {
	b, _ := json.Marshal(n)
	return string(b)
}
