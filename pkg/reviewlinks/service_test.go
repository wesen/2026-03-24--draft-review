package reviewlinks

import (
	"context"
	"testing"
)

type fakeRepository struct {
	createInviteErr error
}

func (f *fakeRepository) ResetShareToken(ctx context.Context, ownerUserID, articleID string) (*ShareLink, error) {
	return &ShareLink{Token: "share-1", URL: "/r/share-1"}, nil
}

func (f *fakeRepository) CreateInvite(ctx context.Context, ownerUserID, articleID string, input InviteInput) (*Reader, error) {
	if f.createInviteErr != nil {
		return nil, f.createInviteErr
	}
	return &Reader{ID: "reader-1", Email: input.Email}, nil
}

func (f *fakeRepository) ResolveToken(ctx context.Context, token string) (*ResolvedLink, error) {
	return &ResolvedLink{}, nil
}

func TestCreateInviteRejectsInvalidEmail(t *testing.T) {
	t.Parallel()

	service := NewService(&fakeRepository{})
	_, err := service.CreateInvite(context.Background(), "owner-1", "article-1", InviteInput{
		Email: "not-an-email",
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}

func TestCreateInviteAllowsNamedInviteWithoutEmail(t *testing.T) {
	t.Parallel()

	service := NewService(&fakeRepository{})
	_, err := service.CreateInvite(context.Background(), "owner-1", "article-1", InviteInput{
		IdentityMode: IdentityModeNamed,
		DisplayName:  "Workshop Reader",
	})
	if err != nil {
		t.Fatalf("expected named invite without email to succeed, got %v", err)
	}
}

func TestCreateInviteAllowsAnonymousInviteWithoutEmail(t *testing.T) {
	t.Parallel()

	service := NewService(&fakeRepository{})
	_, err := service.CreateInvite(context.Background(), "owner-1", "article-1", InviteInput{
		IdentityMode: IdentityModeAnonymous,
	})
	if err != nil {
		t.Fatalf("expected anonymous invite without email to succeed, got %v", err)
	}
}

func TestCreateInviteRejectsAnonymousInviteWithEmail(t *testing.T) {
	t.Parallel()

	service := NewService(&fakeRepository{})
	_, err := service.CreateInvite(context.Background(), "owner-1", "article-1", InviteInput{
		IdentityMode: IdentityModeAnonymous,
		Email:        "reader@example.com",
	})
	if err == nil {
		t.Fatalf("expected anonymous invite with email to fail")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}

func TestDisplayNameFromEmail(t *testing.T) {
	t.Parallel()

	if got := DisplayNameFromEmail("sarah_k@example.com"); got != "Sarah K." {
		t.Fatalf("expected Sarah K., got %q", got)
	}
	if got := AvatarFromName("Sarah K."); got != "SK" {
		t.Fatalf("expected SK avatar, got %q", got)
	}
}

func TestDisplayNameFromInvite(t *testing.T) {
	t.Parallel()

	if got := DisplayNameFromInvite("Workshop Reader", "", IdentityModeNamed); got != "Workshop Reader" {
		t.Fatalf("expected explicit display name, got %q", got)
	}
	if got := DisplayNameFromInvite("", "", IdentityModeAnonymous); got != "Anonymous Reader" {
		t.Fatalf("expected anonymous fallback, got %q", got)
	}
}
