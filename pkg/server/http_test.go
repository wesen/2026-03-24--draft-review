package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"

	"github.com/go-go-golems/draft-review/pkg/analytics"
	"github.com/go-go-golems/draft-review/pkg/articles"
	draftauth "github.com/go-go-golems/draft-review/pkg/auth"
	"github.com/go-go-golems/draft-review/pkg/reviewlinks"
	"github.com/go-go-golems/draft-review/pkg/reviews"
	"github.com/google/uuid"
)

type fakeSessionStore struct {
	session *draftauth.ResolvedSession
}

func (f *fakeSessionStore) CreateAuthorSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (*draftauth.Session, error) {
	if f.session == nil {
		f.session = &draftauth.ResolvedSession{}
	}
	f.session.Session = draftauth.Session{
		ID:         uuid.NewString(),
		UserID:     userID.String(),
		TokenHash:  tokenHash,
		ExpiresAt:  expiresAt.UTC(),
		CreatedAt:  time.Now().UTC(),
		LastUsedAt: time.Now().UTC(),
	}
	return &f.session.Session, nil
}

func (f *fakeSessionStore) FindAuthorSessionByTokenHash(ctx context.Context, tokenHash string) (*draftauth.ResolvedSession, error) {
	if f.session == nil || f.session.Session.TokenHash != tokenHash {
		return nil, draftauth.ErrNotFound
	}
	return f.session, nil
}

func (f *fakeSessionStore) TouchAuthorSession(ctx context.Context, sessionID uuid.UUID, touchedAt time.Time, renewedExpiresAt *time.Time) (*draftauth.Session, error) {
	if f.session == nil || f.session.Session.ID != sessionID.String() {
		return nil, draftauth.ErrNotFound
	}
	f.session.Session.LastUsedAt = touchedAt.UTC()
	if renewedExpiresAt != nil {
		f.session.Session.ExpiresAt = renewedExpiresAt.UTC()
	}
	return &f.session.Session, nil
}

func (f *fakeSessionStore) RevokeAuthorSessionByTokenHash(ctx context.Context, tokenHash string) error {
	if f.session == nil || f.session.Session.TokenHash != tokenHash {
		return draftauth.ErrNotFound
	}
	now := time.Now().UTC()
	f.session.Session.RevokedAt = &now
	return nil
}

func testPublicFS() fs.FS {
	return fstest.MapFS{
		"index.html": {Data: []byte("<!doctype html><html><body>draft-review-shell</body></html>")},
		"assets/app.js": {
			Data: []byte("console.log('draft-review');"),
		},
	}
}

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

