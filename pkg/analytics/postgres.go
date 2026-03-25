package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/go-go-golems/draft-review/pkg/reviewlinks"
	"github.com/go-go-golems/draft-review/pkg/reviews"
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

func (r *PostgresRepository) ListReaders(ctx context.Context, ownerUserID, articleID string) ([]reviewlinks.Reader, error) {
	if r == nil || r.pool == nil {
		return []reviewlinks.Reader{}, nil
	}

	rows, err := r.pool.Query(ctx, `
select
    i.id,
    i.email,
    i.article_id,
    i.invite_token,
    i.sent_at,
    coalesce(ls.reader_name, ''),
    coalesce(ls.progress_percent, 0),
    ls.last_active_at,
    coalesce(ls.is_anonymous, false)
from reader_invites i
join articles a on a.id = i.article_id
left join lateral (
    select reader_name, progress_percent, last_active_at, is_anonymous
    from review_sessions s
    where s.reader_invite_id = i.id
    order by s.started_at desc
    limit 1
) ls on true
where a.owner_user_id = $1
  and i.article_id = $2
  and i.revoked_at is null
order by i.sent_at desc
`, ownerUserID, articleID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list article readers")
	}
	defer rows.Close()

	ret := make([]reviewlinks.Reader, 0)
	for rows.Next() {
		var (
			readerID          uuid.UUID
			email             string
			resolvedArticleID uuid.UUID
			token             string
			invitedAt         time.Time
			sessionName       string
			progress          int
			lastActiveAt      sql.NullTime
			isAnonymous       bool
		)
		if err := rows.Scan(&readerID, &email, &resolvedArticleID, &token, &invitedAt, &sessionName, &progress, &lastActiveAt, &isAnonymous); err != nil {
			return nil, errors.Wrap(err, "failed to scan reader row")
		}

		name := reviewlinks.DisplayNameFromEmail(email)
		if sessionName != "" && !isAnonymous {
			name = sessionName
		}

		ret = append(ret, reviewlinks.Reader{
			ID:           readerID.String(),
			Name:         name,
			Email:        email,
			Avatar:       reviewlinks.AvatarFromName(name),
			ArticleID:    resolvedArticleID.String(),
			Progress:     progress,
			Token:        token,
			InvitedAt:    invitedAt,
			LastActiveAt: nullTimePtr(lastActiveAt),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate reader rows")
	}

	return ret, nil
}

func (r *PostgresRepository) ListReactions(ctx context.Context, ownerUserID, articleID string) ([]reviews.Reaction, error) {
	if r == nil || r.pool == nil {
		return []reviews.Reaction{}, nil
	}

	rows, err := r.pool.Query(ctx, `
select
    r.id,
    r.article_id,
    r.section_id,
    r.paragraph_key,
    s.reader_invite_id,
    s.id,
    coalesce(s.reader_name, ''),
    coalesce(s.reader_email, ''),
    s.is_anonymous,
    r.reaction_type_key,
    r.comment_text,
    r.created_at
from reactions r
join articles a on a.id = r.article_id
join review_sessions s on s.id = r.review_session_id
where a.owner_user_id = $1
  and r.article_id = $2
order by r.created_at desc
`, ownerUserID, articleID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list article reactions")
	}
	defer rows.Close()

	ret := make([]reviews.Reaction, 0)
	for rows.Next() {
		var (
			reactionID        uuid.UUID
			resolvedArticleID uuid.UUID
			sectionID         uuid.UUID
			paragraphKey      string
			inviteID          *uuid.UUID
			sessionID         uuid.UUID
			readerName        string
			readerEmail       string
			isAnonymous       bool
			reactionType      string
			commentText       string
			createdAt         time.Time
		)
		if err := rows.Scan(
			&reactionID,
			&resolvedArticleID,
			&sectionID,
			&paragraphKey,
			&inviteID,
			&sessionID,
			&readerName,
			&readerEmail,
			&isAnonymous,
			&reactionType,
			&commentText,
			&createdAt,
		); err != nil {
			return nil, errors.Wrap(err, "failed to scan reaction row")
		}

		name := readerName
		if isAnonymous {
			name = "Anonymous"
		} else if strings.TrimSpace(name) == "" {
			name = reviewlinks.DisplayNameFromEmail(readerEmail)
		}

		readerID := sessionID.String()
		if inviteID != nil {
			readerID = inviteID.String()
		}

		ret = append(ret, reviews.Reaction{
			ID:          reactionID.String(),
			ArticleID:   resolvedArticleID.String(),
			SectionID:   sectionID.String(),
			ParagraphID: paragraphKey,
			ReaderID:    readerID,
			ReaderName:  name,
			Type:        reactionType,
			Text:        commentText,
			CreatedAt:   createdAt,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate reaction rows")
	}

	return ret, nil
}

func (r *PostgresRepository) GetArticleAnalytics(ctx context.Context, ownerUserID, articleID string) (*ArticleAnalytics, error) {
	if r == nil || r.pool == nil {
		return nil, ErrNotFound
	}

	var (
		resolvedArticleID uuid.UUID
		currentVersionID  uuid.UUID
		versionLabel      string
	)
	err := r.pool.QueryRow(ctx, `
select a.id, a.current_version_id, v.version_label
from articles a
join article_versions v on v.id = a.current_version_id
where a.id = $1
  and a.owner_user_id = $2
`, articleID, ownerUserID).Scan(&resolvedArticleID, &currentVersionID, &versionLabel)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to resolve article analytics subject")
	}

	analytics := &ArticleAnalytics{
		ArticleID: resolvedArticleID.String(),
		Version:   versionLabel,
		Sections:  make([]SectionAnalytics, 0),
	}

	err = r.pool.QueryRow(ctx, `
select
    count(*)::int,
    count(*) filter (where progress_percent > 0 or completed_at is not null)::int,
    coalesce(round(avg(progress_percent))::int, 0),
    (select count(*)::int from reactions where article_id = $1 and article_version_id = $2)
from review_sessions
where article_id = $1
  and article_version_id = $2
`, resolvedArticleID, currentVersionID).Scan(
		&analytics.ReaderCount,
		&analytics.EngagedReaderCount,
		&analytics.AverageProgress,
		&analytics.TotalReactions,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query article analytics summary")
	}

	rows, err := r.pool.Query(ctx, `
select
    s.id,
    s.title,
    s.position,
    count(distinct rsp.review_session_id)::int as readers_reached,
    count(distinct case when rsp.completed_at is not null then rsp.review_session_id end)::int as readers_completed,
    count(r.id)::int as reaction_count
from article_sections s
left join review_section_progress rsp on rsp.section_id = s.id
left join reactions r on r.section_id = s.id and r.article_version_id = $1
where s.article_version_id = $1
group by s.id, s.title, s.position
order by s.position asc
`, currentVersionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list section analytics")
	}
	defer rows.Close()

	sectionsByID := map[string]*SectionAnalytics{}
	for rows.Next() {
		var (
			sectionID        uuid.UUID
			title            string
			position         int
			readersReached   int
			readersCompleted int
			reactionCount    int
		)
		if err := rows.Scan(&sectionID, &title, &position, &readersReached, &readersCompleted, &reactionCount); err != nil {
			return nil, errors.Wrap(err, "failed to scan section analytics")
		}
		section := SectionAnalytics{
			SectionID:        sectionID.String(),
			Title:            title,
			Position:         position,
			ReadersReached:   readersReached,
			ReadersCompleted: readersCompleted,
			ReactionCount:    reactionCount,
			ReactionCounts:   map[string]int{},
		}
		analytics.Sections = append(analytics.Sections, section)
		sectionsByID[section.SectionID] = &analytics.Sections[len(analytics.Sections)-1]
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate section analytics")
	}

	typeRows, err := r.pool.Query(ctx, `
select section_id, reaction_type_key, count(*)::int
from reactions
where article_id = $1
  and article_version_id = $2
group by section_id, reaction_type_key
`, resolvedArticleID, currentVersionID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list reaction type counts")
	}
	defer typeRows.Close()

	for typeRows.Next() {
		var (
			sectionID    uuid.UUID
			reactionType string
			count        int
		)
		if err := typeRows.Scan(&sectionID, &reactionType, &count); err != nil {
			return nil, errors.Wrap(err, "failed to scan reaction type counts")
		}
		if section, ok := sectionsByID[sectionID.String()]; ok {
			section.ReactionCounts[reactionType] = count
		}
	}
	if err := typeRows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate reaction type counts")
	}

	analytics.DraftKiller = computeDraftKiller(analytics.EngagedReaderCount, analytics.Sections)
	return analytics, nil
}

func (r *PostgresRepository) ListReaderDirectory(ctx context.Context, ownerUserID string) ([]ReaderContact, error) {
	if r == nil || r.pool == nil {
		return []ReaderContact{}, nil
	}

	rows, err := r.pool.Query(ctx, `
select
    i.email,
    a.id,
    a.title,
    coalesce(ls.reader_name, ''),
    coalesce(ls.progress_percent, 0),
    ls.last_active_at,
    coalesce(ls.is_anonymous, false),
    coalesce(rc.reaction_count, 0)
from reader_invites i
join articles a on a.id = i.article_id
left join lateral (
    select reader_name, progress_percent, last_active_at, is_anonymous
    from review_sessions s
    where s.reader_invite_id = i.id
    order by s.started_at desc
    limit 1
) ls on true
left join lateral (
    select count(*)::int as reaction_count
    from review_sessions s
    join reactions r on r.review_session_id = s.id
    where s.reader_invite_id = i.id
) rc on true
where a.owner_user_id = $1
  and i.revoked_at is null
order by lower(i.email), a.updated_at desc
`, ownerUserID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list reader directory")
	}
	defer rows.Close()

	directory := map[string]*ReaderContact{}
	for rows.Next() {
		var (
			email         string
			articleID     uuid.UUID
			title         string
			sessionName   string
			progress      int
			lastActiveAt  sql.NullTime
			isAnonymous   bool
			reactionCount int
		)
		if err := rows.Scan(&email, &articleID, &title, &sessionName, &progress, &lastActiveAt, &isAnonymous, &reactionCount); err != nil {
			return nil, errors.Wrap(err, "failed to scan reader directory row")
		}

		key := strings.ToLower(strings.TrimSpace(email))
		contact, ok := directory[key]
		if !ok {
			name := reviewlinks.DisplayNameFromEmail(email)
			if sessionName != "" && !isAnonymous {
				name = sessionName
			}
			contact = &ReaderContact{
				Email:    email,
				Name:     name,
				Avatar:   reviewlinks.AvatarFromName(name),
				Status:   "invited",
				Articles: make([]ReaderArticleActivity, 0),
			}
			directory[key] = contact
		}

		contact.TotalReactions += reactionCount
		contact.Articles = append(contact.Articles, ReaderArticleActivity{
			ArticleID:     articleID.String(),
			Title:         title,
			Progress:      progress,
			ReactionCount: reactionCount,
			LastActiveAt:  nullTimePtr(lastActiveAt),
		})

		if !containsArticle(contact.Articles, articleID.String()) {
			contact.TotalArticles++
		}

		if lastActiveAt.Valid {
			if contact.LastActiveAt == nil || lastActiveAt.Time.After(*contact.LastActiveAt) {
				t := lastActiveAt.Time
				contact.LastActiveAt = &t
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate reader directory rows")
	}

	ret := make([]ReaderContact, 0, len(directory))
	now := time.Now().UTC()
	for _, contact := range directory {
		contact.TotalArticles = uniqueArticleCount(contact.Articles)
		switch {
		case contact.LastActiveAt == nil:
			contact.Status = "invited"
		case now.Sub(*contact.LastActiveAt) <= 7*24*time.Hour:
			contact.Status = "active"
		default:
			contact.Status = "stale"
		}
		sort.Slice(contact.Articles, func(i, j int) bool {
			left := contact.Articles[i].LastActiveAt
			right := contact.Articles[j].LastActiveAt
			if left == nil {
				return false
			}
			if right == nil {
				return true
			}
			return left.After(*right)
		})
		ret = append(ret, *contact)
	}

	sort.Slice(ret, func(i, j int) bool {
		return strings.ToLower(ret[i].Email) < strings.ToLower(ret[j].Email)
	})

	return ret, nil
}

func computeDraftKiller(readerCount int, sections []SectionAnalytics) *DraftKillerInsight {
	if readerCount == 0 || len(sections) == 0 {
		return nil
	}

	var killer *DraftKillerInsight
	for _, section := range sections {
		dropOff := readerCount - section.ReadersReached
		if dropOff < 0 {
			dropOff = 0
		}
		if killer == nil || dropOff > killer.DropOffCount {
			rate := 0
			if readerCount > 0 {
				rate = (100*dropOff + readerCount/2) / readerCount
			}
			killer = &DraftKillerInsight{
				SectionID:    section.SectionID,
				Title:        section.Title,
				DropOffCount: dropOff,
				DropOffRate:  rate,
				Summary:      fmt.Sprintf("%d of %d engaged readers dropped off before or inside %q.", dropOff, readerCount, section.Title),
			}
		}
	}

	if killer != nil && killer.DropOffCount == 0 {
		return nil
	}

	return killer
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}

func uniqueArticleCount(articles []ReaderArticleActivity) int {
	seen := map[string]struct{}{}
	for _, article := range articles {
		seen[article.ArticleID] = struct{}{}
	}
	return len(seen)
}

func containsArticle(articles []ReaderArticleActivity, articleID string) bool {
	for _, article := range articles {
		if article.ArticleID == articleID {
			return true
		}
	}
	return false
}
