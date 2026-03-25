package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	defaultOIDCHTTPTimeout = 10 * time.Second
	defaultJWKSRefresh     = 5 * time.Minute
	authStateCookieName    = "draft_review_auth_state"
	authNonceCookieName    = "draft_review_auth_nonce"
)

type WebHandler interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleCallback(w http.ResponseWriter, r *http.Request)
	HandleLogout(w http.ResponseWriter, r *http.Request)
	HandleLogoutCallback(w http.ResponseWriter, r *http.Request)
}

type oidcDiscoveryDocument struct {
	Issuer                string `json:"issuer"`
	JWKSURI               string `json:"jwks_uri"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	EndSessionEndpoint    string `json:"end_session_endpoint"`
}

type idTokenClaims struct {
	jwt.Claims
	Nonce             string `json:"nonce,omitempty"`
	Email             string `json:"email,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Name              string `json:"name,omitempty"`
	Picture           string `json:"picture,omitempty"`
}

type oauthStatePayload struct {
	ID       string `json:"id"`
	ReturnTo string `json:"rt,omitempty"`
}

type OIDCAuthenticator struct {
	settings      *Settings
	sessions      *SessionManager
	httpClient    *http.Client
	discovery     oidcDiscoveryDocument
	oauthConfig   oauth2.Config
	jwks          *jwksCache
	now           func() time.Time
	newRandom     func() (string, error)
	postLoginPath string
}

func NewOIDCAuthenticator(ctx context.Context, settings *Settings, sessions *SessionManager) (*OIDCAuthenticator, error) {
	return newOIDCAuthenticatorWithClient(ctx, settings, sessions, &http.Client{Timeout: defaultOIDCHTTPTimeout})
}

func newOIDCAuthenticatorWithClient(
	ctx context.Context,
	settings *Settings,
	sessions *SessionManager,
	httpClient *http.Client,
) (*OIDCAuthenticator, error) {
	if settings == nil {
		return nil, errors.New("oidc settings are nil")
	}
	if sessions == nil {
		return nil, errors.New("session manager is required")
	}
	if strings.TrimSpace(settings.OIDCIssuerURL) == "" {
		return nil, errors.New("oidc issuer url is required")
	}
	if strings.TrimSpace(settings.OIDCClientID) == "" {
		return nil, errors.New("oidc client id is required")
	}
	if strings.TrimSpace(settings.OIDCRedirectURL) == "" {
		return nil, errors.New("oidc redirect url is required")
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultOIDCHTTPTimeout}
	}

	discoveryURL := strings.TrimRight(settings.OIDCIssuerURL, "/") + "/.well-known/openid-configuration"
	discovery, err := fetchOIDCDiscovery(ctx, httpClient, discoveryURL)
	if err != nil {
		return nil, err
	}
	if discovery.Issuer == "" || discovery.AuthorizationEndpoint == "" || discovery.TokenEndpoint == "" || discovery.JWKSURI == "" {
		return nil, errors.New("oidc discovery document is incomplete")
	}

	scopes := append([]string(nil), settings.OIDCScopes...)
	if len(scopes) == 0 {
		scopes = []string{"openid", "profile", "email"}
	}

	return &OIDCAuthenticator{
		settings:   settings,
		sessions:   sessions,
		httpClient: httpClient,
		discovery:  discovery,
		oauthConfig: oauth2.Config{
			ClientID:     settings.OIDCClientID,
			ClientSecret: settings.OIDCClientSecret,
			RedirectURL:  settings.OIDCRedirectURL,
			Scopes:       scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  discovery.AuthorizationEndpoint,
				TokenURL: discovery.TokenEndpoint,
			},
		},
		jwks: &jwksCache{
			client:          httpClient,
			jwksURI:         discovery.JWKSURI,
			refreshInterval: defaultJWKSRefresh,
		},
		now:           func() time.Time { return time.Now().UTC() },
		newRandom:     randomToken,
		postLoginPath: "/",
	}, nil
}

func (a *OIDCAuthenticator) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state, err := a.newRandom()
	if err != nil {
		http.Error(w, "failed to create auth state", http.StatusInternalServerError)
		return
	}
	nonce, err := a.newRandom()
	if err != nil {
		http.Error(w, "failed to create auth nonce", http.StatusInternalServerError)
		return
	}

	secureCookies := shouldUseSecureCookies(r, a.oauthConfig.RedirectURL)
	returnTo, err := a.resolveRequestedRedirect(r, r.URL.Query().Get("return_to"))
	if err != nil {
		http.Error(w, "invalid return_to parameter", http.StatusBadRequest)
		return
	}
	stateValue, err := marshalOAuthState(oauthStatePayload{ID: state, ReturnTo: returnTo})
	if err != nil {
		http.Error(w, "failed to encode auth state", http.StatusInternalServerError)
		return
	}

	setShortLivedCookie(w, authStateCookieName, state, secureCookies)
	setShortLivedCookie(w, authNonceCookieName, nonce, secureCookies)

	loginURL := a.oauthConfig.AuthCodeURL(stateValue, oauth2.SetAuthURLParam("nonce", nonce))
	http.Redirect(w, r, loginURL, http.StatusFound)
}

