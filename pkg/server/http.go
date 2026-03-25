package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-go-golems/draft-review/pkg/articles"
	draftauth "github.com/go-go-golems/draft-review/pkg/auth"
	draftdb "github.com/go-go-golems/draft-review/pkg/db"
)

type Options struct {
	Host                string
	Port                int
	Version             string
	AuthSettings        *draftauth.Settings
	Database            *draftdb.DB
	ArticleService      *articles.Service
	FrontendDevProxyURL string
}

type HandlerOptions struct {
	Version             string
	StartedAt           time.Time
	AuthSettings        *draftauth.Settings
	SessionManager      *draftauth.SessionManager
	WebAuth             draftauth.WebHandler
	Database            *draftdb.DB
	ArticleService      *articles.Service
	FrontendDevProxyURL string
}

type infoResponse struct {
	Service            string    `json:"service"`
	Version            string    `json:"version"`
	StartedAt          time.Time `json:"startedAt"`
	AuthMode           string    `json:"authMode"`
	IssuerURL          string    `json:"issuerUrl,omitempty"`
	ClientID           string    `json:"clientId,omitempty"`
	LoginPath          string    `json:"loginPath,omitempty"`
	LogoutPath         string    `json:"logoutPath,omitempty"`
	CallbackPath       string    `json:"callbackPath,omitempty"`
	DatabaseConfigured bool      `json:"databaseConfigured"`
}

type apiEnvelope struct {
	Data any `json:"data,omitempty"`
}

type appHandler struct {
	version        string
	startedAt      time.Time
	authSettings   *draftauth.Settings
	sessionManager *draftauth.SessionManager
	database       *draftdb.DB
	articleService *articles.Service
}

func NewHTTPServer(ctx context.Context, options Options) (*http.Server, error) {
	var (
		sessionManager *draftauth.SessionManager
		webAuth        draftauth.WebHandler
		err            error
	)

	authSettings := options.AuthSettings
	if authSettings == nil {
		authSettings = &draftauth.Settings{Mode: draftauth.AuthModeDev, DevUserID: "local-author"}
	}

	if authSettings.Mode == draftauth.AuthModeOIDC {
		sessionManager, err = draftauth.NewSessionManager(
			authSettings.SessionCookieName,
			authSettings.SessionSecret,
			authSettings.OIDCRedirectURL,
		)
		if err != nil {
			return nil, err
		}

		webAuth, err = draftauth.NewOIDCAuthenticator(ctx, authSettings, sessionManager)
		if err != nil {
			return nil, err
		}
	}

	articleService := options.ArticleService
	if articleService == nil && options.Database != nil && options.Database.Pool() != nil {
		articleService = articles.NewService(articles.NewPostgresRepository(options.Database.Pool()))
	}

	return &http.Server{
		Addr: fmt.Sprintf("%s:%d", options.Host, options.Port),
		Handler: NewHandler(HandlerOptions{
			Version:             options.Version,
			StartedAt:           time.Now().UTC(),
			AuthSettings:        authSettings,
			SessionManager:      sessionManager,
			WebAuth:             webAuth,
			Database:            options.Database,
			ArticleService:      articleService,
			FrontendDevProxyURL: options.FrontendDevProxyURL,
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}, nil
}

func NewHandler(options HandlerOptions) http.Handler {
	authSettings := options.AuthSettings
	if authSettings == nil {
		authSettings = &draftauth.Settings{Mode: draftauth.AuthModeDev, DevUserID: "local-author"}
	}

	h := &appHandler{
		version:        options.Version,
		startedAt:      options.StartedAt,
		authSettings:   authSettings,
		sessionManager: options.SessionManager,
		database:       options.Database,
		articleService: options.ArticleService,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, apiEnvelope{Data: map[string]string{"status": "ok"}})
	})
	mux.HandleFunc("GET /api/info", h.handleInfo)
	mux.HandleFunc("GET /api/me", h.handleMe)
	mux.HandleFunc("GET /api/articles", h.handleArticles)
	mux.HandleFunc("POST /api/articles", h.handleCreateArticle)
	mux.HandleFunc("GET /api/articles/{id}", h.handleArticle)
	mux.HandleFunc("PATCH /api/articles/{id}", h.handleUpdateArticle)
	mux.HandleFunc("GET /api/articles/{id}/readers", h.handleArticleReaders)
	mux.HandleFunc("GET /api/articles/{id}/reactions", h.handleArticleReactions)
	if options.WebAuth != nil {
		mux.HandleFunc("GET /auth/login", options.WebAuth.HandleLogin)
		mux.HandleFunc("GET /auth/callback", options.WebAuth.HandleCallback)
		mux.HandleFunc("GET /auth/logout", options.WebAuth.HandleLogout)
		mux.HandleFunc("GET /auth/logout/callback", options.WebAuth.HandleLogoutCallback)
	}

	return mux
}

