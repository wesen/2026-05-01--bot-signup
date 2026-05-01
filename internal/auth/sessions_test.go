package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSessionRoundTrip(t *testing.T) {
	mgr := NewSessionManager([]byte("secret"), false)
	w := httptest.NewRecorder()
	if err := mgr.WriteSession(w, 42); err != nil {
		t.Fatalf("write session: %v", err)
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one cookie, got %d", len(cookies))
	}
	if !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteLaxMode {
		t.Fatalf("unexpected cookie settings: %+v", cookies[0])
	}
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(cookies[0])
	userID, err := mgr.ReadSession(req)
	if err != nil {
		t.Fatalf("read session: %v", err)
	}
	if userID != 42 {
		t.Fatalf("expected user 42, got %d", userID)
	}
}

func TestOAuthStateRoundTrip(t *testing.T) {
	mgr := NewSessionManager([]byte("secret"), false)
	w := httptest.NewRecorder()
	state, err := mgr.WriteOAuthState(w, "/profile")
	if err != nil {
		t.Fatalf("write oauth state: %v", err)
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected one cookie, got %d", len(cookies))
	}
	req := httptest.NewRequest(http.MethodGet, "/auth/discord/callback?state="+state, nil)
	req.AddCookie(cookies[0])
	returnTo, err := mgr.ConsumeOAuthState(httptest.NewRecorder(), req, state)
	if err != nil {
		t.Fatalf("consume oauth state: %v", err)
	}
	if returnTo != "/profile" {
		t.Fatalf("expected /profile, got %q", returnTo)
	}
}

func TestSanitizeReturnTo(t *testing.T) {
	for _, value := range []string{"", "https://evil.example", "//evil.example"} {
		if got := sanitizeReturnTo(value); got != "/waiting-list" {
			t.Fatalf("expected unsafe %q to sanitize, got %q", value, got)
		}
	}
	if got := sanitizeReturnTo("/admin"); got != "/admin" {
		t.Fatalf("expected /admin, got %q", got)
	}
}
