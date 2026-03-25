package reviews

import (
	"context"
	"testing"

	"github.com/go-go-golems/draft-review/pkg/reviewlinks"
)

type fakeRepository struct{}

func (f *fakeRepository) StartSession(ctx context.Context, link *reviewlinks.ResolvedLink, input StartSessionInput) (*StartSessionResult, error) {
	return &StartSessionResult{}, nil
}

func (f *fakeRepository) RecordProgress(ctx context.Context, sessionID string, input ProgressInput) (*ProgressState, error) {
	return &ProgressState{}, nil
}

func (f *fakeRepository) AddReaction(ctx context.Context, sessionID string, input ReactionInput) (*Reaction, error) {
	return &Reaction{}, nil
}

func (f *fakeRepository) SubmitSummary(ctx context.Context, sessionID string, input SummaryInput) (*Summary, error) {
	return &Summary{}, nil
}

func TestStartSessionRejectsAnonymousWhenArticleDisallowsIt(t *testing.T) {
	t.Parallel()

	service := NewService(&fakeRepository{})
	_, err := service.StartSession(context.Background(), &reviewlinks.ResolvedLink{
		ArticleID:      "article-1",
		AllowAnonymous: false,
	}, StartSessionInput{Anonymous: true})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}

func TestRecordProgressRejectsOutOfRangePercent(t *testing.T) {
	t.Parallel()

	service := NewService(&fakeRepository{})
	progress := 101
	_, err := service.RecordProgress(context.Background(), "session-1", ProgressInput{
		SectionID:       "section-1",
		ProgressPercent: &progress,
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}

func TestSubmitSummaryRejectsUnknownRecommendability(t *testing.T) {
	t.Parallel()

	service := NewService(&fakeRepository{})
	_, err := service.SubmitSummary(context.Background(), "session-1", SummaryInput{
		Recommendability: "nope",
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}
