package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/bot-signup/internal/auth"
	"github.com/go-go-golems/bot-signup/internal/database"
)

type fakeDiscordOAuth struct {
	user *auth.DiscordUser
}

func (f fakeDiscordOAuth) AuthCodeURL(state string) string {
	return "https://discord.example/authorize?state=" + state
}

func (f fakeDiscordOAuth) ExchangeAndFetchUser(ctx context.Context, code string) (*auth.DiscordUser, error) {
	return f.user, nil
}

func newTestServer(t *testing.T) *Server {
	t.Helper()
	db, err := database.Open(context.Background(), filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return New(db, Options{
		Version:       "test",
		SessionSecret: []byte("test-secret"),
		DiscordOAuthOverride: fakeDiscordOAuth{user: &auth.DiscordUser{
			ID:         "123456789",
			Username:   "coolbotdev",
			GlobalName: "CoolBotDev",
			Email:      "user@example.com",
		}},
	})
}

func TestDiscordOAuthLoginCallbackAndMe(t *testing.T) {
	srv := newTestServer(t)
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	loginReq := httptest.NewRequest(http.MethodGet, "/auth/discord/login?return_to=/waiting-list", nil)
	loginResp := httptest.NewRecorder()
	mux.ServeHTTP(loginResp, loginReq)
	if loginResp.Code != http.StatusFound {
		t.Fatalf("expected login redirect 302, got %d body=%s", loginResp.Code, loginResp.Body.String())
	}
	var stateCookie *http.Cookie
	for _, cookie := range loginResp.Result().Cookies() {
		if cookie.Name == auth.OAuthStateCookieName {
			stateCookie = cookie
		}
	}
	if stateCookie == nil || !stateCookie.HttpOnly {
		t.Fatalf("expected http-only oauth state cookie, got %+v", loginResp.Result().Cookies())
	}

	state := loginResp.Result().Header.Get("Location")[len("https://discord.example/authorize?state="):]
	callbackReq := httptest.NewRequest(http.MethodGet, "/auth/discord/callback?code=fake-code&state="+state, nil)
	callbackReq.AddCookie(stateCookie)
	callbackResp := httptest.NewRecorder()
	mux.ServeHTTP(callbackResp, callbackReq)
	if callbackResp.Code != http.StatusFound {
		t.Fatalf("expected callback redirect 302, got %d body=%s", callbackResp.Code, callbackResp.Body.String())
	}

	var sessionCookie *http.Cookie
	for _, cookie := range callbackResp.Result().Cookies() {
		if cookie.Name == auth.SessionCookieName {
			sessionCookie = cookie
		}
	}
	if sessionCookie == nil || !sessionCookie.HttpOnly {
		t.Fatalf("expected http-only session cookie, got %+v", callbackResp.Result().Cookies())
	}

	meReq := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	meReq.AddCookie(sessionCookie)
	meResp := httptest.NewRecorder()
	mux.ServeHTTP(meResp, meReq)
	if meResp.Code != http.StatusOK {
		t.Fatalf("expected me 200, got %d body=%s", meResp.Code, meResp.Body.String())
	}
	var me database.User
	if err := json.NewDecoder(meResp.Body).Decode(&me); err != nil {
		t.Fatalf("decode me: %v", err)
	}
	if me.DiscordID != "123456789" || me.Email != "user@example.com" {
		t.Fatalf("unexpected me response: %+v", me)
	}
}

func TestDiscordCallbackRejectsBadState(t *testing.T) {
	srv := newTestServer(t)
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/auth/discord/callback?code=fake&state=bad", nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)
	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", resp.Code, resp.Body.String())
	}
}
