package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
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

type sessionEnvelope struct {
	Claims SessionClaims `json:"claims"`
}

type SessionManager struct {
	cookieName  string
	secret      []byte
	redirectURL string
}

func NewSessionManager(cookieName, secret, redirectURL string) (*SessionManager, error) {
	cookieName = strings.TrimSpace(cookieName)
	if cookieName == "" {
		cookieName = DefaultSessionCookieName
	}
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, errors.New("session secret is required")
	}

	return &SessionManager{
		cookieName:  cookieName,
		secret:      []byte(secret),
		redirectURL: strings.TrimSpace(redirectURL),
	}, nil
}

func (m *SessionManager) WriteSession(w http.ResponseWriter, r *http.Request, claims SessionClaims) error {
	if claims.IssuedAt.IsZero() {
		claims.IssuedAt = time.Now().UTC()
	}
	if claims.ExpiresAt.IsZero() {
		claims.ExpiresAt = claims.IssuedAt.Add(24 * time.Hour)
	}

	token, err := m.encode(sessionEnvelope{Claims: claims})
	if err != nil {
		return err
	}

	maxAge := int(time.Until(claims.ExpiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(w, &http.Cookie{
		Name:     m.cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   shouldUseSecureCookies(r, m.redirectURL),
		SameSite: http.SameSiteLaxMode,
		Expires:  claims.ExpiresAt.UTC(),
		MaxAge:   maxAge,
	})

	return nil
}

func (m *SessionManager) ReadSession(r *http.Request) (*SessionClaims, error) {
	if r == nil {
		return nil, ErrNoSession
	}

	cookie, err := r.Cookie(m.cookieName)
	if err != nil {
		return nil, ErrNoSession
	}
	if strings.TrimSpace(cookie.Value) == "" {
		return nil, ErrNoSession
	}

	envelope, err := m.decode(cookie.Value)
	if err != nil {
		return nil, err
	}
	if !envelope.Claims.ExpiresAt.IsZero() && time.Now().UTC().After(envelope.Claims.ExpiresAt.UTC()) {
		return nil, errors.New("session has expired")
	}

	return &envelope.Claims, nil
}

func (m *SessionManager) ClearSession(w http.ResponseWriter, r *http.Request) {
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

func (m *SessionManager) encode(envelope sessionEnvelope) (string, error) {
	payload, err := json.Marshal(envelope)
	if err != nil {
		return "", err
	}

	payloadEncoded := base64.RawURLEncoding.EncodeToString(payload)
	signature := m.sign(payloadEncoded)
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)
	return payloadEncoded + "." + signatureEncoded, nil
}

func (m *SessionManager) decode(raw string) (*sessionEnvelope, error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid session token format")
	}

	expectedSignature := m.sign(parts[0])
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid session signature encoding: %w", err)
	}
	if !hmac.Equal(signature, expectedSignature) {
		return nil, errors.New("invalid session signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid session payload encoding: %w", err)
	}

	var envelope sessionEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return nil, fmt.Errorf("invalid session payload: %w", err)
	}

	return &envelope, nil
}

func (m *SessionManager) sign(payload string) []byte {
	hash := hmac.New(sha256.New, m.secret)
	_, _ = hash.Write([]byte(payload))
	return hash.Sum(nil)
}
