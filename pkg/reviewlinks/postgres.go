package reviewlinks

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) ResetShareToken(ctx context.Context, ownerUserID, articleID string) (*ShareLink, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	token := fmt.Sprintf("share-%s-%d", strings.ReplaceAll(uuid.NewString(), "-", ""), time.Now().UTC().Unix())
	commandTag, err := r.pool.Exec(ctx, `
update articles
set share_token = $3,
    updated_at = now()
where id = $1
  and owner_user_id = $2
`, articleID, ownerUserID, token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to reset article share token")
	}
	if commandTag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}

	return &ShareLink{
		Token: token,
		URL:   "/r/" + token,
	}, nil
}

func (r *PostgresRepository) CreateInvite(ctx context.Context, ownerUserID, articleID string, input InviteInput) (*Reader, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin invite transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var verifiedArticleID uuid.UUID
	err = tx.QueryRow(ctx, `
select id
from articles
where id = $1
  and owner_user_id = $2
for update
`, articleID, ownerUserID).Scan(&verifiedArticleID)
	if err != nil {
		return nil, ErrNotFound
	}

	inviteID := uuid.New()
	inviteToken := fmt.Sprintf("invite-%s-%d", strings.ReplaceAll(uuid.NewString(), "-", ""), time.Now().UTC().Unix())

	_, err = tx.Exec(ctx, `
insert into reader_invites (
    id,
    article_id,
    email,
    display_name,
    identity_mode,
    is_preview,
    invite_token,
    invite_note
)
values ($1, $2, nullif($3, ''), nullif($4, ''), $5, $6, $7, $8)
`, inviteID, verifiedArticleID, input.Email, input.DisplayName, input.IdentityMode, input.IsPreview, inviteToken, input.Note)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create reader invite")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit invite transaction")
	}

	name := DisplayNameFromInvite(input.DisplayName, input.Email, input.IdentityMode)
	return &Reader{
		ID:           inviteID.String(),
		Name:         name,
		Email:        input.Email,
		Avatar:       AvatarFromName(name),
		ArticleID:    verifiedArticleID.String(),
		Progress:     0,
		Token:        inviteToken,
		IdentityMode: input.IdentityMode,
		IsPreview:    input.IsPreview,
		InvitedAt:    time.Now().UTC(),
	}, nil
}

