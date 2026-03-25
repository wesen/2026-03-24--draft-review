package articles

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("article not found")

type PostgresRepository struct {
	pool *pgxpool.Pool
}

type articleRow struct {
	ID        uuid.UUID
	VersionID uuid.UUID
	Title     string
	Author    string
	Version   string
	Status    string
	Intro     string
	CreatedAt any
	UpdatedAt any
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) ListArticles(ctx context.Context) ([]Article, error) {
	if r == nil || r.pool == nil {
		return []Article{}, nil
	}

	rows, err := r.pool.Query(ctx, `
select
    a.id,
    a.current_version_id,
    a.title,
    a.author_display_name,
    v.version_label,
    a.status,
    v.intro,
    a.created_at,
    a.updated_at
from articles a
join article_versions v on v.id = a.current_version_id
order by a.updated_at desc
`)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list articles")
	}
	defer rows.Close()

	ret := []Article{}
	for rows.Next() {
		var row articleRow
		var article Article
		if err := rows.Scan(
			&row.ID,
			&row.VersionID,
			&row.Title,
			&row.Author,
			&row.Version,
			&row.Status,
			&row.Intro,
			&article.CreatedAt,
			&article.UpdatedAt,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan article row")
		}
		article.ID = row.ID.String()
		article.Title = row.Title
		article.Author = row.Author
		article.Version = row.Version
		article.Status = row.Status
		article.Intro = row.Intro

		sections, err := r.listSectionsForVersion(ctx, row.VersionID)
		if err != nil {
			return nil, err
		}
		article.Sections = sections
		ret = append(ret, article)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate articles")
	}

	return ret, nil
}

func (r *PostgresRepository) GetArticle(ctx context.Context, id string) (*Article, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	var row articleRow
	article := &Article{}
	err := r.pool.QueryRow(ctx, `
select
    a.id,
    a.current_version_id,
    a.title,
    a.author_display_name,
    v.version_label,
    a.status,
    v.intro,
    a.created_at,
    a.updated_at
from articles a
join article_versions v on v.id = a.current_version_id
where a.id = $1
`, id).Scan(
		&row.ID,
		&row.VersionID,
		&row.Title,
		&row.Author,
		&row.Version,
		&row.Status,
		&row.Intro,
		&article.CreatedAt,
		&article.UpdatedAt,
	)
	if err != nil {
		return nil, ErrNotFound
	}

	article.ID = row.ID.String()
	article.Title = row.Title
	article.Author = row.Author
	article.Version = row.Version
	article.Status = row.Status
	article.Intro = row.Intro

	sections, err := r.listSectionsForVersion(ctx, row.VersionID)
	if err != nil {
		return nil, err
	}
	article.Sections = sections

	return article, nil
}

func (r *PostgresRepository) listSectionsForVersion(ctx context.Context, versionID uuid.UUID) ([]Section, error) {
	rows, err := r.pool.Query(ctx, `
select id, title, body_plaintext
from article_sections
where article_version_id = $1
order by position asc
`, versionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list sections")
	}
	defer rows.Close()

	ret := []Section{}
	for rows.Next() {
		var id uuid.UUID
		var title string
		var body string
		if err := rows.Scan(&id, &title, &body); err != nil {
			return nil, errors.Wrap(err, "failed to scan section row")
		}
		ret = append(ret, Section{
			ID:         id.String(),
			Title:      title,
			Paragraphs: splitParagraphs(body),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate sections")
	}
	return ret, nil
}

func splitParagraphs(body string) []string {
	parts := strings.Split(body, "\n\n")
	ret := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		ret = append(ret, trimmed)
	}
	return ret
}
