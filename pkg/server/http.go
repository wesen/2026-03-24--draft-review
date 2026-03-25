package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-go-golems/draft-review/pkg/articles"
	draftdb "github.com/go-go-golems/draft-review/pkg/db"
)

type Options struct {
	Host                string
	Port                int
	Version             string
	Database            *draftdb.DB
	ArticleService      *articles.Service
	FrontendDevProxyURL string
}

type HandlerOptions struct {
	Version             string
	StartedAt           time.Time
	Database            *draftdb.DB
	ArticleService      *articles.Service
	FrontendDevProxyURL string
}

type infoResponse struct {
	Service            string    `json:"service"`
	Version            string    `json:"version"`
	StartedAt          time.Time `json:"startedAt"`
	DatabaseConfigured bool      `json:"databaseConfigured"`
}

type apiEnvelope struct {
	Data any `json:"data,omitempty"`
}

type appHandler struct {
	version        string
	startedAt      time.Time
	database       *draftdb.DB
	articleService *articles.Service
}

func NewHTTPServer(_ context.Context, options Options) (*http.Server, error) {
	articleService := options.ArticleService
	if articleService == nil && options.Database != nil && options.Database.Pool() != nil {
		articleService = articles.NewService(articles.NewPostgresRepository(options.Database.Pool()))
	}

	return &http.Server{
		Addr: fmt.Sprintf("%s:%d", options.Host, options.Port),
		Handler: NewHandler(HandlerOptions{
			Version:             options.Version,
			StartedAt:           time.Now().UTC(),
			Database:            options.Database,
			ArticleService:      articleService,
			FrontendDevProxyURL: options.FrontendDevProxyURL,
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}, nil
}

func NewHandler(options HandlerOptions) http.Handler {
	h := &appHandler{
		version:        options.Version,
		startedAt:      options.StartedAt,
		database:       options.Database,
		articleService: options.ArticleService,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, apiEnvelope{Data: map[string]string{"status": "ok"}})
	})
	mux.HandleFunc("GET /api/info", h.handleInfo)
	mux.HandleFunc("GET /api/articles", h.handleArticles)
	mux.HandleFunc("GET /api/articles/{id}", h.handleArticle)
	mux.HandleFunc("GET /api/articles/{id}/readers", h.handleArticleReaders)
	mux.HandleFunc("GET /api/articles/{id}/reactions", h.handleArticleReactions)

	return mux
}

func (h *appHandler) handleInfo(w http.ResponseWriter, _ *http.Request) {
	response := infoResponse{
		Service:            "draft-review",
		Version:            h.version,
		StartedAt:          h.startedAt,
		DatabaseConfigured: h.database != nil,
	}
	writeJSON(w, http.StatusOK, apiEnvelope{Data: response})
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
