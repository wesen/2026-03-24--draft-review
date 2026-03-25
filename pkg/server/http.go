package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/go-go-golems/draft-review/pkg/analytics"
	"github.com/go-go-golems/draft-review/pkg/articles"
	draftauth "github.com/go-go-golems/draft-review/pkg/auth"
	draftdb "github.com/go-go-golems/draft-review/pkg/db"
	"github.com/go-go-golems/draft-review/pkg/reviewlinks"
	"github.com/go-go-golems/draft-review/pkg/reviews"
	draftweb "github.com/go-go-golems/draft-review/pkg/web"
)

type Options struct {
	Host                string
	Port                int
	Version             string
	AuthSettings        *draftauth.Settings
	AuthService         *draftauth.Service
	Database            *draftdb.DB
	ArticleService      *articles.Service
	ReviewLinkService   *reviewlinks.Service
	ReviewService       *reviews.Service
	AnalyticsService    *analytics.Service
	FrontendDevProxyURL string
}

type HandlerOptions struct {
	Version             string
	StartedAt           time.Time
	AuthSettings        *draftauth.Settings
	AuthService         *draftauth.Service
	SessionManager      *draftauth.SessionManager
	WebAuth             draftauth.WebHandler
	PublicFS            fs.FS
	Database            *draftdb.DB
	ArticleService      *articles.Service
	ReviewLinkService   *reviewlinks.Service
	ReviewService       *reviews.Service
	AnalyticsService    *analytics.Service
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
	version           string
	startedAt         time.Time
	authSettings      *draftauth.Settings
	authService       *draftauth.Service
	sessionManager    *draftauth.SessionManager
	publicFS          fs.FS
	database          *draftdb.DB
	articleService    *articles.Service
	reviewLinkService *reviewlinks.Service
	reviewService     *reviews.Service
	analyticsService  *analytics.Service
	frontendProxy     *httputil.ReverseProxy
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

	authService := options.AuthService
	if authService == nil && options.Database != nil && options.Database.Pool() != nil {
		authService = draftauth.NewService(draftauth.NewPostgresRepository(options.Database.Pool()))
	}

	articleService := options.ArticleService
	if articleService == nil && options.Database != nil && options.Database.Pool() != nil {
		articleService = articles.NewService(articles.NewPostgresRepository(options.Database.Pool()))
	}

	reviewLinkService := options.ReviewLinkService
	if reviewLinkService == nil && options.Database != nil && options.Database.Pool() != nil {
		reviewLinkService = reviewlinks.NewService(reviewlinks.NewPostgresRepository(options.Database.Pool()))
	}

	reviewService := options.ReviewService
	if reviewService == nil && options.Database != nil && options.Database.Pool() != nil {
		reviewService = reviews.NewService(reviews.NewPostgresRepository(options.Database.Pool()))
	}

	analyticsService := options.AnalyticsService
	if analyticsService == nil && options.Database != nil && options.Database.Pool() != nil {
		analyticsService = analytics.NewService(analytics.NewPostgresRepository(options.Database.Pool()))
	}

	return &http.Server{
		Addr: fmt.Sprintf("%s:%d", options.Host, options.Port),
		Handler: NewHandler(HandlerOptions{
			Version:             options.Version,
			StartedAt:           time.Now().UTC(),
			AuthSettings:        authSettings,
			AuthService:         authService,
			SessionManager:      sessionManager,
			WebAuth:             webAuth,
			PublicFS:            draftweb.PublicFS,
			Database:            options.Database,
			ArticleService:      articleService,
			ReviewLinkService:   reviewLinkService,
			ReviewService:       reviewService,
			AnalyticsService:    analyticsService,
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

	publicFS := options.PublicFS
	if publicFS == nil {
		publicFS = draftweb.PublicFS
	}

	frontendProxy := newFrontendDevProxy(options.FrontendDevProxyURL)

	h := &appHandler{
		version:           options.Version,
		startedAt:         options.StartedAt,
		authSettings:      authSettings,
		authService:       options.AuthService,
		sessionManager:    options.SessionManager,
		publicFS:          publicFS,
		database:          options.Database,
		articleService:    options.ArticleService,
		reviewLinkService: options.ReviewLinkService,
		reviewService:     options.ReviewService,
		analyticsService:  options.AnalyticsService,
		frontendProxy:     frontendProxy,
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
	mux.HandleFunc("POST /api/articles/{id}/versions", h.handleCreateArticleVersion)
	mux.HandleFunc("POST /api/articles/{id}/share-token", h.handleResetShareToken)
	mux.HandleFunc("POST /api/articles/{id}/invite", h.handleInviteReader)
	mux.HandleFunc("GET /api/articles/{id}/readers", h.handleArticleReaders)
	mux.HandleFunc("GET /api/articles/{id}/reactions", h.handleArticleReactions)
	mux.HandleFunc("GET /api/articles/{id}/analytics", h.handleArticleAnalytics)
	mux.HandleFunc("GET /api/articles/{id}/feedback", h.handleArticleFeedback)
	mux.HandleFunc("POST /api/articles/{id}/export", h.handleArticleExport)
	mux.HandleFunc("GET /api/readers", h.handleReadersDirectory)
	mux.HandleFunc("GET /api/r/{token}", h.handleResolveReviewLink)
	mux.HandleFunc("POST /api/r/{token}/start", h.handleStartReview)
	mux.HandleFunc("POST /api/reviews/{sessionId}/progress", h.handleReviewProgress)
	mux.HandleFunc("POST /api/reviews/{sessionId}/reactions", h.handleReviewReaction)
	mux.HandleFunc("POST /api/reviews/{sessionId}/summary", h.handleReviewSummary)
	if options.WebAuth != nil {
		mux.HandleFunc("GET /auth/login", options.WebAuth.HandleLogin)
		mux.HandleFunc("GET /auth/callback", options.WebAuth.HandleCallback)
		mux.HandleFunc("GET /auth/logout", options.WebAuth.HandleLogout)
		mux.HandleFunc("GET /auth/logout/callback", options.WebAuth.HandleLogoutCallback)
	}
	mux.HandleFunc("/", h.handleFrontend)

	return mux
}

func newFrontendDevProxy(rawURL string) *httputil.ReverseProxy {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return nil
	}

	target, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, err error) {
		writeError(w, http.StatusBadGateway, fmt.Sprintf("frontend dev proxy failed: %v", err))
	}
	return proxy
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

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	result, err := h.articleService.ListArticles(r.Context(), author)
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

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	var input articles.CreateArticleInput
	if err := decodeJSONBody(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.articleService.CreateArticle(r.Context(), author, input)
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

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	result, err := h.articleService.GetArticle(r.Context(), author, r.PathValue("id"))
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

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	var input articles.UpdateArticleInput
	if err := decodeJSONBody(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.articleService.UpdateArticle(r.Context(), author, r.PathValue("id"), input)
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

func (h *appHandler) handleCreateArticleVersion(w http.ResponseWriter, r *http.Request) {
	if h.articleService == nil {
		writeError(w, http.StatusServiceUnavailable, "article service is not configured")
		return
	}

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	var input articles.CreateVersionInput
	if err := decodeJSONBody(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.articleService.CreateVersion(r.Context(), author, r.PathValue("id"), input)
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

	writeJSON(w, http.StatusCreated, result)
}

func (h *appHandler) handleResetShareToken(w http.ResponseWriter, r *http.Request) {
	if h.reviewLinkService == nil {
		writeError(w, http.StatusServiceUnavailable, "review link service is not configured")
		return
	}

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	result, err := h.reviewLinkService.ResetShareToken(r.Context(), author.ID, r.PathValue("id"))
	if err != nil {
		switch {
		case errors.Is(err, reviewlinks.ErrNotFound):
			http.NotFound(w, r)
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleInviteReader(w http.ResponseWriter, r *http.Request) {
	if h.reviewLinkService == nil {
		writeError(w, http.StatusServiceUnavailable, "review link service is not configured")
		return
	}

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	var input reviewlinks.InviteInput
	if err := decodeJSONBody(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.reviewLinkService.CreateInvite(r.Context(), author.ID, r.PathValue("id"), input)
	if err != nil {
		switch {
		case errors.Is(err, reviewlinks.ErrNotFound):
			http.NotFound(w, r)
		case reviewlinks.IsValidationError(err):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *appHandler) handleArticleReaders(w http.ResponseWriter, r *http.Request) {
	if h.analyticsService == nil {
		writeJSON(w, http.StatusOK, []reviewlinks.Reader{})
		return
	}

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	result, err := h.analyticsService.ListReaders(r.Context(), author.ID, r.PathValue("id"))
	if err != nil {
		switch {
		case errors.Is(err, analytics.ErrNotFound):
			http.NotFound(w, r)
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleArticleReactions(w http.ResponseWriter, r *http.Request) {
	if h.analyticsService == nil {
		writeJSON(w, http.StatusOK, []reviews.Reaction{})
		return
	}

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	result, err := h.analyticsService.ListReactions(r.Context(), author.ID, r.PathValue("id"))
	if err != nil {
		switch {
		case errors.Is(err, analytics.ErrNotFound):
			http.NotFound(w, r)
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleArticleAnalytics(w http.ResponseWriter, r *http.Request) {
	if h.analyticsService == nil {
		writeError(w, http.StatusServiceUnavailable, "analytics service is not configured")
		return
	}

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	result, err := h.analyticsService.GetArticleAnalytics(r.Context(), author.ID, r.PathValue("id"))
	if err != nil {
		switch {
		case errors.Is(err, analytics.ErrNotFound):
			http.NotFound(w, r)
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleArticleFeedback(w http.ResponseWriter, r *http.Request) {
	h.handleArticleReactions(w, r)
}

func (h *appHandler) handleArticleExport(w http.ResponseWriter, r *http.Request) {
	if h.articleService == nil {
		writeError(w, http.StatusServiceUnavailable, "article service is not configured")
		return
	}

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	article, err := h.articleService.GetArticle(r.Context(), author, r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]string{
		"status":    "stub",
		"articleId": article.ID,
		"message":   "export generation is not implemented yet",
	})
}

func (h *appHandler) handleReadersDirectory(w http.ResponseWriter, r *http.Request) {
	if h.analyticsService == nil {
		writeJSON(w, http.StatusOK, []analytics.ReaderContact{})
		return
	}

	author, ok := h.requireCurrentAuthor(w, r)
	if !ok {
		return
	}

	result, err := h.analyticsService.ListReaderDirectory(r.Context(), author.ID)
	if err != nil {
		switch {
		case errors.Is(err, analytics.ErrNotFound):
			http.NotFound(w, r)
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleResolveReviewLink(w http.ResponseWriter, r *http.Request) {
	if h.reviewLinkService == nil {
		writeError(w, http.StatusServiceUnavailable, "review link service is not configured")
		return
	}

	result, err := h.reviewLinkService.ResolveToken(r.Context(), r.PathValue("token"))
	if err != nil {
		switch {
		case errors.Is(err, reviewlinks.ErrNotFound):
			http.NotFound(w, r)
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleStartReview(w http.ResponseWriter, r *http.Request) {
	if h.reviewLinkService == nil || h.reviewService == nil {
		writeError(w, http.StatusServiceUnavailable, "review services are not configured")
		return
	}

	link, err := h.reviewLinkService.ResolveToken(r.Context(), r.PathValue("token"))
	if err != nil {
		switch {
		case errors.Is(err, reviewlinks.ErrNotFound):
			http.NotFound(w, r)
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	var input reviews.StartSessionInput
	if err := decodeJSONBody(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.reviewService.StartSession(r.Context(), link, input)
	if err != nil {
		switch {
		case errors.Is(err, reviewlinks.ErrNotFound), errors.Is(err, reviews.ErrNotFound):
			http.NotFound(w, r)
		case reviews.IsValidationError(err):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *appHandler) handleReviewProgress(w http.ResponseWriter, r *http.Request) {
	if h.reviewService == nil {
		writeError(w, http.StatusServiceUnavailable, "review service is not configured")
		return
	}

	var input reviews.ProgressInput
	if err := decodeJSONBody(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.reviewService.RecordProgress(r.Context(), r.PathValue("sessionId"), input)
	if err != nil {
		switch {
		case errors.Is(err, reviews.ErrNotFound):
			http.NotFound(w, r)
		case reviews.IsValidationError(err):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleReviewReaction(w http.ResponseWriter, r *http.Request) {
	if h.reviewService == nil {
		writeError(w, http.StatusServiceUnavailable, "review service is not configured")
		return
	}

	var input reviews.ReactionInput
	if err := decodeJSONBody(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.reviewService.AddReaction(r.Context(), r.PathValue("sessionId"), input)
	if err != nil {
		switch {
		case errors.Is(err, reviews.ErrNotFound):
			http.NotFound(w, r)
		case reviews.IsValidationError(err):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *appHandler) handleReviewSummary(w http.ResponseWriter, r *http.Request) {
	if h.reviewService == nil {
		writeError(w, http.StatusServiceUnavailable, "review service is not configured")
		return
	}

	var input reviews.SummaryInput
	if err := decodeJSONBody(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.reviewService.SubmitSummary(r.Context(), r.PathValue("sessionId"), input)
	if err != nil {
		switch {
		case errors.Is(err, reviews.ErrNotFound):
			http.NotFound(w, r)
		case reviews.IsValidationError(err):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *appHandler) handleFrontend(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/auth/") {
		http.NotFound(w, r)
		return
	}

	if h.frontendProxy != nil {
		h.frontendProxy.ServeHTTP(w, r)
		return
	}
	if h.publicFS == nil {
		http.NotFound(w, r)
		return
	}

	requestPath := strings.TrimPrefix(r.URL.Path, "/")
	if requestPath == "" {
		requestPath = "index.html"
	}

	if fileInfo, err := fs.Stat(h.publicFS, requestPath); err == nil && !fileInfo.IsDir() {
		http.FileServer(http.FS(h.publicFS)).ServeHTTP(w, r)
		return
	}

	index, err := h.publicFS.Open("index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer func() { _ = index.Close() }()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.Copy(w, index)
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

func (h *appHandler) requireCurrentAuthor(w http.ResponseWriter, r *http.Request) (*draftauth.User, bool) {
	identity, ok := h.currentIdentity(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return nil, false
	}
	if h.authService == nil {
		writeError(w, http.StatusInternalServerError, "auth service is not configured")
		return nil, false
	}

	user, err := h.authService.EnsureAuthenticatedUser(r.Context(), identity)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to resolve authenticated author")
		return nil, false
	}

	return user, true
}

func (h *appHandler) currentIdentity(r *http.Request) (draftauth.AuthenticatedIdentity, bool) {
	claims, ok := h.currentClaims(r)
	if !ok {
		return draftauth.AuthenticatedIdentity{}, false
	}

	return draftauth.AuthenticatedIdentity{
		Issuer:        claims.Issuer,
		Subject:       claims.Subject,
		Email:         claims.Email,
		DisplayName:   claims.DisplayName,
		EmailVerified: claims.EmailVerified,
	}, true
}

func (h *appHandler) currentClaims(r *http.Request) (*draftauth.SessionClaims, bool) {
	switch h.authSettings.Mode {
	case draftauth.AuthModeDev:
		return &draftauth.SessionClaims{
			Issuer:            "dev",
			Subject:           h.authSettings.DevUserID,
			Email:             "local-author@draft-review.local",
			EmailVerified:     true,
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
