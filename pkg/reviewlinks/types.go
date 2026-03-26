package reviewlinks

import (
	"time"

	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("review link not found")

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

type ShareLink struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

const (
	IdentityModeEmail     = "email"
	IdentityModeNamed     = "named"
	IdentityModeAnonymous = "anonymous"
	IdentityModePreview   = "preview"
)

type InviteInput struct {
	IdentityMode string `json:"identityMode,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`
	Email        string `json:"email,omitempty"`
	Note         string `json:"note"`
	IsPreview    bool   `json:"isPreview,omitempty"`
}

type Reader struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email,omitempty"`
	Avatar       string     `json:"avatar"`
	ArticleID    string     `json:"articleId"`
	Progress     int        `json:"progress"`
	Token        string     `json:"token"`
	IdentityMode string     `json:"identityMode,omitempty"`
	IsPreview    bool       `json:"isPreview,omitempty"`
	InvitedAt    time.Time  `json:"invitedAt"`
	LastActiveAt *time.Time `json:"lastActiveAt,omitempty"`
}

type ReaderIdentity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Section struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
}

type ReaderArticle struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	Version  string    `json:"version"`
	Intro    string    `json:"intro"`
	Sections []Section `json:"sections"`
}

type ResolvedLink struct {
	Token                 string         `json:"-"`
	ArticleID             string         `json:"-"`
	ArticleVersionID      string         `json:"-"`
	InviteID              string         `json:"-"`
	AccessMode            string         `json:"-"`
	ReaderEmail           string         `json:"-"`
	InviteDisplayName     string         `json:"-"`
	IdentityMode          string         `json:"-"`
	IsPreview             bool           `json:"-"`
	AllowAnonymous        bool           `json:"-"`
	RequireNote           bool           `json:"-"`
	ReaderCanSeeReactions bool           `json:"-"`
	ReaderCanSeeNames     bool           `json:"-"`
	ShowAuthorNote        bool           `json:"-"`
	Reader                ReaderIdentity `json:"reader"`
	Article               ReaderArticle  `json:"article"`
}
