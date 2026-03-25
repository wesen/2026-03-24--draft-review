package analytics

import (
	"context"
	"time"

	"github.com/go-go-golems/draft-review/pkg/reviewlinks"
	"github.com/go-go-golems/draft-review/pkg/reviews"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("analytics subject not found")

type SectionAnalytics struct {
	SectionID        string         `json:"sectionId"`
	Title            string         `json:"title"`
	Position         int            `json:"position"`
	ReadersReached   int            `json:"readersReached"`
	ReadersCompleted int            `json:"readersCompleted"`
	ReactionCount    int            `json:"reactionCount"`
	ReactionCounts   map[string]int `json:"reactionCounts"`
}

type DraftKillerInsight struct {
	SectionID    string `json:"sectionId"`
	Title        string `json:"title"`
	DropOffCount int    `json:"dropOffCount"`
	DropOffRate  int    `json:"dropOffRate"`
	Summary      string `json:"summary"`
}

type ArticleAnalytics struct {
	ArticleID          string              `json:"articleId"`
	Version            string              `json:"version"`
	ReaderCount        int                 `json:"readerCount"`
	EngagedReaderCount int                 `json:"engagedReaderCount"`
	AverageProgress    int                 `json:"averageProgress"`
	TotalReactions     int                 `json:"totalReactions"`
	Sections           []SectionAnalytics  `json:"sections"`
	DraftKiller        *DraftKillerInsight `json:"draftKiller,omitempty"`
}

type ReaderArticleActivity struct {
	ArticleID     string     `json:"articleId"`
	Title         string     `json:"title"`
	Progress      int        `json:"progress"`
	ReactionCount int        `json:"reactionCount"`
	LastActiveAt  *time.Time `json:"lastActiveAt,omitempty"`
}

type ReaderContact struct {
	Email          string                  `json:"email"`
	Name           string                  `json:"name"`
	Avatar         string                  `json:"avatar"`
	TotalArticles  int                     `json:"totalArticles"`
	TotalReactions int                     `json:"totalReactions"`
	LastActiveAt   *time.Time              `json:"lastActiveAt,omitempty"`
	Status         string                  `json:"status"`
	Articles       []ReaderArticleActivity `json:"articles"`
}

type Repository interface {
	ListReaders(ctx context.Context, ownerUserID, articleID string) ([]reviewlinks.Reader, error)
	ListReactions(ctx context.Context, ownerUserID, articleID string) ([]reviews.Reaction, error)
	GetArticleAnalytics(ctx context.Context, ownerUserID, articleID string) (*ArticleAnalytics, error)
	ListReaderDirectory(ctx context.Context, ownerUserID string) ([]ReaderContact, error)
}
