package config

import (
	"os"
	"strings"

	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/pkg/errors"
)

const (
	BackendSectionSlug = "backend"
	DefaultAutoMigrate = true
	DefaultAppBaseURL  = "http://127.0.0.1:8080"
)

type BackendSettings struct {
	AutoMigrate         bool   `glazed:"auto-migrate"`
	AppBaseURL          string `glazed:"app-base-url"`
	FrontendDevProxyURL string `glazed:"frontend-dev-proxy-url"`
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