func TestHandleFrontendServesEmbeddedIndex(t *testing.T) {
	t.Parallel()

	handler := NewHandler(HandlerOptions{
		Version:   "test",
		StartedAt: time.Now().UTC(),
		AuthSettings: &draftauth.Settings{
			Mode: draftauth.AuthModeDev,
		},
		PublicFS: testPublicFS(),
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if body := recorder.Body.String(); !bytes.Contains([]byte(body), []byte("draft-review-shell")) {
		t.Fatalf("expected embedded shell, got %q", body)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Fatalf("expected html content type, got %q", contentType)
	}
	if cacheControl := recorder.Header().Get("Cache-Control"); cacheControl != "no-cache" {
		t.Fatalf("expected no-cache html header, got %q", cacheControl)
	}
}

func TestHandleFrontendServesEmbeddedAssets(t *testing.T) {
	t.Parallel()

	handler := NewHandler(HandlerOptions{
		Version:   "test",
		StartedAt: time.Now().UTC(),
		AuthSettings: &draftauth.Settings{
			Mode: draftauth.AuthModeDev,
		},
		PublicFS: testPublicFS(),
	})

	request := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if body := recorder.Body.String(); !bytes.Contains([]byte(body), []byte("draft-review")) {
		t.Fatalf("expected embedded asset body, got %q", body)
	}
}

func TestHandleFrontendFallsBackToIndexForReaderRoutes(t *testing.T) {
	t.Parallel()

	handler := NewHandler(HandlerOptions{
		Version:   "test",
		StartedAt: time.Now().UTC(),
		AuthSettings: &draftauth.Settings{
			Mode: draftauth.AuthModeDev,
		},
		PublicFS: testPublicFS(),
	})

	request := httptest.NewRequest(http.MethodGet, "/r/invite-123", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if body := recorder.Body.String(); !bytes.Contains([]byte(body), []byte("draft-review-shell")) {
		t.Fatalf("expected index fallback, got %q", body)
	}
	if cacheControl := recorder.Header().Get("Cache-Control"); cacheControl != "no-cache" {
		t.Fatalf("expected no-cache html header, got %q", cacheControl)
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

func TestHandleMeOIDCAuthenticatedViaOpaqueSession(t *testing.T) {
	t.Parallel()

	store := &fakeSessionStore{}
	manager, err := draftauth.NewSessionManager("test_session", "super-secret", "http://127.0.0.1:8080/auth/callback", 12*time.Hour, store)
	if err != nil {
		t.Fatalf("NewSessionManager returned error: %v", err)
	}

	user := &draftauth.User{
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
	store.session.User = *user

	request := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	for _, cookie := range writeRecorder.Result().Cookies() {
		request.AddCookie(cookie)
	}
	recorder := httptest.NewRecorder()

	handler := NewHandler(HandlerOptions{
		Version:        "test",
		StartedAt:      time.Now().UTC(),
		AuthSettings:   &draftauth.Settings{Mode: draftauth.AuthModeOIDC},
		SessionManager: manager,
	})

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
		t.Fatalf("expected authenticated oidc user")
	}
	if response.Data.Email != "alice@example.com" {
		t.Fatalf("expected email alice@example.com, got %q", response.Data.Email)
	}
}

type fakeAuthRepo struct {
	foundUser  *draftauth.User
	findErr    error
	createUser *draftauth.User
}

func (f *fakeAuthRepo) FindByAuthIdentity(ctx context.Context, issuer, subject string) (*draftauth.User, error) {
	return f.foundUser, f.findErr
}

func (f *fakeAuthRepo) CreateAuthenticatedUser(ctx context.Context, identity draftauth.AuthenticatedIdentity) (*draftauth.User, error) {
	return f.createUser, nil
}

func (f *fakeAuthRepo) UpdateAuthenticatedUser(ctx context.Context, userID uuid.UUID, identity draftauth.AuthenticatedIdentity) (*draftauth.User, error) {
	if f.foundUser != nil {
		return f.foundUser, nil
	}
	return f.createUser, nil
}

type fakeArticleRepo struct {
	lastOwnerUserID string
	deleteArticleID string
	deleteErr       error
}

func (f *fakeArticleRepo) ListArticles(ctx context.Context, ownerUserID string) ([]articles.Article, error) {
	f.lastOwnerUserID = ownerUserID
	return []articles.Article{{ID: "article-1", Title: "Owned Article"}}, nil
}

func (f *fakeArticleRepo) GetArticle(ctx context.Context, ownerUserID, id string) (*articles.Article, error) {
	f.lastOwnerUserID = ownerUserID
	return &articles.Article{ID: id, Title: "Owned Article"}, nil
}

func (f *fakeArticleRepo) CreateArticle(ctx context.Context, ownerUserID string, input articles.CreateArticleInput) (*articles.Article, error) {
	f.lastOwnerUserID = ownerUserID
	return &articles.Article{ID: "article-1", Title: input.Title, Author: input.Author}, nil
}

func (f *fakeArticleRepo) UpdateArticle(ctx context.Context, ownerUserID, id string, input articles.UpdateArticleInput) (*articles.Article, error) {
	f.lastOwnerUserID = ownerUserID
	return &articles.Article{ID: id, Title: "Owned Article"}, nil
}

func (f *fakeArticleRepo) CreateVersion(ctx context.Context, ownerUserID, id string, input articles.CreateVersionInput) (*articles.Article, error) {
	f.lastOwnerUserID = ownerUserID
	return &articles.Article{ID: id, Title: "Owned Article", Version: input.Label}, nil
}

func (f *fakeArticleRepo) DeleteArticle(ctx context.Context, ownerUserID, id string) error {
	f.lastOwnerUserID = ownerUserID
	f.deleteArticleID = id
	return f.deleteErr
}

type fakeReviewLinkRepo struct {
	link *reviewlinks.ResolvedLink
}

func (f *fakeReviewLinkRepo) ResetShareToken(ctx context.Context, ownerUserID, articleID string) (*reviewlinks.ShareLink, error) {
	return &reviewlinks.ShareLink{Token: "share-1", URL: "/r/share-1"}, nil
}

func (f *fakeReviewLinkRepo) CreateInvite(ctx context.Context, ownerUserID, articleID string, input reviewlinks.InviteInput) (*reviewlinks.Reader, error) {
	name := reviewlinks.DisplayNameFromInvite(input.DisplayName, input.Email, input.IdentityMode)
	return &reviewlinks.Reader{ID: "reader-1", Name: name, Email: input.Email, IdentityMode: input.IdentityMode}, nil
}

func (f *fakeReviewLinkRepo) ResolveToken(ctx context.Context, token string) (*reviewlinks.ResolvedLink, error) {
	if f.link != nil {
		return f.link, nil
	}
	return &reviewlinks.ResolvedLink{
		Token:            token,
		ArticleID:        "article-1",
		ArticleVersionID: "version-1",
		InviteID:         "invite-1",
		AllowAnonymous:   true,
		Reader: reviewlinks.ReaderIdentity{
			ID:   "reader-1",
			Name: "Reader One",
		},
		Article: reviewlinks.ReaderArticle{
			ID:      "article-1",
			Title:   "Shared Article",
			Author:  "Author",
			Version: "Draft 1",
			Intro:   "Intro",
			Sections: []reviewlinks.Section{
				{ID: "section-1", Title: "Section 1", Paragraphs: []string{"Paragraph 1"}},
			},
		},
	}, nil
}

type fakeReviewRepo struct{}

func (f *fakeReviewRepo) StartSession(ctx context.Context, link *reviewlinks.ResolvedLink, input reviews.StartSessionInput) (*reviews.StartSessionResult, error) {
	return &reviews.StartSessionResult{
		Session: &reviews.Session{
			ID:               "session-1",
			ArticleID:        link.Article.ID,
			ArticleVersionID: link.ArticleVersionID,
			ReaderID:         link.Reader.ID,
			ReaderName:       link.Reader.Name,
			ProgressPercent:  0,
			StartedAt:        time.Now().UTC(),
			LastActiveAt:     time.Now().UTC(),
		},
		Reader:  link.Reader,
		Article: link.Article,
	}, nil
}

func (f *fakeReviewRepo) RecordProgress(ctx context.Context, sessionID string, input reviews.ProgressInput) (*reviews.ProgressState, error) {
	return &reviews.ProgressState{SessionID: sessionID, ProgressPercent: 42, LastActiveAt: time.Now().UTC()}, nil
}

func (f *fakeReviewRepo) AddReaction(ctx context.Context, sessionID string, input reviews.ReactionInput) (*reviews.Reaction, error) {
	return &reviews.Reaction{ID: "reaction-1"}, nil
}

func (f *fakeReviewRepo) SubmitSummary(ctx context.Context, sessionID string, input reviews.SummaryInput) (*reviews.Summary, error) {
	return &reviews.Summary{SessionID: sessionID, SubmittedAt: time.Now().UTC()}, nil
}

type fakeAnalyticsRepo struct {
	lastOwnerUserID string
}

func (f *fakeAnalyticsRepo) ListReaders(ctx context.Context, ownerUserID, articleID string) ([]reviewlinks.Reader, error) {
	f.lastOwnerUserID = ownerUserID
	return []reviewlinks.Reader{{ID: "reader-1", Name: "Reader One", ArticleID: articleID}}, nil
}

func (f *fakeAnalyticsRepo) ListReactions(ctx context.Context, ownerUserID, articleID string) ([]reviews.Reaction, error) {
	f.lastOwnerUserID = ownerUserID
	return []reviews.Reaction{{ID: "reaction-1", ArticleID: articleID}}, nil
}

func (f *fakeAnalyticsRepo) GetArticleAnalytics(ctx context.Context, ownerUserID, articleID string) (*analytics.ArticleAnalytics, error) {
	f.lastOwnerUserID = ownerUserID
	return &analytics.ArticleAnalytics{ArticleID: articleID}, nil
}

func (f *fakeAnalyticsRepo) ListReaderDirectory(ctx context.Context, ownerUserID string) ([]analytics.ReaderContact, error) {
	f.lastOwnerUserID = ownerUserID
	return []analytics.ReaderContact{{Email: "reader@example.com"}}, nil
}

func TestHandleArticlesUsesAuthenticatedAuthor(t *testing.T) {
	t.Parallel()

	authRepo := &fakeAuthRepo{
		findErr: draftauth.ErrNotFound,
		createUser: &draftauth.User{
			ID:    "11111111-1111-1111-1111-111111111111",
			Email: "local-author@draft-review.local",
			Name:  "Development Author",
		},
	}
	articleRepo := &fakeArticleRepo{}

	handler := NewHandler(HandlerOptions{
		Version:   "test",
		StartedAt: time.Now().UTC(),
		AuthSettings: &draftauth.Settings{
			Mode:      draftauth.AuthModeDev,
			DevUserID: "local-author",
		},
		AuthService:    draftauth.NewService(authRepo),
		ArticleService: articles.NewService(articleRepo),
	})

	request := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	if articleRepo.lastOwnerUserID != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("expected article owner id to come from ensured auth user, got %q", articleRepo.lastOwnerUserID)
	}
}

func TestHandleDeleteArticleUsesAuthenticatedAuthor(t *testing.T) {
	t.Parallel()

	authRepo := &fakeAuthRepo{
		findErr: draftauth.ErrNotFound,
		createUser: &draftauth.User{
			ID:    "11111111-1111-1111-1111-111111111111",
			Email: "local-author@draft-review.local",
			Name:  "Development Author",
		},
	}
	articleRepo := &fakeArticleRepo{}

	handler := NewHandler(HandlerOptions{
		Version:   "test",
		StartedAt: time.Now().UTC(),
		AuthSettings: &draftauth.Settings{
			Mode:      draftauth.AuthModeDev,
			DevUserID: "local-author",
		},
		AuthService:    draftauth.NewService(authRepo),
		ArticleService: articles.NewService(articleRepo),
	})

	request := httptest.NewRequest(http.MethodDelete, "/api/articles/article-123", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", recorder.Code)
	}

	if articleRepo.lastOwnerUserID != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("expected article owner id to come from ensured auth user, got %q", articleRepo.lastOwnerUserID)
	}
	if articleRepo.deleteArticleID != "article-123" {
		t.Fatalf("expected delete to target article-123, got %q", articleRepo.deleteArticleID)
	}
}

func TestHandleResolveReviewLinkHidesInternalFields(t *testing.T) {
	t.Parallel()

	handler := NewHandler(HandlerOptions{
		Version:           "test",
		StartedAt:         time.Now().UTC(),
		AuthSettings:      &draftauth.Settings{Mode: draftauth.AuthModeDev, DevUserID: "local-author"},
		ReviewLinkService: reviewlinks.NewService(&fakeReviewLinkRepo{}),
	})

	request := httptest.NewRequest(http.MethodGet, "/api/r/token-1", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	var payload map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode payload: %v", err)
	}

	if _, ok := payload["reader"]; !ok {
		t.Fatalf("expected reader payload")
	}
	if _, ok := payload["article"]; !ok {
		t.Fatalf("expected article payload")
	}
	for _, field := range []string{"articleId", "articleVersionId", "inviteId", "allowAnonymous", "token"} {
		if _, ok := payload[field]; ok {
			t.Fatalf("did not expect internal field %q in payload", field)
		}
	}
}

func TestHandleArticleReadersUsesAnalyticsService(t *testing.T) {
	t.Parallel()

	authRepo := &fakeAuthRepo{
		findErr: draftauth.ErrNotFound,
		createUser: &draftauth.User{
			ID:    "11111111-1111-1111-1111-111111111111",
			Email: "local-author@draft-review.local",
			Name:  "Development Author",
		},
	}
	analyticsRepo := &fakeAnalyticsRepo{}

	handler := NewHandler(HandlerOptions{
		Version:   "test",
		StartedAt: time.Now().UTC(),
		AuthSettings: &draftauth.Settings{
			Mode:      draftauth.AuthModeDev,
			DevUserID: "local-author",
		},
		AuthService:      draftauth.NewService(authRepo),
		AnalyticsService: analytics.NewService(analyticsRepo),
	})

	request := httptest.NewRequest(http.MethodGet, "/api/articles/article-1/readers", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if analyticsRepo.lastOwnerUserID != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("expected analytics owner id to come from ensured auth user, got %q", analyticsRepo.lastOwnerUserID)
	}
}

func TestHandleReviewProgressRejectsInvalidPercent(t *testing.T) {
	t.Parallel()

	handler := NewHandler(HandlerOptions{
		Version:       "test",
		StartedAt:     time.Now().UTC(),
		AuthSettings:  &draftauth.Settings{Mode: draftauth.AuthModeDev, DevUserID: "local-author"},
		ReviewService: reviews.NewService(&fakeReviewRepo{}),
	})

	request := httptest.NewRequest(http.MethodPost, "/api/reviews/session-1/progress", bytes.NewBufferString(`{"sectionId":"section-1","progressPercent":101}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestHandleFrontendProxiesToDevServer(t *testing.T) {
	t.Parallel()

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = io.WriteString(w, "proxied:"+r.URL.Path)
	}))
	defer upstream.Close()

	handler := NewHandler(HandlerOptions{
		Version:             "test",
		StartedAt:           time.Now().UTC(),
		AuthSettings:        &draftauth.Settings{Mode: draftauth.AuthModeDev, DevUserID: "local-author"},
		FrontendDevProxyURL: upstream.URL,
	})

	request := httptest.NewRequest(http.MethodGet, "/r/token-1", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if got := recorder.Body.String(); got != "proxied:/r/token-1" {
		t.Fatalf("expected proxied body, got %q", got)
	}
}
