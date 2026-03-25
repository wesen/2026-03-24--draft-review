package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestResolveRequestedRedirectAllowsRelativePath(t *testing.T) {
	authenticator := &OIDCAuthenticator{}
	request := httptest.NewRequest("GET", "http://127.0.0.1:8080/auth/login?return_to=%2Fapp", nil)

	redirect, err := authenticator.resolveRequestedRedirect(request, "/app")
	if err != nil {
		t.Fatalf("resolveRequestedRedirect returned error: %v", err)
	}
	if redirect != "/app" {
		t.Fatalf("expected /app, got %q", redirect)
	}
}

func TestResolveRequestedRedirectAllowsSameHostDifferentPort(t *testing.T) {
	authenticator := &OIDCAuthenticator{}
	authenticator.oauthConfig.RedirectURL = "http://127.0.0.1:8080/auth/callback"
	request := httptest.NewRequest("GET", "http://127.0.0.1:8080/auth/login", nil)

	redirect, err := authenticator.resolveRequestedRedirect(request, "http://127.0.0.1:5173/app")
	if err != nil {
		t.Fatalf("resolveRequestedRedirect returned error: %v", err)
	}
	if redirect != "http://127.0.0.1:5173/app" {
		t.Fatalf("unexpected redirect target %q", redirect)
	}
}

func TestResolveRequestedRedirectRejectsDifferentHost(t *testing.T) {
	authenticator := &OIDCAuthenticator{}
	authenticator.oauthConfig.RedirectURL = "http://127.0.0.1:8080/auth/callback"
	request := httptest.NewRequest("GET", "http://127.0.0.1:8080/auth/login", nil)

	if _, err := authenticator.resolveRequestedRedirect(request, "https://evil.example.com/app"); err == nil {
		t.Fatal("expected resolveRequestedRedirect to reject different host")
	}
}

func TestBuildLogoutCallbackURLUsesBackendHostAndReturnTo(t *testing.T) {
	redirectURL, err := buildLogoutCallbackURL("http://127.0.0.1:8080/auth/callback", "http://127.0.0.1:5173/")
	if err != nil {
		t.Fatalf("buildLogoutCallbackURL returned error: %v", err)
	}
	if redirectURL != "http://127.0.0.1:8080/auth/logout/callback?return_to=http%3A%2F%2F127.0.0.1%3A5173%2F" {
		t.Fatalf("expected backend logout callback redirect, got %q", redirectURL)
	}
}

func TestOAuthStateRoundTripPreservesReturnTo(t *testing.T) {
	value, err := marshalOAuthState(oauthStatePayload{
		ID:       "state-123",
		ReturnTo: "http://127.0.0.1:5173/app",
	})
	if err != nil {
		t.Fatalf("marshalOAuthState returned error: %v", err)
	}

	payload, err := unmarshalOAuthState(value)
	if err != nil {
		t.Fatalf("unmarshalOAuthState returned error: %v", err)
	}
	if payload.ID != "state-123" {
		t.Fatalf("expected state id to round-trip, got %q", payload.ID)
	}
	if payload.ReturnTo != "http://127.0.0.1:5173/app" {
		t.Fatalf("expected return_to to round-trip, got %q", payload.ReturnTo)
	}
}

func TestUnmarshalOAuthStateRejectsMissingID(t *testing.T) {
	raw, err := json.Marshal(oauthStatePayload{ReturnTo: "/app"})
	if err != nil {
		t.Fatalf("json.Marshal returned error: %v", err)
	}

	value := base64.RawURLEncoding.EncodeToString(raw)
	if _, err := unmarshalOAuthState(value); err == nil {
		t.Fatal("expected unmarshalOAuthState to reject missing state id")
	}
}