func (a *OIDCAuthenticator) HandleCallback(w http.ResponseWriter, r *http.Request) {
	queryState := strings.TrimSpace(r.URL.Query().Get("state"))
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	if queryState == "" || code == "" {
		http.Error(w, "missing oauth callback parameters", http.StatusBadRequest)
		return
	}

	statePayload, err := unmarshalOAuthState(queryState)
	if err != nil {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}

	stateCookie, err := r.Cookie(authStateCookieName)
	if err != nil || strings.TrimSpace(stateCookie.Value) == "" || stateCookie.Value != statePayload.ID {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}
	nonceCookie, err := r.Cookie(authNonceCookieName)
	if err != nil || strings.TrimSpace(nonceCookie.Value) == "" {
		http.Error(w, "missing oauth nonce", http.StatusBadRequest)
		return
	}

	secureCookies := shouldUseSecureCookies(r, a.oauthConfig.RedirectURL)
	clearCookie(w, authStateCookieName, secureCookies)
	clearCookie(w, authNonceCookieName, secureCookies)

	ctx := context.WithValue(r.Context(), oauth2.HTTPClient, a.httpClient)
	token, err := a.oauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, fmt.Sprintf("token exchange failed: %v", err), http.StatusBadGateway)
		return
	}

	rawIDToken, _ := token.Extra("id_token").(string)
	if strings.TrimSpace(rawIDToken) == "" {
		http.Error(w, "missing id_token in token response", http.StatusBadGateway)
		return
	}

	claims, err := a.verifyIDToken(ctx, rawIDToken, nonceCookie.Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("id_token verification failed: %v", err), http.StatusBadGateway)
		return
	}

	expiry := token.Expiry
	if expiry.IsZero() {
		expiry = a.now().Add(24 * time.Hour)
	}

	sessionClaims := SessionClaims{
		Issuer:            claims.Issuer,
		Subject:           claims.Subject,
		Email:             claims.Email,
		EmailVerified:     claims.EmailVerified,
		PreferredUsername: claims.PreferredUsername,
		DisplayName:       claims.Name,
		Picture:           claims.Picture,
		Scopes:            append([]string(nil), a.oauthConfig.Scopes...),
		IssuedAt:          a.now(),
		ExpiresAt:         expiry.UTC(),
	}

	if err := a.sessions.WriteSession(w, r, sessionClaims); err != nil {
		http.Error(w, fmt.Sprintf("session creation failed: %v", err), http.StatusInternalServerError)
		return
	}

	redirectTarget := a.postLoginPath
	if strings.TrimSpace(statePayload.ReturnTo) != "" {
		redirectTarget = statePayload.ReturnTo
	}
	http.Redirect(w, r, redirectTarget, http.StatusSeeOther)
}

