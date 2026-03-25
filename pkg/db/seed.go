package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

var (
	devUserID           = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	devArticleID        = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	devArticleVersionID = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	devSectionID        = uuid.MustParse("44444444-4444-4444-4444-444444444444")
)

func SeedDev(ctx context.Context, pool *pgxpool.Pool) error {
	if pool == nil {
		return ErrNotConfigured
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to begin dev seed transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, `
insert into users (id, email, name, password_hash, email_verified_at)
values ($1, 'manuel@example.com', 'Manuel', 'dev-only-password-hash', now())
on conflict (id) do nothing
`, devUserID)
	if err != nil {
		return errors.Wrap(err, "failed to seed user")
	}

	_, err = tx.Exec(ctx, `
insert into articles (
	id,
	owner_user_id,
	title,
	author_display_name,
	status,
	access_mode,
	share_token,
	show_author_note,
	reader_can_see_reactions,
	reader_can_see_names,
	require_note,
	allow_anonymous
)
values (
	$1,
	$2,
	'Why Design Systems Fail',
	'Manuel',
	'in_review',
	'invite_link',
	'tok-dev-article-1',
	true,
	true,
	false,
	false,
	true
)
on conflict (id) do nothing
`, devArticleID, devUserID)
	if err != nil {
		return errors.Wrap(err, "failed to seed article")
	}

	_, err = tx.Exec(ctx, `
insert into article_versions (id, article_id, version_number, version_label, intro, author_note)
values (
	$1,
	$2,
	1,
	'Draft 1',
	'Thanks for reading an early draft!',
	'I am especially looking for feedback on what feels useful, confusing, or slow.'
)
on conflict (id) do nothing
`, devArticleVersionID, devArticleID)
	if err != nil {
		return errors.Wrap(err, "failed to seed article version")
	}

	_, err = tx.Exec(ctx, `
update articles
set current_version_id = $2
where id = $1 and current_version_id is distinct from $2
`, devArticleID, devArticleVersionID)
	if err != nil {
		return errors.Wrap(err, "failed to update article current version")
	}

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
values (
	$1,
	$2,
	'introduction',
	1,
	'Introduction',
	'Design systems promise consistency and speed, but many teams abandon them within a year.',
	'Design systems promise consistency and speed, but many teams abandon them within a year.',
	45
)
on conflict (id) do nothing
`, devSectionID, devArticleVersionID)
	if err != nil {
		return errors.Wrap(err, "failed to seed section")
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
`, devArticleID)
	if err != nil {
		return errors.Wrap(err, "failed to seed article reaction types")
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "failed to commit dev seed transaction")
	}

	return nil
}
