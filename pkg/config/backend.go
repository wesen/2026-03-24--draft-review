package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/pkg/errors"
)

const (
	BackendSectionSlug    = "backend"
	DefaultAutoMigrate    = true
	DefaultAppBaseURL     = "http://127.0.0.1:8080"
	DefaultMediaRoot      = ".draft-review/media"
	DefaultMaxUploadBytes = 10 * 1024 * 1024
)

type BackendSettings struct {
	AutoMigrate         bool   `glazed:"auto-migrate"`
	AppBaseURL          string `glazed:"app-base-url"`
	FrontendDevProxyURL string `glazed:"frontend-dev-proxy-url"`
	MediaRoot           string `glazed:"media-root"`
	MaxUploadBytes      int64  `glazed:"max-upload-bytes"`
}

func NewBackendSection() (schema.Section, error) {
	return schema.NewSection(
		BackendSectionSlug,
		"Backend Settings",
		schema.WithFields(
			fields.New(
				"auto-migrate",
				fields.TypeBool,
				fields.WithHelp("Apply embedded SQL migrations automatically on startup when a database is configured"),
				fields.WithDefault(envBoolOr("DRAFT_REVIEW_AUTO_MIGRATE", DefaultAutoMigrate)),
			),
			fields.New(
				"app-base-url",
				fields.TypeString,
				fields.WithHelp("Base URL used to build internal links and future invite URLs"),
				fields.WithDefault(envOr("DRAFT_REVIEW_APP_BASE_URL", DefaultAppBaseURL)),
			),
			fields.New(
				"frontend-dev-proxy-url",
				fields.TypeString,
				fields.WithHelp("Optional frontend dev-server origin to proxy non-API browser routes to during local development"),
				fields.WithDefault(strings.TrimSpace(os.Getenv("DRAFT_REVIEW_FRONTEND_DEV_PROXY_URL"))),
			),
			fields.New(
				"media-root",
				fields.TypeString,
				fields.WithHelp("Filesystem root for persisted uploaded article media; use a persistent mounted volume in hosted environments"),
				fields.WithDefault(envOr("DRAFT_REVIEW_MEDIA_ROOT", DefaultMediaRoot)),
			),
			fields.New(
				"max-upload-bytes",
				fields.TypeInteger,
				fields.WithHelp("Maximum allowed size in bytes for uploaded article images"),
				fields.WithDefault(envInt64Or("DRAFT_REVIEW_MAX_UPLOAD_BYTES", DefaultMaxUploadBytes)),
			),
		),
	)
}

func LoadBackendSettings(parsedValues *values.Values) (*BackendSettings, error) {
	if parsedValues == nil {
		return nil, errors.New("parsed values are nil")
	}

	settings := &BackendSettings{}
	if err := parsedValues.DecodeSectionInto(BackendSectionSlug, settings); err != nil {
		return nil, errors.Wrap(err, "failed to decode backend settings")
	}

	return NormalizeBackendSettings(settings), nil
}

func NormalizeBackendSettings(settings *BackendSettings) *BackendSettings {
	if settings == nil {
		settings = &BackendSettings{}
	}
	settings.AppBaseURL = strings.TrimSpace(settings.AppBaseURL)
	if settings.AppBaseURL == "" {
		settings.AppBaseURL = DefaultAppBaseURL
	}
	settings.FrontendDevProxyURL = strings.TrimSpace(settings.FrontendDevProxyURL)
	settings.MediaRoot = strings.TrimSpace(settings.MediaRoot)
	if settings.MediaRoot == "" {
		settings.MediaRoot = DefaultMediaRoot
	}
	if settings.MaxUploadBytes <= 0 {
		settings.MaxUploadBytes = DefaultMaxUploadBytes
	}
	return settings
}

func envOr(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envBoolOr(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}

func envInt64Or(key string, fallback int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	var parsed int64
	_, err := fmt.Sscan(value, &parsed)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}
