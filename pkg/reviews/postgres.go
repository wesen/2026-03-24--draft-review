package reviews

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/go-go-golems/draft-review/pkg/reviewlinks"
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

func (r *PostgresRepository) StartSession(ctx context.Context, link *reviewlinks.ResolvedLink, input StartSessionInput) (*StartSessionResult, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin review session transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	readerName, readerEmail := resolveSessionReader(link, input)
	readerID := link.Reader.ID

	if link.InviteID != "" {
		var existing Session
		err = tx.QueryRow(ctx, `
select
    id,
    article_id,
    article_version_id,
    coalesce(nullif(reader_name, ''), ''),
    coalesce(reader_email, ''),
    is_anonymous,
    progress_percent,
    started_at,
    last_active_at,
    completed_at
from review_sessions
where reader_invite_id = $1
  and article_version_id = $2
order by started_at desc
limit 1
`, link.InviteID, link.ArticleVersionID).Scan(
			&existing.ID,
			&existing.ArticleID,
			&existing.ArticleVersionID,
			&existing.ReaderName,
			&existing.ReaderEmail,
			&existing.IsAnonymous,
			&existing.ProgressPercent,
			&existing.StartedAt,
			&existing.LastActiveAt,
			&existing.CompletedAt,
		)
		if err == nil {
			existing.ReaderID = readerID
			if existing.ReaderName == "" {
				existing.ReaderName = readerName
			}
			_, err = tx.Exec(ctx, `
update reader_invites
set opened_at = coalesce(opened_at, now())
where id = $1
`, link.InviteID)
			if err != nil {
				return nil, errors.Wrap(err, "failed to mark invite as opened")
			}
			if err := tx.Commit(ctx); err != nil {
				return nil, errors.Wrap(err, "failed to commit existing review session")
			}
			return &StartSessionResult{
				Session: &existing,
				Reader: reviewlinks.ReaderIdentity{
					ID:   readerID,
					Name: existing.ReaderName,
				},
				Article: link.Article,
			}, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(err, "failed to query existing review session")
		}
	}

	sessionID := uuid.New()
	var startedAt time.Time
	var lastActiveAt time.Time
	var completedAt sql.NullTime
	var inviteID any
	if link.InviteID != "" {
		inviteID = link.InviteID
	}

	err = tx.QueryRow(ctx, `
insert into review_sessions (
    id,
    article_id,
    article_version_id,
    reader_invite_id,
    reader_name,
    reader_email,
    is_anonymous
)
values ($1, $2, $3, $4, $5, nullif($6, ''), $7)
returning started_at, last_active_at, completed_at
`, sessionID, link.ArticleID, link.ArticleVersionID, inviteID, readerName, readerEmail, input.Anonymous).Scan(&startedAt, &lastActiveAt, &completedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create review session")
	}

	if link.InviteID != "" {
		_, err = tx.Exec(ctx, `
update reader_invites
set opened_at = coalesce(opened_at, now())
where id = $1
`, link.InviteID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to mark invite as opened")
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit review session transaction")
	}

	return &StartSessionResult{
		Session: &Session{
			ID:               sessionID.String(),
			ArticleID:        link.ArticleID,
			ArticleVersionID: link.ArticleVersionID,
			ReaderID:         readerID,
			ReaderName:       readerName,
			ReaderEmail:      readerEmail,
			IsAnonymous:      input.Anonymous,
			ProgressPercent:  0,
			StartedAt:        startedAt,
			LastActiveAt:     lastActiveAt,
			CompletedAt:      nullTimePtr(completedAt),
		},
		Reader: reviewlinks.ReaderIdentity{
			ID:   readerID,
			Name: readerName,
		},
		Article: link.Article,
	}, nil
}

func (r *PostgresRepository) RecordProgress(ctx context.Context, sessionID string, input ProgressInput) (*ProgressState, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin progress transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var articleVersionID uuid.UUID
	err = tx.QueryRow(ctx, `
select article_version_id
from review_sessions
where id = $1
for update
`, sessionID).Scan(&articleVersionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to load review session")
	}

	var sectionID uuid.UUID
	err = tx.QueryRow(ctx, `
select id
from article_sections
where id = $1
  and article_version_id = $2
`, input.SectionID, articleVersionID).Scan(&sectionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to resolve review section")
	}

	_, err = tx.Exec(ctx, `
insert into review_section_progress (
    review_session_id,
    section_id,
    reached_at,
    completed_at
)
values ($1, $2, now(), case when $3 then now() else null end)
on conflict (review_session_id, section_id) do update
set reached_at = least(review_section_progress.reached_at, excluded.reached_at),
    completed_at = case
        when excluded.completed_at is null then review_section_progress.completed_at
        else coalesce(review_section_progress.completed_at, excluded.completed_at)
    end
`, sessionID, sectionID, input.Completed)
	if err != nil {
		return nil, errors.Wrap(err, "failed to store section progress")
	}

	if input.ParagraphID != "" {
		_, err = tx.Exec(ctx, `
insert into review_paragraph_progress (
    review_session_id,
    section_id,
    paragraph_key,
    reached_at,
    completed_at
)
values ($1, $2, $3, now(), case when $4 then now() else null end)
on conflict (review_session_id, section_id, paragraph_key) do update
set reached_at = least(review_paragraph_progress.reached_at, excluded.reached_at),
    completed_at = case
        when excluded.completed_at is null then review_paragraph_progress.completed_at
        else coalesce(review_paragraph_progress.completed_at, excluded.completed_at)
    end
`, sessionID, sectionID, input.ParagraphID, input.Completed)
		if err != nil {
			return nil, errors.Wrap(err, "failed to store paragraph progress")
		}
	}

	progressPercent := 0
	if input.ProgressPercent != nil {
		progressPercent = *input.ProgressPercent
	} else {
		var reachedSections int
		var totalSections int
		err = tx.QueryRow(ctx, `
select
    (select count(*) from review_section_progress where review_session_id = $1),
    (select count(*) from article_sections where article_version_id = $2)
`, sessionID, articleVersionID).Scan(&reachedSections, &totalSections)
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate review progress")
		}
		if totalSections > 0 {
			progressPercent = (100*reachedSections + totalSections/2) / totalSections
		}
	}

	var lastActiveAt time.Time
	err = tx.QueryRow(ctx, `
update review_sessions
set last_active_at = now(),
    progress_percent = $2
where id = $1
returning last_active_at
`, sessionID, progressPercent).Scan(&lastActiveAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update session progress")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit progress transaction")
	}

	return &ProgressState{
		SessionID:       sessionID,
		ProgressPercent: progressPercent,
		LastActiveAt:    lastActiveAt,
	}, nil
}

