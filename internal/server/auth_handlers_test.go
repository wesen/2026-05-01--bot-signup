package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/bot-signup/internal/database"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	db, err := database.Open(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return New(db, []byte("test-secret"), "test")
}

func TestSignupLoginAndMe(t *testing.T) {
	srv := newTestServer(t)
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	signupBody := []byte(`{"discord_id":"123456789","email":"user@example.com","display_name":"CoolBotDev","password":"password123"}`)
	signupReq := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(signupBody))
	signupReq.Header.Set("Content-Type", "application/json")
	signupResp := httptest.NewRecorder()
	mux.ServeHTTP(signupResp, signupReq)
	if signupResp.Code != http.StatusCreated {
		t.Fatalf("expected signup 201, got %d body=%s", signupResp.Code, signupResp.Body.String())
	}
	var signup authResponse
	if err := json.NewDecoder(signupResp.Body).Decode(&signup); err != nil {
		t.Fatalf("decode signup: %v", err)
	}
	if signup.Token == "" || signup.User.Status != database.UserStatusWaiting {
		t.Fatalf("unexpected signup response: %+v", signup)
	}

	loginBody := []byte(`{"email":"user@example.com","password":"password123"}`)
	loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp := httptest.NewRecorder()
	mux.ServeHTTP(loginResp, loginReq)
	if loginResp.Code != http.StatusOK {
		t.Fatalf("expected login 200, got %d body=%s", loginResp.Code, loginResp.Body.String())
	}
	var login authResponse
	if err := json.NewDecoder(loginResp.Body).Decode(&login); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	if login.Token == "" {
		t.Fatal("expected login token")
	}

	meReq := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	meReq.Header.Set("Authorization", "Bearer "+login.Token)
	meResp := httptest.NewRecorder()
	mux.ServeHTTP(meResp, meReq)
	if meResp.Code != http.StatusOK {
		t.Fatalf("expected me 200, got %d body=%s", meResp.Code, meResp.Body.String())
	}
	var me database.User
	if err := json.NewDecoder(meResp.Body).Decode(&me); err != nil {
		t.Fatalf("decode me: %v", err)
	}
	if me.Email != "user@example.com" {
		t.Fatalf("unexpected me response: %+v", me)
	}
}

func TestSignupValidation(t *testing.T) {
	srv := newTestServer(t)
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	body := []byte(`{"discord_id":"abc","email":"bad","display_name":"x","password":"short"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", resp.Code, resp.Body.String())
	}
}
