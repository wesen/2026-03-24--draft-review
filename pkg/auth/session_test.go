package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeSessionStore struct {
	createdSession *Session
	resolved       *ResolvedSession
	revokedHash    string
	touchedAt      time.Time
}

func (f *fakeSessionStore) CreateAuthorSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (*Session, error) {
	f.createdSession = &Session{
		ID:         uuid.NewString(),
		UserID:     userID.String(),
		TokenHash:  tokenHash,
		ExpiresAt:  expiresAt.UTC(),
		CreatedAt:  time.Now().UTC(),
		LastUsedAt: time.Now().UTC(),
	}
	return f.createdSession, nil
}

func (f *fakeSessionStore) FindAuthorSessionByTokenHash(ctx context.Context, tokenHash string) (*ResolvedSession, error) {
	if f.resolved == nil {
		return nil, ErrNotFound
	}
	if f.resolved.Session.TokenHash != tokenHash {
		return nil, ErrNotFound
	}
	return f.resolved, nil
}

func (f *fakeSessionStore) TouchAuthorSession(ctx context.Context, sessionID uuid.UUID, touchedAt time.Time, renewedExpiresAt *time.Time) (*Session, error) {
	if f.resolved == nil || f.resolved.Session.ID != sessionID.String() {
		return nil, ErrNotFound
	}
	f.touchedAt = touchedAt.UTC()
	f.resolved.Session.LastUsedAt = touchedAt.UTC()
	if renewedExpiresAt != nil {
		f.resolved.Session.ExpiresAt = renewedExpiresAt.UTC()
	}
	return &f.resolved.Session, nil
}

func (f *fakeSessionStore) RevokeAuthorSessionByTokenHash(ctx context.Context, tokenHash string) error {
	f.revokedHash = tokenHash
	return nil
}

func TestSessionManagerRoundTrip(t *testing.T) {
	store := &fakeSessionStore{}
	manager, err := NewSessionManager("test_session", "super-secret", "http://127.0.0.1:8080/auth/callback", 12*time.Hour, false, 0, store)
	if err != nil {
		t.Fatalf("NewSessionManager returned error: %v", err)
	}
	manager.now = func() time.Time {
		return time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC)
	}
	manager.newToken = func() (string, error) {
		return "opaque-token-123", nil
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/", nil)

	user := &User{
		ID:              "11111111-1111-1111-1111-111111111111",
		AuthIssuer:      "https://auth.example.com/realms/draft-review",
		AuthSubject:     "user-123",
		Email:           "alice@example.com",
		Name:            "Alice",
		EmailVerifiedAt: ptrTime(time.Date(2026, 3, 26, 11, 0, 0, 0, time.UTC)),
	}

	if err := manager.WriteSession(context.Background(), recorder, request, user); err != nil {
		t.Fatalf("WriteSession returned error: %v", err)
	}

	if store.createdSession == nil {
		t.Fatal("expected session to be created in store")
	}

	response := recorder.Result()
	for _, cookie := range response.Cookies() {
		request.AddCookie(cookie)
	}

	store.resolved = &ResolvedSession{
		Session: *store.createdSession,
		User:    *user,
	}

	claims, err := manager.ReadSession(context.Background(), nil, request)
	if err != nil {
		t.Fatalf("ReadSession returned error: %v", err)
	}

	if claims.Subject != "user-123" {
		t.Fatalf("expected subject user-123, got %q", claims.Subject)
	}
	if claims.Email != "alice@example.com" {
		t.Fatalf("expected email alice@example.com, got %q", claims.Email)
	}
	if claims.PreferredUsername != "alice" {
		t.Fatalf("expected preferred username alice, got %q", claims.PreferredUsername)
	}
	if !claims.ExpiresAt.Equal(store.createdSession.ExpiresAt) {
		t.Fatalf("expected expiry %s, got %s", store.createdSession.ExpiresAt, claims.ExpiresAt)
	}
	if store.touchedAt.IsZero() {
		t.Fatal("expected session read to touch last_used_at")
	}
	if !store.resolved.Session.LastUsedAt.Equal(store.touchedAt) {
		t.Fatalf("expected last_used_at %s, got %s", store.touchedAt, store.resolved.Session.LastUsedAt)
	}
}