func (r *PostgresRepository) AddReaction(ctx context.Context, sessionID string, input ReactionInput) (*Reaction, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin reaction transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var (
		articleID        uuid.UUID
		articleVersionID uuid.UUID
		inviteID         *uuid.UUID
		readerName       string
		readerEmail      string
		isAnonymous      bool
		requireNote      bool
	)

	err = tx.QueryRow(ctx, `
select
    s.article_id,
    s.article_version_id,
    s.reader_invite_id,
    coalesce(nullif(s.reader_name, ''), ''),
    coalesce(s.reader_email, ''),
    s.is_anonymous,
    a.require_note
from review_sessions s
join articles a on a.id = s.article_id
where s.id = $1
for update
`, sessionID).Scan(&articleID, &articleVersionID, &inviteID, &readerName, &readerEmail, &isAnonymous, &requireNote)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to load review session for reaction")
	}

	if requireNote && input.Text == "" {
		return nil, NewValidationError("this article requires a note for every reaction")
	}

	var sectionID uuid.UUID
	err = tx.QueryRow(ctx, `
select id
from article_sections
where id = $1
  and article_version_id = $2
`, input.SectionID, articleVersionID).Scan(&sectionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to resolve reaction section")
	}

	var reactionTypeKey string
	err = tx.QueryRow(ctx, `
select type_key
from article_reaction_types
where article_id = $1
  and type_key = $2
  and enabled = true
`, articleID, input.Type).Scan(&reactionTypeKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewValidationError("reaction type is not enabled for this article")
		}
		return nil, errors.Wrap(err, "failed to resolve reaction type")
	}

	reactionID := uuid.New()
	var createdAt time.Time
	err = tx.QueryRow(ctx, `
insert into reactions (
    id,
    review_session_id,
    article_id,
    article_version_id,
    section_id,
    paragraph_key,
    reaction_type_key,
    comment_text
)
values ($1, $2, $3, $4, $5, $6, $7, $8)
returning created_at
`, reactionID, sessionID, articleID, articleVersionID, sectionID, input.ParagraphID, reactionTypeKey, input.Text).Scan(&createdAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert reaction")
	}

	_, err = tx.Exec(ctx, `
update review_sessions
set last_active_at = now()
where id = $1
`, sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to bump review session after reaction")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit reaction transaction")
	}

	return &Reaction{
		ID:          reactionID.String(),
		ArticleID:   articleID.String(),
		SectionID:   sectionID.String(),
		ParagraphID: input.ParagraphID,
		ReaderID:    resolveReaderID(inviteID, sessionID),
		ReaderName:  resolveReactionReaderName(readerName, readerEmail, isAnonymous),
		Type:        reactionTypeKey,
		Text:        input.Text,
		CreatedAt:   createdAt,
	}, nil
}

