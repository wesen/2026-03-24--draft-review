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

	normalized, err := NormalizeInviteInput(input)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateInvite(ctx, ownerUserID, articleID, normalized)
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

func NormalizeInviteInput(input InviteInput) (InviteInput, error) {
	input.IdentityMode = strings.ToLower(strings.TrimSpace(input.IdentityMode))
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Note = strings.TrimSpace(input.Note)

	if input.IdentityMode == "" {
		switch {
		case input.IsPreview:
			input.IdentityMode = IdentityModePreview
		case input.Email != "":
			input.IdentityMode = IdentityModeEmail
		case input.DisplayName != "":
			input.IdentityMode = IdentityModeNamed
		default:
			input.IdentityMode = IdentityModeAnonymous
		}
	}

	switch input.IdentityMode {
	case IdentityModeEmail:
		if input.Email == "" {
			return input, NewValidationError("email is required for email invite links")
		}
		if _, err := mail.ParseAddress(input.Email); err != nil {
			return input, NewValidationError("email must be a valid email address")
		}
	case IdentityModeNamed:
		if input.DisplayName == "" {
			return input, NewValidationError("displayName is required for named invite links")
		}
		if input.Email != "" {
			if _, err := mail.ParseAddress(input.Email); err != nil {
				return input, NewValidationError("email must be a valid email address")
			}
		}
	case IdentityModeAnonymous:
		if input.Email != "" || input.DisplayName != "" {
			return input, NewValidationError("anonymous invite links cannot include email or displayName")
		}
	case IdentityModePreview:
		if input.Email != "" {
			if _, err := mail.ParseAddress(input.Email); err != nil {
				return input, NewValidationError("email must be a valid email address")
			}
		}
		if input.DisplayName == "" {
			input.DisplayName = "Preview Reader"
		}
		input.IsPreview = true
	default:
		return input, NewValidationError("identityMode must be one of email, named, anonymous, preview")
	}

	return input, nil
}

func DisplayNameFromInvite(displayName, email, identityMode string) string {
	displayName = strings.TrimSpace(displayName)
	if displayName != "" {
		return displayName
	}

	email = strings.TrimSpace(email)
	if email != "" {
		return DisplayNameFromEmail(email)
	}

	switch strings.TrimSpace(identityMode) {
	case IdentityModeAnonymous:
		return "Anonymous Reader"
	case IdentityModePreview:
		return "Preview Reader"
	default:
		return "Reader"
	}
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
