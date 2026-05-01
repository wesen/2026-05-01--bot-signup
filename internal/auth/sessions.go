package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	SessionCookieName    = "bot_signup_session"
	OAuthStateCookieName = "bot_signup_oauth_state"
)

// SessionManager signs and verifies small cookie payloads without storing
// browser-readable tokens in localStorage.
type SessionManager struct {
	secret []byte
	secure bool
}

func NewSessionManager(secret []byte, secure bool) *SessionManager {
	return &SessionManager{secret: secret, secure: secure}
}

func (m *SessionManager) WriteSession(w http.ResponseWriter, userID int64) error {
	expires := time.Now().Add(7 * 24 * time.Hour)
	payload := fmt.Sprintf("%d|%d", userID, expires.Unix())
	value := m.sign(payload)
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		MaxAge:   int(time.Until(expires).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   m.secure,
	})
	return nil
}

func (m *SessionManager) ReadSession(r *http.Request) (int64, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return 0, err
	}
	payload, err := m.verify(cookie.Value)
	if err != nil {
		return 0, err
	}
	parts := strings.Split(payload, "|")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid session payload")
	}
	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid session user id: %w", err)
	}
	expires, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid session expiry: %w", err)
	}
	if time.Now().Unix() > expires {
		return 0, fmt.Errorf("session expired")
	}
	return userID, nil
}

func (m *SessionManager) ClearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   m.secure,
	})
}

func (m *SessionManager) WriteOAuthState(w http.ResponseWriter, returnTo string) (string, error) {
	state, err := randomToken(32)
	if err != nil {
		return "", err
	}
	expires := time.Now().Add(10 * time.Minute)
	payload := fmt.Sprintf("%s|%d|%s", state, expires.Unix(), sanitizeReturnTo(returnTo))
	http.SetCookie(w, &http.Cookie{
		Name:     OAuthStateCookieName,
		Value:    m.sign(payload),
		Path:     "/auth/discord",
		Expires:  expires,
		MaxAge:   int(time.Until(expires).Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   m.secure,
	})
	return state, nil
}

func (m *SessionManager) ConsumeOAuthState(w http.ResponseWriter, r *http.Request, actual string) (string, error) {
	cookie, err := r.Cookie(OAuthStateCookieName)
	if err != nil {
		return "", err
	}
	payload, err := m.verify(cookie.Value)
	if err != nil {
		return "", err
	}
	m.clearOAuthState(w)

	parts := strings.SplitN(payload, "|", 3)
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid oauth state payload")
	}
	if !hmac.Equal([]byte(parts[0]), []byte(actual)) {
		return "", fmt.Errorf("oauth state mismatch")
	}
	expires, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid oauth state expiry: %w", err)
	}
	if time.Now().Unix() > expires {
		return "", fmt.Errorf("oauth state expired")
	}
	return sanitizeReturnTo(parts[2]), nil
}

func (m *SessionManager) clearOAuthState(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     OAuthStateCookieName,
		Value:    "",
		Path:     "/auth/discord",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   m.secure,
	})
}

func (m *SessionManager) sign(payload string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(payload))
	sig := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString([]byte(payload)) + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func (m *SessionManager) verify(value string) (string, error) {
	encodedPayload, encodedSig, ok := strings.Cut(value, ".")
	if !ok {
		return "", fmt.Errorf("missing signature")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return "", err
	}
	sig, err := base64.RawURLEncoding.DecodeString(encodedSig)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, m.secret)
	mac.Write(payloadBytes)
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return "", fmt.Errorf("invalid signature")
	}
	return string(payloadBytes), nil
}

func randomToken(bytes int) (string, error) {
	buf := make([]byte, bytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func sanitizeReturnTo(returnTo string) string {
	if returnTo == "" || !strings.HasPrefix(returnTo, "/") || strings.HasPrefix(returnTo, "//") || strings.Contains(returnTo, "://") {
		return "/waiting-list"
	}
	return returnTo
}