func (r *PostgresRepository) SubmitSummary(ctx context.Context, sessionID string, input SummaryInput) (*Summary, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to begin summary transaction")
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var submittedAt time.Time
	err = tx.QueryRow(ctx, `
insert into review_summaries (
    review_session_id,
    overall_thoughts,
    recommendability,
    notify_new_version,
    submitted_at
)
values ($1, nullif($2, ''), nullif($3, ''), $4, now())
on conflict (review_session_id) do update
set overall_thoughts = excluded.overall_thoughts,
    recommendability = excluded.recommendability,
    notify_new_version = excluded.notify_new_version,
    submitted_at = now()
returning submitted_at
`, sessionID, input.OverallThoughts, input.Recommendability, input.NotifyNewVersion).Scan(&submittedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upsert review summary")
	}

	commandTag, err := tx.Exec(ctx, `
update review_sessions
set completed_at = coalesce(completed_at, now()),
    last_active_at = now(),
    progress_percent = 100
where id = $1
`, sessionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to complete review session")
	}
	if commandTag.RowsAffected() == 0 {
		return nil, ErrNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to commit summary transaction")
	}

	return &Summary{
		SessionID:        sessionID,
		OverallThoughts:  input.OverallThoughts,
		Recommendability: input.Recommendability,
		NotifyNewVersion: input.NotifyNewVersion,
		SubmittedAt:      submittedAt,
	}, nil
}

func resolveSessionReader(link *reviewlinks.ResolvedLink, input StartSessionInput) (string, string) {
	name := reviewlinks.DisplayNameFromInvite(link.InviteDisplayName, link.ReaderEmail, link.IdentityMode)
	email := strings.TrimSpace(link.ReaderEmail)

	if input.ReaderName != "" && !input.Anonymous {
		name = input.ReaderName
	}
	if name == "" {
		switch {
		case email != "":
			name = reviewlinks.DisplayNameFromEmail(email)
		case link.IdentityMode == reviewlinks.IdentityModeAnonymous:
			name = "Anonymous"
		default:
			name = "Guest Reader"
		}
	}
	if input.Anonymous || link.IdentityMode == reviewlinks.IdentityModeAnonymous {
		name = "Anonymous"
		email = ""
	}

	return name, email
}

func resolveReactionReaderName(readerName, readerEmail string, anonymous bool) string {
	if anonymous {
		return "Anonymous"
	}
	readerName = strings.TrimSpace(readerName)
	if readerName != "" {
		return readerName
	}
	if readerEmail != "" {
		return reviewlinks.DisplayNameFromEmail(readerEmail)
	}
	return "Reader"
}

func resolveReaderID(inviteID *uuid.UUID, sessionID string) string {
	if inviteID != nil {
		return inviteID.String()
	}
	return sessionID
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	ret := value.Time
	return &ret
}
