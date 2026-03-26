package auth

import (
	"os"
	"strings"
	"time"

	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/pkg/errors"
)

const (
	AuthSectionSlug          = "auth"
	AuthModeDev              = "dev"
	AuthModeOIDC             = "oidc"
	DefaultSessionCookieName = "draft_review_session"
)

type Settings struct {
	Mode                    string   `glazed:"auth-mode"`
	DevUserID               string   `glazed:"auth-dev-user-id"`
	SessionCookieName       string   `glazed:"auth-session-cookie-name"`
	SessionSecret           string   `glazed:"auth-session-secret"`
	SessionTTL              string   `glazed:"auth-session-ttl"`
	SessionSlidingRenewal   bool     `glazed:"auth-session-sliding-renewal"`
	SessionRenewBefore      string   `glazed:"auth-session-renew-before"`
	OIDCIssuerURL           string   `glazed:"oidc-issuer-url"`
	OIDCClientID            string   `glazed:"oidc-client-id"`
	OIDCClientSecret        string   `glazed:"oidc-client-secret"`
	OIDCRedirectURL         string   `glazed:"oidc-redirect-url"`
	OIDCScopes              []string `glazed:"oidc-scopes"`
	SessionTTLValue         time.Duration
	SessionRenewBeforeValue time.Duration
}

type UserInfo struct {
	Authenticated     bool     `json:"authenticated"`
	AuthMode          string   `json:"authMode"`
	Issuer            string   `json:"issuer,omitempty"`
	Subject           string   `json:"subject,omitempty"`
	Email             string   `json:"email,omitempty"`
	EmailVerified     bool     `json:"emailVerified,omitempty"`
	PreferredUsername string   `json:"preferredUsername,omitempty"`
	DisplayName       string   `json:"displayName,omitempty"`
	Picture           string   `json:"picture,omitempty"`
	Scopes            []string `json:"scopes,omitempty"`
	SessionExpiresAt  string   `json:"sessionExpiresAt,omitempty"`
}

func NewSection() (schema.Section, error) {
	return schema.NewSection(
		AuthSectionSlug,
		"Authentication Settings",
		schema.WithFields(
			fields.New(
				"auth-mode",
				fields.TypeChoice,
				fields.WithChoices(AuthModeDev, AuthModeOIDC),
				fields.WithHelp("Authentication mode for Draft Review"),
				fields.WithDefault(envOr("DRAFT_REVIEW_AUTH_MODE", AuthModeDev)),
			),
			fields.New(
				"auth-dev-user-id",
				fields.TypeString,
				fields.WithHelp("Development user ID returned by /api/me in auth-mode=dev"),
				fields.WithDefault(envOr("DRAFT_REVIEW_AUTH_DEV_USER_ID", "local-author")),
			),
			fields.New(
				"auth-session-cookie-name",
				fields.TypeString,
				fields.WithHelp("Cookie name used for browser sessions"),
				fields.WithDefault(envOr("DRAFT_REVIEW_AUTH_SESSION_COOKIE_NAME", DefaultSessionCookieName)),
			),
			fields.New(
				"auth-session-secret",
				fields.TypeString,
				fields.WithHelp("HMAC secret used to sign browser session cookies"),
				fields.WithDefault(strings.TrimSpace(os.Getenv("DRAFT_REVIEW_AUTH_SESSION_SECRET"))),
			),
			fields.New(
				"auth-session-ttl",
				fields.TypeString,
				fields.WithHelp("Application-managed browser session TTL used for opaque author sessions"),
				fields.WithDefault(envOr("DRAFT_REVIEW_AUTH_SESSION_TTL", "12h")),
			),
			fields.New(
				"auth-session-sliding-renewal",
				fields.TypeBool,
				fields.WithHelp("Whether authenticated browser sessions should renew before expiry while the user is active"),
				fields.WithDefault(envOrBool("DRAFT_REVIEW_AUTH_SESSION_SLIDING_RENEWAL", true)),
			),
			fields.New(
				"auth-session-renew-before",
				fields.TypeString,
				fields.WithHelp("Renew the session when the remaining lifetime is at or below this threshold"),
				fields.WithDefault(envOr("DRAFT_REVIEW_AUTH_SESSION_RENEW_BEFORE", "1h")),
			),
			fields.New(
				"oidc-issuer-url",
				fields.TypeString,
				fields.WithHelp("OIDC issuer URL exposed by Keycloak"),
				fields.WithDefault(strings.TrimSpace(os.Getenv("DRAFT_REVIEW_OIDC_ISSUER_URL"))),
			),
			fields.New(
				"oidc-client-id",
				fields.TypeString,
				fields.WithHelp("OIDC client ID for the Draft Review web app"),
				fields.WithDefault(strings.TrimSpace(os.Getenv("DRAFT_REVIEW_OIDC_CLIENT_ID"))),
			),
			fields.New(
				"oidc-client-secret",
				fields.TypeString,
				fields.WithHelp("OIDC client secret for confidential clients"),
				fields.WithDefault(strings.TrimSpace(os.Getenv("DRAFT_REVIEW_OIDC_CLIENT_SECRET"))),
			),
			fields.New(
				"oidc-redirect-url",
				fields.TypeString,
				fields.WithHelp("Redirect URL handled by this app after Keycloak login"),
				fields.WithDefault(envOr("DRAFT_REVIEW_OIDC_REDIRECT_URL", "http://127.0.0.1:8080/auth/callback")),
			),
			fields.New(
				"oidc-scopes",
				fields.TypeStringList,
				fields.WithHelp("Scopes requested during browser login"),
				fields.WithDefault([]string{"openid", "profile", "email"}),
			),
		),
	)
}

