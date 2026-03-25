package auth

import (
	"context"
	"database/sql"
	"strings"

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
