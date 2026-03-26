package auth

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) FindByAuthIdentity(ctx context.Context, issuer, subject string) (*User, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("postgres pool is not configured")
	}

	row := r.pool.QueryRow(ctx, `
select id, auth_subject, auth_issuer, email, name, email_verified_at, created_at, updated_at
from users
where auth_issuer = $1 and auth_subject = $2
`, issuer, subject)

	user := &User{}
	var authSubject sql.NullString
	var authIssuer sql.NullString
	var emailVerifiedAt sql.NullTime
	if err := row.Scan(
		&user.ID,
		&authSubject,
		&authIssuer,
		&user.Email,
		&user.Name,
		&emailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to load user by auth identity")
	}

	assignNullableUserFields(user, authSubject, authIssuer, emailVerifiedAt)
	return user, nil
}

func (r *PostgresRepository) CreateAuthenticatedUser(ctx context.Context, identity AuthenticatedIdentity) (*User, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("postgres pool is not configured")
	}

	user := &User{}
	row := r.pool.QueryRow(ctx, `
insert into users(id, auth_subject, auth_issuer, email, name, password_hash, email_verified_at)
values($1, $2, $3, $4, $5, 'oidc-managed', case when $6 then now() else null end)
on conflict (email) do update
set auth_subject = excluded.auth_subject,
    auth_issuer = excluded.auth_issuer,
    name = case
      when coalesce(nullif(trim(users.name), ''), '') = '' then excluded.name
      else users.name
    end,
    email_verified_at = case
      when excluded.email_verified_at is not null then coalesce(users.email_verified_at, excluded.email_verified_at)
      else users.email_verified_at
    end,
    updated_at = now()
where users.auth_subject is null
   or (users.auth_subject = excluded.auth_subject and users.auth_issuer = excluded.auth_issuer)
returning id, auth_subject, auth_issuer, email, name, email_verified_at, created_at, updated_at
`, uuid.New(), identity.Subject, identity.Issuer, identity.Email, identity.DisplayName, identity.EmailVerified)

	var authSubject sql.NullString
	var authIssuer sql.NullString
	var emailVerifiedAt sql.NullTime
	if err := row.Scan(
		&user.ID,
		&authSubject,
		&authIssuer,
		&user.Email,
		&user.Name,
		&emailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "failed to create authenticated user")
	}

	assignNullableUserFields(user, authSubject, authIssuer, emailVerifiedAt)
	return user, nil
}

func (r *PostgresRepository) UpdateAuthenticatedUser(ctx context.Context, userID uuid.UUID, identity AuthenticatedIdentity) (*User, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("postgres pool is not configured")
	}

	user := &User{}
	row := r.pool.QueryRow(ctx, `
update users
set name = case
      when coalesce(nullif(trim(name), ''), '') = '' and nullif($2::text, '') is not null then $2::text
      else name
    end,
    email = case
      when coalesce(nullif(trim(email), ''), '') = '' and nullif($3::text, '') is not null then $3::text
      else email
    end,
    email_verified_at = case
      when $4 then coalesce(email_verified_at, now())
      else email_verified_at
    end,
    updated_at = now()
where id = $1
returning id, auth_subject, auth_issuer, email, name, email_verified_at, created_at, updated_at
`, userID, identity.DisplayName, identity.Email, identity.EmailVerified)

	var authSubject sql.NullString
	var authIssuer sql.NullString
	var emailVerifiedAt sql.NullTime
	if err := row.Scan(
		&user.ID,
		&authSubject,
		&authIssuer,
		&user.Email,
		&user.Name,
		&emailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "failed to update authenticated user")
	}

	assignNullableUserFields(user, authSubject, authIssuer, emailVerifiedAt)
	return user, nil
}

