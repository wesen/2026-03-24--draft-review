package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrNoSession = errors.New("no session present")

type SessionClaims struct {
	Issuer            string    `json:"iss"`
	Subject           string    `json:"sub"`
	Email             string    `json:"email,omitempty"`
	EmailVerified     bool      `json:"email_verified,omitempty"`
	PreferredUsername string    `json:"preferred_username,omitempty"`
	DisplayName       string    `json:"name,omitempty"`
	Picture           string    `json:"picture,omitempty"`
	Scopes            []string  `json:"scopes,omitempty"`
	IssuedAt          time.Time `json:"iat"`
	ExpiresAt         time.Time `json:"exp"`
}

type SessionStore interface {
	CreateAuthorSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (*Session, error)
	FindAuthorSessionByTokenHash(ctx context.Context, tokenHash string) (*ResolvedSession, error)
	TouchAuthorSession(ctx context.Context, sessionID uuid.UUID, touchedAt time.Time, renewedExpiresAt *time.Time) (*Session, error)
	RevokeAuthorSessionByTokenHash(ctx context.Context, tokenHash string) error
}

type ActiveSession struct {
	Resolved *ResolvedSession
	Claims   SessionClaims
	Renewed  bool
}

type SessionManager struct {
	cookieName  string
	secret      []byte
	redirectURL string
	sessionTTL  time.Duration
	store       SessionStore
	now         func() time.Time
	newToken    func() (string, error)
}

func NewSessionManager(cookieName, secret, redirectURL string, sessionTTL time.Duration, store SessionStore) (*SessionManager, error) {
	cookieName = strings.TrimSpace(cookieName)
	if cookieName == "" {
		cookieName = DefaultSessionCookieName
	}
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, errors.New("session secret is required")
	}
	if sessionTTL <= 0 {
		return nil, errors.New("session ttl is required")
	}
	if store == nil {
		return nil, errors.New("session store is required")
	}

	return &SessionManager{
		cookieName:  cookieName,
		secret:      []byte(secret),
		redirectURL: strings.TrimSpace(redirectURL),
		sessionTTL:  sessionTTL,
		store:       store,
		now:         func() time.Time { return time.Now().UTC() },
		newToken:    randomToken,
	}, nil
}

func (m *SessionManager) WriteSession(ctx context.Context, w http.ResponseWriter, r *http.Request, user *User) error {
	if user == nil {
		return errors.New("user is required")
	}

	userID, err := uuid.Parse(strings.TrimSpace(user.ID))
	if err != nil {
		return err
	}

	rawToken, err := m.newToken()
	if err != nil {
		return err
	}

	expiresAt := m.now().Add(m.sessionTTL)
	session, err := m.store.CreateAuthorSession(ctx, userID, m.tokenHash(rawToken), expiresAt)
	if err != nil {
		return err
	}

	maxAge := int(time.Until(session.ExpiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Value:    rawToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   shouldUseSecureCookies(r, m.redirectURL),
		SameSite: http.SameSiteLaxMode,
		Expires:  session.ExpiresAt.UTC(),
		MaxAge:   maxAge,
	})

	return nil
}

func (m *SessionManager) ReadSession(ctx context.Context, r *http.Request) (*SessionClaims, error) {
	active, err := m.ReadSessionState(ctx, r)
	if err != nil {
		return nil, err
	}
	return &active.Claims, nil
}

func (m *SessionManager) ReadSessionState(ctx context.Context, r *http.Request) (*ActiveSession, error) {
	if r == nil {
		return nil, ErrNoSession
	}

	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return nil, ErrNoSession
	}
	rawToken := strings.TrimSpace(cookie.Value)
	if rawToken == "" {
		return nil, ErrNoSession
	}

	resolved, err := m.store.FindAuthorSessionByTokenHash(ctx, m.tokenHash(rawToken))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNoSession
		}
		return nil, err
	}
	if resolved == nil {
		return nil, ErrNoSession
	}
	now := m.now()
	if resolved.Session.RevokedAt != nil {
		return nil, errors.New("session has been revoked")
	}
	if !resolved.Session.ExpiresAt.IsZero() && now.After(resolved.Session.ExpiresAt.UTC()) {
		return nil, errors.New("session has expired")
	}

	sessionID, err := uuid.Parse(strings.TrimSpace(resolved.Session.ID))
	if err != nil {
		return nil, err
	}
	touchedSession, err := m.store.TouchAuthorSession(ctx, sessionID, now, nil)
	if err != nil {
		return nil, err
	}
	resolved.Session = *touchedSession

	claims := SessionClaims{
		Issuer:            resolved.User.AuthIssuer,
		Subject:           resolved.User.AuthSubject,
		Email:             resolved.User.Email,
		EmailVerified:     resolved.User.EmailVerifiedAt != nil,
		PreferredUsername: derivePreferredUsername(resolved.User),
		DisplayName:       resolved.User.Name,
		IssuedAt:          resolved.Session.CreatedAt.UTC(),
		ExpiresAt:         resolved.Session.ExpiresAt.UTC(),
	}

	return &ActiveSession{
		Resolved: resolved,
		Claims:   claims,
		Renewed:  false,
	}, nil
}

func (m *SessionManager) ClearSession(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if r != nil {
		if cookie, err := r.Cookie(m.cookieName); err == nil {
			rawToken := strings.TrimSpace(cookie.Value)
			if rawToken != "" {
				_ = m.store.RevokeAuthorSessionByTokenHash(ctx, m.tokenHash(rawToken))
			}
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   shouldUseSecureCookies(r, m.redirectURL),
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
	})
}

func (c SessionClaims) UserInfo(authMode string) UserInfo {
	info := UserInfo{
		Authenticated:     true,
		AuthMode:          authMode,
		Issuer:            c.Issuer,
		Subject:           c.Subject,
		Email:             c.Email,
		EmailVerified:     c.EmailVerified,
		PreferredUsername: c.PreferredUsername,
		DisplayName:       c.DisplayName,
		Picture:           c.Picture,
		Scopes:            append([]string(nil), c.Scopes...),
	}
	if !c.ExpiresAt.IsZero() {
		info.SessionExpiresAt = c.ExpiresAt.UTC().Format(time.RFC3339)
	}
	return info
}

func shouldUseSecureCookies(r *http.Request, redirectURL string) bool {
	if r != nil {
		if r.TLS != nil {
			return true
		}
		if strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")), "https") {
			return true
		}
	}
	if strings.TrimSpace(redirectURL) == "" {
		return false
	}
	parsed, err := url.Parse(redirectURL)
	return err == nil && strings.EqualFold(parsed.Scheme, "https")
}

func (m *SessionManager) tokenHash(raw string) string {
	hash := hmac.New(sha256.New, m.secret)
	_, _ = hash.Write([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func derivePreferredUsername(user User) string {
	email := strings.TrimSpace(user.Email)
	if email == "" {
		return ""
	}
	if at := strings.Index(email, "@"); at > 0 {
		return email[:at]
	}
	return email
}
