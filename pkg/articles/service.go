package articles

import "context"

type Repository interface {
	ListArticles(ctx context.Context) ([]Article, error)
	GetArticle(ctx context.Context, id string) (*Article, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListArticles(ctx context.Context) ([]Article, error) {
	if s == nil || s.repo == nil {
		return []Article{}, nil
	}
	return s.repo.ListArticles(ctx)
}

func (s *Service) GetArticle(ctx context.Context, id string) (*Article, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}
	return s.repo.GetArticle(ctx, id)
}