func (h *appHandler) handleInfo(w http.ResponseWriter, _ *http.Request) {
	response := infoResponse{
		Service:            "draft-review",
		Version:            h.version,
		StartedAt:          h.startedAt,
		AuthMode:           h.authSettings.Mode,
		IssuerURL:          h.authSettings.OIDCIssuerURL,
		ClientID:           h.authSettings.OIDCClientID,
		LoginPath:          "/auth/login",
		LogoutPath:         "/auth/logout",
		CallbackPath:       "/auth/callback",
		DatabaseConfigured: h.database != nil,
	}
	writeJSON(w, http.StatusOK, apiEnvelope{Data: response})
}

func (h *appHandler) handleMe(w http.ResponseWriter, r *http.Request) {
	user, ok := h.currentUser(r)
	if !ok {
		writeJSON(w, http.StatusOK, apiEnvelope{Data: draftauth.UserInfo{
			Authenticated: false,
			AuthMode:      h.authSettings.Mode,
		}})
		return
	}

	writeJSON(w, http.StatusOK, apiEnvelope{Data: user})
}

func (h *appHandler) handleArticles(w http.ResponseWriter, r *http.Request) {
	if h.articleService == nil {
		writeJSON(w, http.StatusOK, []articles.Article{})
		return
	}

	result, err := h.articleService.ListArticles(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleCreateArticle(w http.ResponseWriter, r *http.Request) {
	if h.articleService == nil {
		writeError(w, http.StatusServiceUnavailable, "article service is not configured")
		return
	}

	var input articles.CreateArticleInput
	if err := decodeJSONBody(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.articleService.CreateArticle(r.Context(), input)
	if err != nil {
		if articles.IsValidationError(err) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *appHandler) handleArticle(w http.ResponseWriter, r *http.Request) {
	if h.articleService == nil {
		http.NotFound(w, r)
		return
	}

	result, err := h.articleService.GetArticle(r.Context(), r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleUpdateArticle(w http.ResponseWriter, r *http.Request) {
	if h.articleService == nil {
		writeError(w, http.StatusServiceUnavailable, "article service is not configured")
		return
	}

	var input articles.UpdateArticleInput
	if err := decodeJSONBody(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.articleService.UpdateArticle(r.Context(), r.PathValue("id"), input)
	if err != nil {
		switch {
		case errors.Is(err, articles.ErrNotFound):
			http.NotFound(w, r)
		case articles.IsValidationError(err):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleArticleReaders(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, []any{})
}

func (h *appHandler) handleArticleReactions(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, []any{})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func decodeJSONBody(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}

	var trailing any
	if err := decoder.Decode(&trailing); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("request body must contain a single JSON object: %w", err)
	}

	return nil
}

func (h *appHandler) currentUser(r *http.Request) (*draftauth.UserInfo, bool) {
	claims, ok := h.currentClaims(r)
	if !ok {
		return nil, false
	}

	user := claims.UserInfo(h.authSettings.Mode)
	return &user, true
}

func (h *appHandler) currentClaims(r *http.Request) (*draftauth.SessionClaims, bool) {
	switch h.authSettings.Mode {
	case draftauth.AuthModeDev:
		return &draftauth.SessionClaims{
			Issuer:            "dev",
			Subject:           h.authSettings.DevUserID,
			PreferredUsername: h.authSettings.DevUserID,
			DisplayName:       "Development Author",
			IssuedAt:          time.Now().UTC(),
			ExpiresAt:         time.Now().UTC().Add(24 * time.Hour),
		}, true
	case draftauth.AuthModeOIDC:
		if h.sessionManager == nil {
			return nil, false
		}
		claims, err := h.sessionManager.ReadSession(r)
		if err != nil {
			return nil, false
		}
		return claims, true
	default:
		return nil, false
	}
}
