package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-go-golems/draft-review/pkg/articles"
	draftauth "github.com/go-go-golems/draft-review/pkg/auth"
	"github.com/google/uuid"
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