func TestSessionManagerRejectsUnknownToken(t *testing.T) {
	store := &fakeSessionStore{}
	manager, err := NewSessionManager("test_session", "super-secret", "http://127.0.0.1:8080/auth/callback", 12*time.Hour, false, 0, store)
	if err != nil {
		t.Fatalf("NewSessionManager returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/", nil)
	request.AddCookie(&http.Cookie{Name: "test_session", Value: "unknown-token"})

	if _, err := manager.ReadSession(context.Background(), nil, request); err == nil {
		t.Fatal("expected ReadSession to reject unknown token")
	}
}

func TestSessionManagerClearSessionRevokesServerSession(t *testing.T) {
	store := &fakeSessionStore{}
	manager, err := NewSessionManager("test_session", "super-secret", "http://127.0.0.1:8080/auth/callback", 12*time.Hour, false, 0, store)
	if err != nil {
		t.Fatalf("NewSessionManager returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/", nil)
	request.AddCookie(&http.Cookie{Name: "test_session", Value: "opaque-token-123"})
	recorder := httptest.NewRecorder()

	manager.ClearSession(context.Background(), recorder, request)

	if store.revokedHash == "" {
		t.Fatal("expected server-side session to be revoked")
	}
}

func TestSessionManagerRenewsActiveSessionNearExpiry(t *testing.T) {
	store := &fakeSessionStore{}
	manager, err := NewSessionManager("test_session", "super-secret", "http://127.0.0.1:8080/auth/callback", 12*time.Hour, true, time.Hour, store)
	if err != nil {
		t.Fatalf("NewSessionManager returned error: %v", err)
	}
	now := time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC)
	manager.now = func() time.Time {
		return now
	}
	manager.newToken = func() (string, error) {
		return "opaque-token-456", nil
	}

	user := &User{
		ID:          "11111111-1111-1111-1111-111111111111",
		AuthIssuer:  "https://auth.example.com/realms/draft-review",
		AuthSubject: "user-123",
		Email:       "alice@example.com",
		Name:        "Alice",
	}

	writeRequest := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/", nil)
	writeRecorder := httptest.NewRecorder()
	if err := manager.WriteSession(context.Background(), writeRecorder, writeRequest, user); err != nil {
		t.Fatalf("WriteSession returned error: %v", err)
	}

	store.resolved = &ResolvedSession{
		Session: *store.createdSession,
		User:    *user,
	}
	store.resolved.Session.ExpiresAt = now.Add(45 * time.Minute)

	request := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/api/me", nil)
	for _, cookie := range writeRecorder.Result().Cookies() {
		request.AddCookie(cookie)
	}
	recorder := httptest.NewRecorder()

	active, err := manager.ReadSessionState(context.Background(), recorder, request)
	if err != nil {
		t.Fatalf("ReadSessionState returned error: %v", err)
	}

	if !active.Renewed {
		t.Fatal("expected session renewal")
	}
	expectedExpiry := now.Add(12 * time.Hour)
	if !store.resolved.Session.ExpiresAt.Equal(expectedExpiry) {
		t.Fatalf("expected renewed expiry %s, got %s", expectedExpiry, store.resolved.Session.ExpiresAt)
	}
	cookies := recorder.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected renewed cookie to be issued")
	}
	if !cookies[0].Expires.Equal(expectedExpiry) {
		t.Fatalf("expected renewed cookie expiry %s, got %s", expectedExpiry, cookies[0].Expires)
	}
}

func ptrTime(value time.Time) *time.Time {
	return &value
}
