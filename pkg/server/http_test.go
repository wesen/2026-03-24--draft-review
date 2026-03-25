package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	draftauth "github.com/go-go-golems/draft-review/pkg/auth"
)

func TestHandleMeDevMode(t *testing.T) {
	t.Parallel()

	handler := NewHandler(HandlerOptions{
		Version:   "test",
		StartedAt: time.Now().UTC(),
		AuthSettings: &draftauth.Settings{
			Mode:      draftauth.AuthModeDev,
			DevUserID: "local-author",
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var response struct {
		Data draftauth.UserInfo `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Data.Authenticated {
		t.Fatalf("expected authenticated dev user")
	}
	if response.Data.AuthMode != draftauth.AuthModeDev {
		t.Fatalf("expected auth mode dev, got %q", response.Data.AuthMode)
	}
	if response.Data.Subject != "local-author" {
		t.Fatalf("expected subject local-author, got %q", response.Data.Subject)
	}
}

func TestHandleMeOIDCUnauthenticated(t *testing.T) {
	t.Parallel()

	handler := NewHandler(HandlerOptions{
		Version:   "test",
		StartedAt: time.Now().UTC(),
		AuthSettings: &draftauth.Settings{
			Mode: draftauth.AuthModeOIDC,
		},
	})

	request := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var response struct {
		Data draftauth.UserInfo `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Data.Authenticated {
		t.Fatalf("expected unauthenticated oidc user")
	}
	if response.Data.AuthMode != draftauth.AuthModeOIDC {
		t.Fatalf("expected auth mode oidc, got %q", response.Data.AuthMode)
	}
}
