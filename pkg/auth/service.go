package auth

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrNotFound = errors.New("user not found")

type Repository interface {
	FindByAuthIdentity(ctx context.Context, issuer, subject string) (*User, error)
	CreateAuthenticatedUser(ctx context.Context, identity AuthenticatedIdentity) (*User, error)
	UpdateAuthenticatedUser(ctx context.Context, userID uuid.UUID, identity AuthenticatedIdentity) (*User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) EnsureAuthenticatedUser(ctx context.Context, identity AuthenticatedIdentity) (*User, error) {
	if s == nil || s.repo == nil {
		return nil, errors.New("auth service repository is not configured")
	}

	identity = normalizeIdentity(identity)
	if identity.Subject == "" {
		return nil, errors.New("authenticated identity subject is required")
	}

	user, err := s.repo.FindByAuthIdentity(ctx, identity.Issuer, identity.Subject)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return nil, err
		}

		user, err = s.repo.CreateAuthenticatedUser(ctx, identity)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	parsedUserID, err := uuid.Parse(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse authenticated user id")
	}

	return s.repo.UpdateAuthenticatedUser(ctx, parsedUserID, identity)
}

func normalizeIdentity(identity AuthenticatedIdentity) AuthenticatedIdentity {
	identity.Issuer = strings.TrimSpace(identity.Issuer)
	identity.Subject = strings.TrimSpace(identity.Subject)
	identity.Email = strings.ToLower(strings.TrimSpace(identity.Email))
	identity.DisplayName = strings.TrimSpace(identity.DisplayName)
	if identity.DisplayName == "" {
		switch {
		case identity.Email != "":
			identity.DisplayName = identity.Email
		case identity.Subject != "":
			identity.DisplayName = identity.Subject
		default:
			identity.DisplayName = "Author"
		}
	}
	if identity.Email == "" {
		subjectSlug := strings.NewReplacer(":", "-", "/", "-", " ", "-").Replace(identity.Subject)
		if subjectSlug == "" {
			subjectSlug = "unknown"
		}
		identity.Email = subjectSlug + "@draft-review.local"
	}
	return identity
}