func (a *OIDCAuthenticator) HandleLogout(w http.ResponseWriter, r *http.Request) {
	a.sessions.ClearSession(w, r)
	returnTo, err := a.resolveRequestedRedirect(r, r.URL.Query().Get("return_to"))
	if err != nil {
		http.Error(w, "invalid return_to parameter", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, a.buildLogoutRedirectURL(returnTo), http.StatusSeeOther)
}

func (a *OIDCAuthenticator) HandleLogoutCallback(w http.ResponseWriter, r *http.Request) {
	returnTo, err := a.resolveRequestedRedirect(r, r.URL.Query().Get("return_to"))
	if err != nil {
		http.Error(w, "invalid return_to parameter", http.StatusBadRequest)
		return
	}

	redirectTarget := a.postLoginPath
	if strings.TrimSpace(returnTo) != "" {
		redirectTarget = returnTo
	}
	http.Redirect(w, r, redirectTarget, http.StatusSeeOther)
}

func (a *OIDCAuthenticator) buildLogoutRedirectURL(returnTo string) string {
	fallback := a.postLoginPath
	if strings.TrimSpace(returnTo) != "" {
		fallback = returnTo
	}

	if strings.TrimSpace(a.discovery.EndSessionEndpoint) == "" {
		return fallback
	}

	endSessionURL, err := url.Parse(a.discovery.EndSessionEndpoint)
	if err != nil {
		return fallback
	}

	postLogoutURL, err := buildLogoutCallbackURL(a.oauthConfig.RedirectURL, returnTo)
	if err != nil {
		return fallback
	}

	query := endSessionURL.Query()
	query.Set("post_logout_redirect_uri", postLogoutURL)
	if strings.TrimSpace(a.settings.OIDCClientID) != "" {
		query.Set("client_id", a.settings.OIDCClientID)
	}
	endSessionURL.RawQuery = query.Encode()
	return endSessionURL.String()
}

func buildLogoutCallbackURL(redirectURL, returnTo string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(redirectURL))
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", errors.New("redirect url must include scheme and host")
	}
	parsed.Path = "/auth/logout/callback"
	parsed.RawPath = ""
	parsed.Fragment = ""

	query := parsed.Query()
	if strings.TrimSpace(returnTo) != "" {
		query.Set("return_to", returnTo)
	} else {
		query.Del("return_to")
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

func (a *OIDCAuthenticator) resolveRequestedRedirect(r *http.Request, rawValue string) (string, error) {
	value := strings.TrimSpace(rawValue)
	if value == "" {
		return "", nil
	}

	if strings.HasPrefix(value, "/") {
		return value, nil
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errors.New("return_to must use http or https")
	}
	if parsed.Host == "" {
		return "", errors.New("return_to must include a host")
	}

	redirectURL, err := url.Parse(strings.TrimSpace(a.oauthConfig.RedirectURL))
	if err != nil {
		return "", err
	}

	allowedHosts := []string{
		hostWithoutPort(r.Host),
		hostWithoutPort(redirectURL.Host),
		redirectURL.Hostname(),
	}
	for _, host := range allowedHosts {
		if hostnamesMatch(parsed.Hostname(), host) {
			return parsed.String(), nil
		}
	}

	return "", errors.New("return_to host is not allowed")
}

func hostWithoutPort(value string) string {
	host := strings.TrimSpace(value)
	if host == "" {
		return ""
	}
	if parsedHost, _, err := net.SplitHostPort(host); err == nil {
		return parsedHost
	}
	return host
}

func hostnamesMatch(a, b string) bool {
	left := strings.ToLower(strings.TrimSpace(a))
	right := strings.ToLower(strings.TrimSpace(b))
	if left == "" || right == "" {
		return false
	}
	if left == right {
		return true
	}
	return (left == "localhost" && right == "127.0.0.1") || (left == "127.0.0.1" && right == "localhost")
}

func (a *OIDCAuthenticator) verifyIDToken(ctx context.Context, rawIDToken, expectedNonce string) (*idTokenClaims, error) {
	parsed, err := jwt.ParseSigned(rawIDToken)
	if err != nil {
		return nil, err
	}

	kid := ""
	if len(parsed.Headers) > 0 {
		kid = parsed.Headers[0].KeyID
	}

	keys, err := a.jwks.Keys(ctx, kid, false)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for _, key := range keys {
		var claims idTokenClaims
		if err := parsed.Claims(key.Key, &claims); err != nil {
			lastErr = err
			continue
		}

		expected := jwt.Expected{
			Issuer:   a.discovery.Issuer,
			Audience: jwt.Audience{a.settings.OIDCClientID},
			Time:     a.now(),
		}
		if err := claims.Validate(expected); err != nil {
			lastErr = err
			continue
		}
		if strings.TrimSpace(expectedNonce) != "" && claims.Nonce != expectedNonce {
			lastErr = errors.New("unexpected nonce")
			continue
		}

		return &claims, nil
	}

	if lastErr == nil {
		lastErr = errors.New("no valid jwks keys")
	}
	return nil, lastErr
}

type jwksCache struct {
	client          *http.Client
	jwksURI         string
	refreshInterval time.Duration
	lastFetched     time.Time
	keySet          jose.JSONWebKeySet
}

func (c *jwksCache) Keys(ctx context.Context, kid string, forceRefresh bool) ([]jose.JSONWebKey, error) {
	if forceRefresh || len(c.keySet.Keys) == 0 || time.Since(c.lastFetched) >= c.refreshInterval {
		if err := c.refresh(ctx); err != nil {
			return nil, err
		}
	}

	if kid != "" {
		if keys := c.keySet.Key(kid); len(keys) > 0 {
			return keys, nil
		}
	}
	if len(c.keySet.Keys) == 0 {
		return nil, errors.New("jwks key set is empty")
	}

	return append([]jose.JSONWebKey(nil), c.keySet.Keys...), nil
}

func (c *jwksCache) refresh(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.jwksURI, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jwks request failed with status %d", resp.StatusCode)
	}

	var keySet jose.JSONWebKeySet
	if err := json.NewDecoder(resp.Body).Decode(&keySet); err != nil {
		return err
	}
	c.keySet = keySet
	c.lastFetched = time.Now()
	return nil
}

func fetchOIDCDiscovery(ctx context.Context, client *http.Client, discoveryURL string) (oidcDiscoveryDocument, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, nil)
	if err != nil {
		return oidcDiscoveryDocument{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return oidcDiscoveryDocument{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return oidcDiscoveryDocument{}, fmt.Errorf("oidc discovery request failed with status %d", resp.StatusCode)
	}

	var doc oidcDiscoveryDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return oidcDiscoveryDocument{}, err
	}
	return doc, nil
}

func setShortLivedCookie(w http.ResponseWriter, name, value string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((10 * time.Minute).Seconds()),
	})
}

func clearCookie(w http.ResponseWriter, name string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0).UTC(),
	})
}

func marshalOAuthState(payload oauthStatePayload) (string, error) {
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func unmarshalOAuthState(value string) (oauthStatePayload, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return oauthStatePayload{}, err
	}

	var payload oauthStatePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return oauthStatePayload{}, err
	}
	if strings.TrimSpace(payload.ID) == "" {
		return oauthStatePayload{}, errors.New("oauth state id is required")
	}

	return payload, nil
}

func randomToken() (string, error) {
	raw := make([]byte, 24)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
