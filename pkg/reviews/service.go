package reviews

import (
	"context"
	"strings"

	"github.com/go-go-golems/draft-review/pkg/reviewlinks"
)

type Repository interface {
	StartSession(ctx context.Context, link *reviewlinks.ResolvedLink, input StartSessionInput) (*StartSessionResult, error)
	RecordProgress(ctx context.Context, sessionID string, input ProgressInput) (*ProgressState, error)
	AddReaction(ctx context.Context, sessionID string, input ReactionInput) (*Reaction, error)
	SubmitSummary(ctx context.Context, sessionID string, input SummaryInput) (*Summary, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) StartSession(ctx context.Context, link *reviewlinks.ResolvedLink, input StartSessionInput) (*StartSessionResult, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}
	if link == nil || strings.TrimSpace(link.ArticleID) == "" {
		return nil, reviewlinks.ErrNotFound
	}

	input.ReaderName = strings.TrimSpace(input.ReaderName)
	if input.Anonymous && !link.AllowAnonymous {
		return nil, NewValidationError("anonymous reviews are not allowed for this article")
	}

	return s.repo.StartSession(ctx, link, input)
}

func (s *Service) RecordProgress(ctx context.Context, sessionID string, input ProgressInput) (*ProgressState, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}

	sessionID = strings.TrimSpace(sessionID)
	input.SectionID = strings.TrimSpace(input.SectionID)
	input.ParagraphID = strings.TrimSpace(input.ParagraphID)
	if sessionID == "" || input.SectionID == "" {
		return nil, ErrNotFound
	}

	if input.ProgressPercent != nil {
		normalized := *input.ProgressPercent
		if normalized < 0 || normalized > 100 {
			return nil, NewValidationError("progressPercent must be between 0 and 100")
		}
		input.ProgressPercent = &normalized
	}

	return s.repo.RecordProgress(ctx, sessionID, input)
}

func (s *Service) AddReaction(ctx context.Context, sessionID string, input ReactionInput) (*Reaction, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}

	sessionID = strings.TrimSpace(sessionID)
	input.SectionID = strings.TrimSpace(input.SectionID)
	input.ParagraphID = strings.TrimSpace(input.ParagraphID)
	input.Type = strings.TrimSpace(input.Type)
	input.Text = strings.TrimSpace(input.Text)

	if sessionID == "" || input.SectionID == "" || input.ParagraphID == "" {
		return nil, ErrNotFound
	}
	if input.Type == "" {
		return nil, NewValidationError("reaction type is required")
	}

	return s.repo.AddReaction(ctx, sessionID, input)
}

func (s *Service) SubmitSummary(ctx context.Context, sessionID string, input SummaryInput) (*Summary, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}

	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil, ErrNotFound
	}

	input.OverallThoughts = strings.TrimSpace(input.OverallThoughts)
	input.Recommendability = strings.TrimSpace(input.Recommendability)
	if input.Recommendability != "" && input.Recommendability != "maybe" && input.Recommendability != "yes" && input.Recommendability != "absolutely" {
		return nil, NewValidationError("recommendability must be one of maybe, yes, absolutely")
	}

	return s.repo.SubmitSummary(ctx, sessionID, input)
}
