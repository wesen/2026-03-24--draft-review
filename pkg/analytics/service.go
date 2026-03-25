package analytics

import (
	"context"
	"strings"

	"github.com/go-go-golems/draft-review/pkg/reviewlinks"
	"github.com/go-go-golems/draft-review/pkg/reviews"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListReaders(ctx context.Context, ownerUserID, articleID string) ([]reviewlinks.Reader, error) {
	if s == nil || s.repo == nil {
		return []reviewlinks.Reader{}, nil
	}
	if strings.TrimSpace(ownerUserID) == "" || strings.TrimSpace(articleID) == "" {
		return nil, ErrNotFound
	}
	return s.repo.ListReaders(ctx, ownerUserID, articleID)
}

func (s *Service) ListReactions(ctx context.Context, ownerUserID, articleID string) ([]reviews.Reaction, error) {
	if s == nil || s.repo == nil {
		return []reviews.Reaction{}, nil
	}
	if strings.TrimSpace(ownerUserID) == "" || strings.TrimSpace(articleID) == "" {
		return nil, ErrNotFound
	}
	return s.repo.ListReactions(ctx, ownerUserID, articleID)
}

func (s *Service) GetArticleAnalytics(ctx context.Context, ownerUserID, articleID string) (*ArticleAnalytics, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}
	if strings.TrimSpace(ownerUserID) == "" || strings.TrimSpace(articleID) == "" {
		return nil, ErrNotFound
	}
	return s.repo.GetArticleAnalytics(ctx, ownerUserID, articleID)
}

func (s *Service) ListReaderDirectory(ctx context.Context, ownerUserID string) ([]ReaderContact, error) {
	if s == nil || s.repo == nil {
		return []ReaderContact{}, nil
	}
	if strings.TrimSpace(ownerUserID) == "" {
		return nil, ErrNotFound
	}
	return s.repo.ListReaderDirectory(ctx, ownerUserID)
}
