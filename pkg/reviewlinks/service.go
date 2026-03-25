package reviewlinks

import (
	"context"
	"net/mail"
	"strings"
	"unicode"
)

type Repository interface {
	ResetShareToken(ctx context.Context, ownerUserID, articleID string) (*ShareLink, error)
	CreateInvite(ctx context.Context, ownerUserID, articleID string, input InviteInput) (*Reader, error)
	ResolveToken(ctx context.Context, token string) (*ResolvedLink, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ResetShareToken(ctx context.Context, ownerUserID, articleID string) (*ShareLink, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}

	ownerUserID = strings.TrimSpace(ownerUserID)
	articleID = strings.TrimSpace(articleID)
	if ownerUserID == "" || articleID == "" {
		return nil, ErrNotFound
	}

	return s.repo.ResetShareToken(ctx, ownerUserID, articleID)
}

func (s *Service) CreateInvite(ctx context.Context, ownerUserID, articleID string, input InviteInput) (*Reader, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}

	ownerUserID = strings.TrimSpace(ownerUserID)
	articleID = strings.TrimSpace(articleID)
	if ownerUserID == "" || articleID == "" {
		return nil, ErrNotFound
	}

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Note = strings.TrimSpace(input.Note)

	if _, err := mail.ParseAddress(input.Email); err != nil {
		return nil, NewValidationError("email must be a valid email address")
	}

	return s.repo.CreateInvite(ctx, ownerUserID, articleID, input)
}

func (s *Service) ResolveToken(ctx context.Context, token string) (*ResolvedLink, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, ErrNotFound
	}

	return s.repo.ResolveToken(ctx, token)
}

func DisplayNameFromEmail(email string) string {
	localPart := strings.TrimSpace(strings.Split(strings.ToLower(email), "@")[0])
	if localPart == "" {
		return "Reader"
	}

	replacer := strings.NewReplacer(".", " ", "_", " ", "-", " ", "+", " ")
	parts := strings.Fields(replacer.Replace(localPart))
	if len(parts) == 0 {
		return "Reader"
	}

	for i, part := range parts {
		runes := []rune(part)
		if len(runes) == 0 {
			continue
		}
		runes[0] = unicode.ToUpper(runes[0])
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}
		parts[i] = string(runes)
	}

	if len(parts) == 1 {
		return parts[0]
	}

	last := []rune(parts[len(parts)-1])
	if len(last) == 0 {
		return parts[0]
	}

	return parts[0] + " " + string(unicode.ToUpper(last[0])) + "."
}

func AvatarFromName(name string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) == 0 {
		return "RD"
	}

	if len(parts) == 1 {
		runes := []rune(parts[0])
		if len(runes) == 1 {
			return strings.ToUpper(string(runes[0])) + "R"
		}
		return strings.ToUpper(string(runes[0]) + string(runes[1]))
	}

	first := []rune(parts[0])
	last := []rune(parts[len(parts)-1])
	return strings.ToUpper(string(first[0]) + string(last[0]))
}
