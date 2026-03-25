package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

type fakeRepository struct {
	foundUser   *User
	findErr     error
	createInput AuthenticatedIdentity
	updateInput AuthenticatedIdentity
	createUser  *User
	updateUser  *User
}

func (f *fakeRepository) FindByAuthIdentity(ctx context.Context, issuer, subject string) (*User, error) {
	return f.foundUser, f.findErr
}

func (f *fakeRepository) CreateAuthenticatedUser(ctx context.Context, identity AuthenticatedIdentity) (*User, error) {
	f.createInput = identity
	return f.createUser, nil
}

func (f *fakeRepository) UpdateAuthenticatedUser(ctx context.Context, userID uuid.UUID, identity AuthenticatedIdentity) (*User, error) {
	f.updateInput = identity
	return f.updateUser, nil
}

func TestEnsureAuthenticatedUserCreatesMissingUser(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{
		findErr: ErrNotFound,
		createUser: &User{
			ID:    uuid.NewString(),
			Email: "author@example.com",
			Name:  "Author Name",
		},
	}
	service := NewService(repo)

	user, err := service.EnsureAuthenticatedUser(context.Background(), AuthenticatedIdentity{
		Issuer:        "issuer",
		Subject:       "subject",
		Email:         "Author@example.com",
		DisplayName:   "Author Name",
		EmailVerified: true,
	})
	if err != nil {
		t.Fatalf("EnsureAuthenticatedUser returned error: %v", err)
	}

	if user.Email != "author@example.com" {
		t.Fatalf("expected normalized email author@example.com, got %q", user.Email)
	}
	if repo.createInput.Subject != "subject" {
		t.Fatalf("expected create subject subject, got %q", repo.createInput.Subject)
	}
	if repo.createInput.Email != "author@example.com" {
		t.Fatalf("expected normalized create email, got %q", repo.createInput.Email)
	}
}

func TestEnsureAuthenticatedUserUpdatesExistingUser(t *testing.T) {
	t.Parallel()

	existingID := uuid.NewString()
	repo := &fakeRepository{
		foundUser: &User{ID: existingID, Email: "author@example.com", Name: "Existing Author"},
		updateUser: &User{
			ID:    existingID,
			Email: "author@example.com",
			Name:  "Existing Author",
		},
	}
	service := NewService(repo)

	_, err := service.EnsureAuthenticatedUser(context.Background(), AuthenticatedIdentity{
		Issuer:      "issuer",
		Subject:     "subject",
		Email:       "author@example.com",
		DisplayName: "Existing Author",
	})
	if err != nil {
		t.Fatalf("EnsureAuthenticatedUser returned error: %v", err)
	}

	if repo.updateInput.Subject != "subject" {
		t.Fatalf("expected update subject subject, got %q", repo.updateInput.Subject)
	}
}