func (r *PostgresRepository) CreateAuthorSession(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) (*Session, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("postgres pool is not configured")
	}

	session := &Session{}
	var revokedAt sql.NullTime
	row := r.pool.QueryRow(ctx, `
insert into author_sessions(id, user_id, token_hash, expires_at)
values($1, $2, $3, $4)
returning id, user_id, token_hash, expires_at, created_at, last_used_at, revoked_at
`, uuid.New(), userID, tokenHash, expiresAt)

	if err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.LastUsedAt,
		&revokedAt,
	); err != nil {
		return nil, errors.Wrap(err, "failed to create author session")
	}
	if revokedAt.Valid {
		value := revokedAt.Time.UTC()
		session.RevokedAt = &value
	}

	return session, nil
}

func (r *PostgresRepository) FindAuthorSessionByTokenHash(ctx context.Context, tokenHash string) (*ResolvedSession, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("postgres pool is not configured")
	}

	row := r.pool.QueryRow(ctx, `
select
    s.id,
    s.user_id,
    s.token_hash,
    s.expires_at,
    s.created_at,
    s.last_used_at,
    s.revoked_at,
    u.id,
    u.auth_subject,
    u.auth_issuer,
    u.email,
    u.name,
    u.email_verified_at,
    u.created_at,
    u.updated_at
from author_sessions s
join users u on u.id = s.user_id
where s.token_hash = $1
`, tokenHash)

	resolved := &ResolvedSession{}
	var sessionRevokedAt sql.NullTime
	var authSubject sql.NullString
	var authIssuer sql.NullString
	var emailVerifiedAt sql.NullTime
	if err := row.Scan(
		&resolved.Session.ID,
		&resolved.Session.UserID,
		&resolved.Session.TokenHash,
		&resolved.Session.ExpiresAt,
		&resolved.Session.CreatedAt,
		&resolved.Session.LastUsedAt,
		&sessionRevokedAt,
		&resolved.User.ID,
		&authSubject,
		&authIssuer,
		&resolved.User.Email,
		&resolved.User.Name,
		&emailVerifiedAt,
		&resolved.User.CreatedAt,
		&resolved.User.UpdatedAt,
	); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to load author session")
	}
	if sessionRevokedAt.Valid {
		value := sessionRevokedAt.Time.UTC()
		resolved.Session.RevokedAt = &value
	}
	assignNullableUserFields(&resolved.User, authSubject, authIssuer, emailVerifiedAt)
	return resolved, nil
}

func (r *PostgresRepository) TouchAuthorSession(ctx context.Context, sessionID uuid.UUID, touchedAt time.Time, renewedExpiresAt *time.Time) (*Session, error) {
	if r == nil || r.pool == nil {
		return nil, errors.New("postgres pool is not configured")
	}

	session := &Session{}
	var revokedAt sql.NullTime
	var renewedAt any = nil
	if renewedExpiresAt != nil {
		renewedAt = renewedExpiresAt.UTC()
	}
	row := r.pool.QueryRow(ctx, `
update author_sessions
set last_used_at = $2,
    expires_at = coalesce($3, expires_at)
where id = $1
returning id, user_id, token_hash, expires_at, created_at, last_used_at, revoked_at
`, sessionID, touchedAt.UTC(), renewedAt)

	if err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.LastUsedAt,
		&revokedAt,
	); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "failed to touch author session")
	}
	if revokedAt.Valid {
		value := revokedAt.Time.UTC()
		session.RevokedAt = &value
	}

	return session, nil
}

func (r *PostgresRepository) RevokeAuthorSessionByTokenHash(ctx context.Context, tokenHash string) error {
	if r == nil || r.pool == nil {
		return errors.New("postgres pool is not configured")
	}

	tag, err := r.pool.Exec(ctx, `
update author_sessions
set revoked_at = coalesce(revoked_at, now())
where token_hash = $1
`, tokenHash)
	if err != nil {
		return errors.Wrap(err, "failed to revoke author session")
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func assignNullableUserFields(user *User, authSubject, authIssuer sql.NullString, emailVerifiedAt sql.NullTime) {
	if user == nil {
		return
	}

	if authSubject.Valid {
		user.AuthSubject = authSubject.String
	}
	if authIssuer.Valid {
		user.AuthIssuer = authIssuer.String
	}
	if emailVerifiedAt.Valid {
		value := emailVerifiedAt.Time.UTC()
		user.EmailVerifiedAt = &value
	}
}
