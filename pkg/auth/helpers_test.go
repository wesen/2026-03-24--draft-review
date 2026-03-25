package auth

import (
	"testing"
	"time"
)

func TestHashPasswordRoundTrip(t *testing.T) {
	t.Parallel()

	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if err := ComparePasswordHash(hash, "correct horse battery staple"); err != nil {
		t.Fatalf("ComparePasswordHash returned error: %v", err)
	}

	if err := ComparePasswordHash(hash, "wrong password"); err == nil {
		t.Fatalf("ComparePasswordHash should reject the wrong password")
	}
}

func TestGenerateOpaqueToken(t *testing.T) {
	t.Parallel()

	token, tokenHash, err := GenerateOpaqueToken()
	if err != nil {
		t.Fatalf("GenerateOpaqueToken returned error: %v", err)
	}

	if token == "" {
		t.Fatalf("token should not be empty")
	}

	if tokenHash == "" {
		t.Fatalf("token hash should not be empty")
	}

	if HashToken(token) != tokenHash {
		t.Fatalf("token hash should match HashToken output")
	}
}

func TestSessionCookies(t *testing.T) {
	t.Parallel()

	expiresAt := time.Now().Add(2 * time.Hour)
	cookie := NewSessionCookie("plain-token", expiresAt, true)

	if cookie.Name != DefaultSessionCookieName {
		t.Fatalf("unexpected cookie name: %s", cookie.Name)
	}

	if cookie.Value != "plain-token" {
		t.Fatalf("unexpected cookie value: %s", cookie.Value)
	}

	if !cookie.HttpOnly {
		t.Fatalf("session cookie should be httpOnly")
	}

	if !cookie.Secure {
		t.Fatalf("session cookie should be secure when requested")
	}

	expired := ExpireSessionCookie(false)
	if expired.MaxAge != -1 {
		t.Fatalf("expired cookie should set MaxAge=-1, got %d", expired.MaxAge)
	}

	if expired.Value != "" {
		t.Fatalf("expired cookie should clear the value")
	}
}
