package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-go-golems/bot-signup/internal/auth"
	"github.com/go-go-golems/bot-signup/internal/database"
)

func createTestUser(t *testing.T, srv *Server, discordID, email, role string) (*database.User, *http.Cookie) {
	t.Helper()
	user, err := srv.db.UpsertDiscordUser(t.Context(), discordID, email, "Test User", "")
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
	w := httptest.NewRecorder()
	if err := srv.sessions.WriteSession(w, user.ID); err != nil {
		t.Fatalf("write session: %v", err)
	}
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == auth.SessionCookieName {
			return user, cookie
		}
	}
	t.Fatal("session cookie not written")
	return nil, nil
}

func TestProfileUpdate(t *testing.T) {
	srv := newTestServer(t)
	_, session := createTestUser(t, srv, "111", "user@example.com", "user")
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
	req.AddCookie(session)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected profile 200, got %d body=%s", resp.Code, resp.Body.String())
	}

	updateBody := []byte(`{"email":"new@example.com","display_name":"New Name"}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/profile", bytes.NewReader(updateBody))
	updateReq.AddCookie(session)
	updateResp := httptest.NewRecorder()
	mux.ServeHTTP(updateResp, updateReq)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected update profile 200, got %d body=%s", updateResp.Code, updateResp.Body.String())
	}
}

func TestAdminApproveUser(t *testing.T) {
	srv := newTestServer(t)
	waitingUser, _ := createTestUser(t, srv, "222", "waiting@example.com", "user")
	_, adminSession := createTestUser(t, srv, "999", "admin@example.com", "admin")
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	waitlistReq := httptest.NewRequest(http.MethodGet, "/api/admin/waitlist", nil)
	waitlistReq.AddCookie(adminSession)
	waitlistResp := httptest.NewRecorder()
	mux.ServeHTTP(waitlistResp, waitlistReq)
	if waitlistResp.Code != http.StatusOK {
		t.Fatalf("expected waitlist 200, got %d body=%s", waitlistResp.Code, waitlistResp.Body.String())
	}

	approveBody := []byte(`{"application_id":"987","bot_token":"token","guild_id":"111","public_key":"abcdef"}`)
	approveReq := httptest.NewRequest(http.MethodPost, "/api/admin/users/"+strconv.FormatInt(waitingUser.ID, 10)+"/approve", bytes.NewReader(approveBody))
	approveReq.AddCookie(adminSession)
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
	_, session := createTestUser(t, srv, "333", "user@example.com", "user")
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/waitlist", nil)
	req.AddCookie(session)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", resp.Code, resp.Body.String())
	}
}
