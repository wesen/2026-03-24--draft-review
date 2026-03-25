package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	draftdb "github.com/go-go-golems/draft-review/pkg/db"
)

type Options struct {
	Host                string
	Port                int
	Version             string
	Database            *draftdb.DB
	FrontendDevProxyURL string
}

type HandlerOptions struct {
	Version             string
	StartedAt           time.Time
	Database            *draftdb.DB
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
	version   string
	startedAt time.Time
	database  *draftdb.DB
}

func NewHTTPServer(_ context.Context, options Options) (*http.Server, error) {
	return &http.Server{
		Addr: fmt.Sprintf("%s:%d", options.Host, options.Port),
		Handler: NewHandler(HandlerOptions{
			Version:             options.Version,
			StartedAt:           time.Now().UTC(),
			Database:            options.Database,
			FrontendDevProxyURL: options.FrontendDevProxyURL,
		}),
		ReadHeaderTimeout: 10 * time.Second,
	}, nil
}

func NewHandler(options HandlerOptions) http.Handler {
	h := &appHandler{
		version:   options.Version,
		startedAt: options.StartedAt,
		database:  options.Database,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, apiEnvelope{Data: map[string]string{"status": "ok"}})
	})
	mux.HandleFunc("GET /api/info", h.handleInfo)

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

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