func (r *PostgresRepository) ResolveToken(ctx context.Context, token string) (*ResolvedLink, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	link, err := r.resolveInviteToken(ctx, token)
	if err == nil {
		return link, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	return r.resolveShareToken(ctx, token)
}

func (r *PostgresRepository) resolveInviteToken(ctx context.Context, token string) (*ResolvedLink, error) {
	type inviteRow struct {
		InviteID              uuid.UUID
		Email                 sql.NullString
		DisplayName           sql.NullString
		IdentityMode          string
		IsPreview             bool
		ArticleID             uuid.UUID
		VersionID             uuid.UUID
		Title                 string
		Author                string
		Version               string
		Intro                 string
		AccessMode            string
		AllowAnonymous        bool
		RequireNote           bool
		ReaderCanSeeReactions bool
		ReaderCanSeeNames     bool
		ShowAuthorNote        bool
	}

	var row inviteRow
	err := r.pool.QueryRow(ctx, `
	select
	    i.id,
	    i.email,
	    i.display_name,
	    i.identity_mode,
	    i.is_preview,
	    a.id,
	    a.current_version_id,
    a.title,
    a.author_display_name,
    v.version_label,
    v.intro,
    a.access_mode,
    a.allow_anonymous,
    a.require_note,
    a.reader_can_see_reactions,
    a.reader_can_see_names,
    a.show_author_note
from reader_invites i
join articles a on a.id = i.article_id
join article_versions v on v.id = a.current_version_id
where i.invite_token = $1
  and i.revoked_at is null
`, token).Scan(
		&row.InviteID,
		&row.Email,
		&row.DisplayName,
		&row.IdentityMode,
		&row.IsPreview,
		&row.ArticleID,
		&row.VersionID,
		&row.Title,
		&row.Author,
		&row.Version,
		&row.Intro,
		&row.AccessMode,
		&row.AllowAnonymous,
		&row.RequireNote,
		&row.ReaderCanSeeReactions,
		&row.ReaderCanSeeNames,
		&row.ShowAuthorNote,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to resolve invite token")
	}

	sections, err := r.listSectionsForVersion(ctx, row.VersionID)
	if err != nil {
		return nil, err
	}

	name := DisplayNameFromInvite(row.DisplayName.String, row.Email.String, row.IdentityMode)
	return &ResolvedLink{
		Token:                 token,
		ArticleID:             row.ArticleID.String(),
		ArticleVersionID:      row.VersionID.String(),
		InviteID:              row.InviteID.String(),
		AccessMode:            row.AccessMode,
		ReaderEmail:           row.Email.String,
		InviteDisplayName:     row.DisplayName.String,
		IdentityMode:          row.IdentityMode,
		IsPreview:             row.IsPreview,
		AllowAnonymous:        row.AllowAnonymous,
		RequireNote:           row.RequireNote,
		ReaderCanSeeReactions: row.ReaderCanSeeReactions,
		ReaderCanSeeNames:     row.ReaderCanSeeNames,
		ShowAuthorNote:        row.ShowAuthorNote,
		Reader: ReaderIdentity{
			ID:   row.InviteID.String(),
			Name: name,
		},
		Article: ReaderArticle{
			ID:       row.ArticleID.String(),
			Title:    row.Title,
			Author:   row.Author,
			Version:  row.Version,
			Intro:    row.Intro,
			Sections: sections,
		},
	}, nil
}

func (r *PostgresRepository) resolveShareToken(ctx context.Context, token string) (*ResolvedLink, error) {
	type shareRow struct {
		ArticleID             uuid.UUID
		VersionID             uuid.UUID
		Title                 string
		Author                string
		Version               string
		Intro                 string
		AccessMode            string
		AllowAnonymous        bool
		RequireNote           bool
		ReaderCanSeeReactions bool
		ReaderCanSeeNames     bool
		ShowAuthorNote        bool
	}

	var row shareRow
	err := r.pool.QueryRow(ctx, `
select
    a.id,
    a.current_version_id,
    a.title,
    a.author_display_name,
    v.version_label,
    v.intro,
    a.access_mode,
    a.allow_anonymous,
    a.require_note,
    a.reader_can_see_reactions,
    a.reader_can_see_names,
    a.show_author_note
from articles a
join article_versions v on v.id = a.current_version_id
where a.share_token = $1
`, token).Scan(
		&row.ArticleID,
		&row.VersionID,
		&row.Title,
		&row.Author,
		&row.Version,
		&row.Intro,
		&row.AccessMode,
		&row.AllowAnonymous,
		&row.RequireNote,
		&row.ReaderCanSeeReactions,
		&row.ReaderCanSeeNames,
		&row.ShowAuthorNote,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to resolve share token")
	}

	sections, err := r.listSectionsForVersion(ctx, row.VersionID)
	if err != nil {
		return nil, err
	}

	return &ResolvedLink{
		Token:                 token,
		ArticleID:             row.ArticleID.String(),
		ArticleVersionID:      row.VersionID.String(),
		AccessMode:            row.AccessMode,
		AllowAnonymous:        row.AllowAnonymous,
		RequireNote:           row.RequireNote,
		ReaderCanSeeReactions: row.ReaderCanSeeReactions,
		ReaderCanSeeNames:     row.ReaderCanSeeNames,
		ShowAuthorNote:        row.ShowAuthorNote,
		Reader: ReaderIdentity{
			ID:   "share:" + row.ArticleID.String(),
			Name: "Guest Reader",
		},
		Article: ReaderArticle{
			ID:       row.ArticleID.String(),
			Title:    row.Title,
			Author:   row.Author,
			Version:  row.Version,
			Intro:    row.Intro,
			Sections: sections,
		},
	}, nil
}

func (r *PostgresRepository) listSectionsForVersion(ctx context.Context, versionID uuid.UUID) ([]Section, error) {
	rows, err := r.pool.Query(ctx, `
select id, title, body_markdown
from article_sections
where article_version_id = $1
order by position asc
`, versionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list reader sections")
	}
	defer rows.Close()

	ret := make([]Section, 0)
	for rows.Next() {
		var (
			sectionID uuid.UUID
			title     string
			body      string
		)
		if err := rows.Scan(&sectionID, &title, &body); err != nil {
			return nil, errors.Wrap(err, "failed to scan reader section")
		}
		ret = append(ret, Section{
			ID:           sectionID.String(),
			Title:        title,
			BodyMarkdown: body,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate reader sections")
	}

	return ret, nil
}
