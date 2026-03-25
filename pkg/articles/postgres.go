package articles

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("article not found")

var defaultOwnerUserID = uuid.MustParse("11111111-1111-1111-1111-111111111111")

type PostgresRepository struct {
	pool *pgxpool.Pool
}

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

func (r *PostgresRepository) CreateArticle(ctx context.Context, input CreateArticleInput) (*Article, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin article creation transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := ensureDefaultOwner(ctx, tx); err != nil {
		return nil, err
	}

	articleID := uuid.New()
	versionID := uuid.New()
	sectionID := uuid.New()

	_, err = tx.Exec(ctx, `
insert into articles (
    id,
    owner_user_id,
    current_version_id,
    title,
    author_display_name,
    status,
    access_mode,
    show_author_note,
    reader_can_see_reactions,
    reader_can_see_names,
    require_note,
    allow_anonymous
)
values ($1, $2, $3, $4, $5, 'draft', 'invite_link', true, true, false, false, true)
`, articleID, defaultOwnerUserID, versionID, input.Title, input.Author)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert article")
	}

	_, err = tx.Exec(ctx, `
insert into article_versions (
    id,
    article_id,
    version_number,
    version_label,
    intro,
    author_note
)
values ($1, $2, 1, 'Draft 1', $3, '')
`, versionID, articleID, input.Intro)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert article version")
	}

	body := strings.Join([]string{""}, "\n\n")
	_, err = tx.Exec(ctx, `
insert into article_sections (
    id,
    article_version_id,
    section_key,
    position,
    title,
    body_markdown,
    body_plaintext,
    estimated_read_seconds
)
values ($1, $2, 'section-1', 1, 'Untitled Section', $3, $3, 0)
`, sectionID, versionID, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert default article section")
	}

	_, err = tx.Exec(ctx, `
insert into article_reaction_types (article_id, type_key, label, icon, is_default, enabled, position)
select
    $1,
    key,
    label,
    icon,
    true,
    true,
    position
from default_reaction_types
on conflict (article_id, type_key) do nothing
`, articleID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to seed article reaction types")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit article creation transaction")
	}

	return r.GetArticle(ctx, articleID.String())
}

func (r *PostgresRepository) UpdateArticle(ctx context.Context, id string, input UpdateArticleInput) (*Article, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin article update transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var articleID uuid.UUID
	var versionID uuid.UUID
	err = tx.QueryRow(ctx, `
select id, current_version_id
from articles
where id = $1
for update
`, id).Scan(&articleID, &versionID)
	if err != nil {
		return nil, ErrNotFound
	}

	title := input.Title
	author := input.Author
	status := input.Status
	if title != nil || author != nil || status != nil {
		_, err = tx.Exec(ctx, `
update articles
set
    title = coalesce($2, title),
    author_display_name = coalesce($3, author_display_name),
    status = coalesce($4, status),
    updated_at = now()
where id = $1
`, articleID, derefString(title), derefString(author), derefString(status))
		if err != nil {
			return nil, errors.Wrap(err, "failed to update article metadata")
		}
	}

	if input.Intro != nil {
		_, err = tx.Exec(ctx, `
update article_versions
set intro = $2
where id = $1
`, versionID, *input.Intro)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update article intro")
		}
	}

	if input.Sections != nil {
		if err := replaceSections(ctx, tx, versionID, *input.Sections); err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx, `
update articles
set updated_at = now()
where id = $1
`, articleID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to bump article timestamp after section update")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit article update transaction")
	}

	return r.GetArticle(ctx, articleID.String())
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

func ensureDefaultOwner(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
insert into users (id, email, name, password_hash, email_verified_at)
values ($1, 'manuel@example.com', 'Manuel', 'dev-only-password-hash', now())
on conflict (id) do nothing
`, defaultOwnerUserID)
	if err != nil {
		return errors.Wrap(err, "failed to ensure default owner user")
	}
	return nil
}

func replaceSections(ctx context.Context, tx pgx.Tx, versionID uuid.UUID, sections []SectionInput) error {
	_, err := tx.Exec(ctx, `
delete from article_sections
where article_version_id = $1
`, versionID)
	if err != nil {
		return errors.Wrap(err, "failed to delete existing sections")
	}

	for i, section := range sections {
		sectionID, err := parseOrGenerateSectionID(section.ID)
		if err != nil {
			return err
		}

		body := strings.Join(section.Paragraphs, "\n\n")
		_, err = tx.Exec(ctx, `
insert into article_sections (
    id,
    article_version_id,
    section_key,
    position,
    title,
    body_markdown,
    body_plaintext,
    estimated_read_seconds
)
values ($1, $2, $3, $4, $5, $6, $6, $7)
`, sectionID, versionID, fmt.Sprintf("section-%d", i+1), i+1, section.Title, body, estimateReadSeconds(body))
		if err != nil {
			return errors.Wrapf(err, "failed to insert section %d", i+1)
		}
	}

	return nil
}

func parseOrGenerateSectionID(raw string) (uuid.UUID, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return uuid.New(), nil
	}

	parsed, err := uuid.Parse(raw)
	if err == nil {
		return parsed, nil
	}

	if isTemporarySectionID(raw) {
		return uuid.New(), nil
	}

	return uuid.Nil, NewValidationError("section id must be a UUID when provided")
}

func isTemporarySectionID(value string) bool {
	return strings.HasPrefix(value, "s-new-")
}

var whitespaceRE = regexp.MustCompile(`\s+`)

func estimateReadSeconds(body string) int {
	normalized := strings.TrimSpace(whitespaceRE.ReplaceAllString(body, " "))
	if normalized == "" {
		return 0
	}

	words := len(strings.Split(normalized, " "))
	seconds := (words * 60) / 200
	if seconds < 15 {
		return 15
	}
	return seconds
}

func derefString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
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
	if len(ret) == 0 {
		return []string{""}
	}
	return ret
}
