package articles

import (
	"context"
	"strings"

	draftauth "github.com/go-go-golems/draft-review/pkg/auth"
)

type Repository interface {
	ListArticles(ctx context.Context, ownerUserID string) ([]Article, error)
	GetArticle(ctx context.Context, ownerUserID, id string) (*Article, error)
	CreateArticle(ctx context.Context, ownerUserID string, input CreateArticleInput) (*Article, error)
	UpdateArticle(ctx context.Context, ownerUserID, id string, input UpdateArticleInput) (*Article, error)
	CreateVersion(ctx context.Context, ownerUserID, id string, input CreateVersionInput) (*Article, error)
	DeleteArticle(ctx context.Context, ownerUserID, id string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListArticles(ctx context.Context, owner *draftauth.User) ([]Article, error) {
	if s == nil || s.repo == nil {
		return []Article{}, nil
	}
	if owner == nil || strings.TrimSpace(owner.ID) == "" {
		return nil, ErrNotFound
	}
	return s.repo.ListArticles(ctx, owner.ID)
}

func (s *Service) GetArticle(ctx context.Context, owner *draftauth.User, id string) (*Article, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}
	if owner == nil || strings.TrimSpace(owner.ID) == "" {
		return nil, ErrNotFound
	}
	return s.repo.GetArticle(ctx, owner.ID, id)
}

func (s *Service) CreateArticle(ctx context.Context, owner *draftauth.User, input CreateArticleInput) (*Article, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}
	if owner == nil || strings.TrimSpace(owner.ID) == "" {
		return nil, ErrNotFound
	}

	input.Title = defaultIfBlank(input.Title, "Untitled Article")
	input.Author = defaultIfBlank(input.Author, "You")
	input.Intro = strings.TrimSpace(input.Intro)

	return s.repo.CreateArticle(ctx, owner.ID, input)
}

func (s *Service) UpdateArticle(ctx context.Context, owner *draftauth.User, id string, input UpdateArticleInput) (*Article, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}
	if owner == nil || strings.TrimSpace(owner.ID) == "" {
		return nil, ErrNotFound
	}

	if input.Title != nil {
		normalized := strings.TrimSpace(*input.Title)
		if normalized == "" {
			return nil, NewValidationError("title cannot be empty")
		}
		input.Title = &normalized
	}

	if input.Author != nil {
		normalized := strings.TrimSpace(*input.Author)
		if normalized == "" {
			return nil, NewValidationError("author cannot be empty")
		}
		input.Author = &normalized
	}

	if input.Intro != nil {
		normalized := strings.TrimSpace(*input.Intro)
		input.Intro = &normalized
	}

	if input.Status != nil {
		normalized := strings.TrimSpace(*input.Status)
		if !isValidStatus(normalized) {
			return nil, NewValidationError("status must be one of draft, in_review, complete, archived")
		}
		input.Status = &normalized
	}

	if input.Sections != nil {
		sections := *input.Sections
		if len(sections) == 0 {
			return nil, NewValidationError("sections must contain at least one section")
		}

		normalized := make([]SectionInput, 0, len(sections))
		for _, section := range sections {
			title := strings.TrimSpace(section.Title)
			if title == "" {
				return nil, NewValidationError("section title cannot be empty")
			}

			paragraphs := normalizeParagraphs(section.Paragraphs)
			if len(paragraphs) == 0 {
				return nil, NewValidationError("section paragraphs must contain at least one paragraph")
			}

			normalized = append(normalized, SectionInput{
				ID:         strings.TrimSpace(section.ID),
				Title:      title,
				Paragraphs: paragraphs,
			})
		}

		input.Sections = &normalized
	}

	return s.repo.UpdateArticle(ctx, owner.ID, id, input)
}

func (s *Service) CreateVersion(ctx context.Context, owner *draftauth.User, id string, input CreateVersionInput) (*Article, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}
	if owner == nil || strings.TrimSpace(owner.ID) == "" {
		return nil, ErrNotFound
	}

	input.Label = strings.TrimSpace(input.Label)

	if input.Intro != nil {
		normalized := strings.TrimSpace(*input.Intro)
		input.Intro = &normalized
	}

	if input.AuthorNote != nil {
		normalized := strings.TrimSpace(*input.AuthorNote)
		input.AuthorNote = &normalized
	}

	return s.repo.CreateVersion(ctx, owner.ID, id, input)
}

func (s *Service) DeleteArticle(ctx context.Context, owner *draftauth.User, id string) error {
	if s == nil || s.repo == nil {
		return ErrNotFound
	}
	if owner == nil || strings.TrimSpace(owner.ID) == "" {
		return ErrNotFound
	}

	id = strings.TrimSpace(id)
	if id == "" {
		return ErrNotFound
	}

	return s.repo.DeleteArticle(ctx, owner.ID, id)
}

func defaultIfBlank(value, fallback string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fallback
	}
	return trimmed
}

func normalizeParagraphs(paragraphs []string) []string {
	ret := make([]string, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		trimmed := strings.TrimSpace(paragraph)
		if trimmed == "" {
			continue
		}
		ret = append(ret, trimmed)
	}
	if len(ret) == 0 {
		ret = append(ret, "")
	}
	return ret
}

func isValidStatus(status string) bool {
	switch status {
	case "draft", "in_review", "complete", "archived":
		return true
	default:
		return false
	}
}
