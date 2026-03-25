package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSessionManagerRoundTrip(t *testing.T) {
	manager, err := NewSessionManager("test_session", "super-secret", "http://127.0.0.1:8080/auth/callback")
	if err != nil {
		t.Fatalf("NewSessionManager returned error: %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/", nil)

	expectedExpiry := time.Now().UTC().Add(30 * time.Minute).Round(time.Second)
	err = manager.WriteSession(recorder, request, SessionClaims{
		Issuer:            "https://auth.example.com/realms/draft-review",
		Subject:           "user-123",
		Email:             "alice@example.com",
		PreferredUsername: "alice",
		DisplayName:       "Alice",
		Scopes:            []string{"openid", "profile", "email"},
		IssuedAt:          time.Now().UTC(),
		ExpiresAt:         expectedExpiry,
	})
	if err != nil {
		t.Fatalf("WriteSession returned error: %v", err)
	}

	response := recorder.Result()
	for _, cookie := range response.Cookies() {
		request.AddCookie(cookie)
	}

	claims, err := manager.ReadSession(request)
	if err != nil {
		t.Fatalf("ReadSession returned error: %v", err)
	}

	if claims.Subject != "user-123" {
		t.Fatalf("expected subject user-123, got %q", claims.Subject)
	}
	if claims.Email != "alice@example.com" {
		t.Fatalf("expected email alice@example.com, got %q", claims.Email)
	}
	if !claims.ExpiresAt.Equal(expectedExpiry) {
		t.Fatalf("expected expiry %s, got %s", expectedExpiry, claims.ExpiresAt)
	}
}

func TestSessionManagerRejectsTampering(t *testing.T) {
	manager, err := NewSessionManager("test_session", "super-secret", "http://127.0.0.1:8080/auth/callback")
	if err != nil {
		t.Fatalf("NewSessionManager returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/", nil)
	request.AddCookie(&http.Cookie{Name: "test_session", Value: "not-a-valid-token"})

	if _, err := manager.ReadSession(request); err == nil {
		t.Fatal("expected ReadSession to reject malformed token")
	}
}