func LoadSettingsFromParsedValues(parsedValues *values.Values) (*Settings, error) {
	if parsedValues == nil {
		return nil, errors.New("parsed values are nil")
	}

	settings := &Settings{}
	if err := parsedValues.DecodeSectionInto(AuthSectionSlug, settings); err != nil {
		return nil, errors.Wrap(err, "failed to decode auth section")
	}

	settings.Mode = normalizeMode(settings.Mode)
	settings.DevUserID = strings.TrimSpace(settings.DevUserID)
	if settings.DevUserID == "" {
		settings.DevUserID = "local-author"
	}
	settings.SessionCookieName = strings.TrimSpace(settings.SessionCookieName)
	if settings.SessionCookieName == "" {
		settings.SessionCookieName = DefaultSessionCookieName
	}
	settings.SessionSecret = strings.TrimSpace(settings.SessionSecret)
	settings.SessionTTL = strings.TrimSpace(settings.SessionTTL)
	if settings.SessionTTL == "" {
		settings.SessionTTL = "12h"
	}
	sessionTTL, err := time.ParseDuration(settings.SessionTTL)
	if err != nil {
		return nil, errors.Wrap(err, "invalid auth-session-ttl")
	}
	if sessionTTL <= 0 {
		return nil, errors.New("auth-session-ttl must be greater than zero")
	}
	settings.SessionTTLValue = sessionTTL
	settings.SessionRenewBefore = strings.TrimSpace(settings.SessionRenewBefore)
	if settings.SessionRenewBefore == "" {
		settings.SessionRenewBefore = "1h"
	}
	sessionRenewBefore, err := time.ParseDuration(settings.SessionRenewBefore)
	if err != nil {
		return nil, errors.Wrap(err, "invalid auth-session-renew-before")
	}
	if sessionRenewBefore < 0 {
		return nil, errors.New("auth-session-renew-before must not be negative")
	}
	settings.SessionRenewBeforeValue = sessionRenewBefore
	settings.OIDCIssuerURL = strings.TrimSpace(settings.OIDCIssuerURL)
	settings.OIDCClientID = strings.TrimSpace(settings.OIDCClientID)
	settings.OIDCClientSecret = strings.TrimSpace(settings.OIDCClientSecret)
	settings.OIDCRedirectURL = strings.TrimSpace(settings.OIDCRedirectURL)
	settings.OIDCScopes = compact(settings.OIDCScopes)

	if settings.Mode == AuthModeOIDC {
		if settings.SessionSecret == "" {
			return nil, errors.New("auth-session-secret is required when auth-mode=oidc")
		}
		if settings.OIDCIssuerURL == "" {
			return nil, errors.New("oidc-issuer-url is required when auth-mode=oidc")
		}
		if settings.OIDCClientID == "" {
			return nil, errors.New("oidc-client-id is required when auth-mode=oidc")
		}
		if settings.OIDCRedirectURL == "" {
			return nil, errors.New("oidc-redirect-url is required when auth-mode=oidc")
		}
	}

	return settings, nil
}

func normalizeMode(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case AuthModeOIDC:
		return AuthModeOIDC
	case AuthModeDev:
		return AuthModeDev
	default:
		return AuthModeDev
	}
}

func envOr(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func envOrBool(key string, fallback bool) bool {
	value := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	case "":
		return fallback
	default:
		return fallback
	}
}

func compact(values []string) []string {
	ret := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		ret = append(ret, value)
	}
	return ret
}
