package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if hash == "password123" {
		t.Fatal("hash should not equal plaintext")
	}
	if !CheckPassword(hash, "password123") {
		t.Fatal("expected password check to pass")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatal("expected wrong password to fail")
	}
}

func TestGenerateAndParseToken(t *testing.T) {
	secret := []byte("secret")
	token, err := GenerateToken(42, "admin", secret)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}
	claims, err := ParseToken(token, secret)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if claims.UserID != 42 || claims.Role != "admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if _, err := ParseToken(token, []byte("wrong-secret")); err == nil {
		t.Fatal("expected wrong secret to fail")
	}
}
