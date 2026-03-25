package reviews

import (
	"time"

	"github.com/go-go-golems/draft-review/pkg/reviewlinks"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("review session not found")

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

type Session struct {
	ID               string     `json:"id"`
	ArticleID        string     `json:"articleId"`
	ArticleVersionID string     `json:"articleVersionId"`
	ReaderID         string     `json:"readerId"`
	ReaderName       string     `json:"readerName"`
	ReaderEmail      string     `json:"readerEmail,omitempty"`
	IsAnonymous      bool       `json:"isAnonymous"`
	ProgressPercent  int        `json:"progressPercent"`
	StartedAt        time.Time  `json:"startedAt"`
	LastActiveAt     time.Time  `json:"lastActiveAt"`
	CompletedAt      *time.Time `json:"completedAt,omitempty"`
}

type StartSessionInput struct {
	ReaderName string `json:"readerName,omitempty"`
	Anonymous  bool   `json:"anonymous,omitempty"`
}

type StartSessionResult struct {
	Session   *Session                   `json:"session"`
	Reader    reviewlinks.ReaderIdentity `json:"reader"`
	Article   reviewlinks.ReaderArticle  `json:"article"`
	Reactions []Reaction                 `json:"reactions,omitempty"`
}

type ProgressInput struct {
	SectionID       string `json:"sectionId"`
	ParagraphID     string `json:"paragraphId,omitempty"`
	Completed       bool   `json:"completed,omitempty"`
	ProgressPercent *int   `json:"progressPercent,omitempty"`
}

type ProgressState struct {
	SessionID       string    `json:"sessionId"`
	ProgressPercent int       `json:"progressPercent"`
	LastActiveAt    time.Time `json:"lastActiveAt"`
}

type ReactionInput struct {
	SectionID   string `json:"sectionId"`
	ParagraphID string `json:"paragraphId"`
	Type        string `json:"type"`
	Text        string `json:"text"`
}

type Reaction struct {
	ID          string    `json:"id"`
	ArticleID   string    `json:"articleId"`
	SectionID   string    `json:"sectionId"`
	ParagraphID string    `json:"paragraphId"`
	ReaderID    string    `json:"readerId"`
	ReaderName  string    `json:"readerName"`
	Type        string    `json:"type"`
	Text        string    `json:"text"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SummaryInput struct {
	OverallThoughts  string `json:"overallThoughts,omitempty"`
	Recommendability string `json:"recommendability,omitempty"`
	NotifyNewVersion bool   `json:"notifyNewVersion"`
}

type Summary struct {
	SessionID        string    `json:"sessionId"`
	OverallThoughts  string    `json:"overallThoughts,omitempty"`
	Recommendability string    `json:"recommendability,omitempty"`
	NotifyNewVersion bool      `json:"notifyNewVersion"`
	SubmittedAt      time.Time `json:"submittedAt"`
}
