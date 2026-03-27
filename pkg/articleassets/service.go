package articleassets

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	draftauth "github.com/go-go-golems/draft-review/pkg/auth"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	ErrNotFound         = errors.New("article asset not found")
	ErrStorageNotReady  = errors.New("article asset storage is not configured")
	ErrValidation       = errors.New("article asset validation failed")
	safeFilenamePattern = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	allowedContentTypes = map[string]string{
		"image/png":  ".png",
		"image/jpeg": ".jpg",
		"image/gif":  ".gif",
		"image/webp": ".webp",
	}
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

type Repository interface {
	CreateAssetRecord(ctx context.Context, ownerUserID, articleID string, input AssetRecordInput) (*AssetRecord, error)
	GetAssetByID(ctx context.Context, assetID string) (*AssetRecord, error)
}

type Service struct {
	repo           Repository
	storage        Storage
	maxUploadBytes int64
	publicBasePath string
}

func NewService(repo Repository, storage Storage, maxUploadBytes int64, publicBasePath string) *Service {
	publicBasePath = strings.TrimSpace(publicBasePath)
	if publicBasePath == "" {
		publicBasePath = "/media/article-assets"
	}

	return &Service{
		repo:           repo,
		storage:        storage,
		maxUploadBytes: maxUploadBytes,
		publicBasePath: strings.TrimRight(publicBasePath, "/"),
	}
}

func (s *Service) UploadImage(ctx context.Context, owner *draftauth.User, articleID string, input UploadInput) (*Asset, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}
	if s.storage == nil {
		return nil, ErrStorageNotReady
	}
	if owner == nil || strings.TrimSpace(owner.ID) == "" {
		return nil, ErrNotFound
	}

	filename := sanitizeFilename(input.Filename)
	if filename == "" {
		return nil, NewValidationError("uploaded file must have a valid filename")
	}

	content := input.Content
	if len(content) == 0 {
		return nil, NewValidationError("uploaded file is empty")
	}
	if s.maxUploadBytes > 0 && int64(len(content)) > s.maxUploadBytes {
		return nil, NewValidationError(fmt.Sprintf("uploaded file exceeds %d bytes", s.maxUploadBytes))
	}

	contentType := http.DetectContentType(content)
	extension, ok := allowedContentTypes[contentType]
	if !ok {
		return nil, NewValidationError("only png, jpeg, gif, and webp images are supported")
	}

	assetID := uuid.NewString()
	storageKey := fmt.Sprintf("article-assets/%s/%s%s", strings.TrimSpace(articleID), assetID, extension)
	if err := s.storage.Save(ctx, storageKey, content); err != nil {
		return nil, errors.Wrap(err, "failed to persist article asset content")
	}

	record, err := s.repo.CreateAssetRecord(ctx, owner.ID, articleID, AssetRecordInput{
		ID:               assetID,
		StorageKey:       storageKey,
		OriginalFilename: filename,
		ContentType:      contentType,
		ByteSize:         int64(len(content)),
	})
	if err != nil {
		return nil, err
	}

	return s.assetResponse(record), nil
}

func (s *Service) ResolveAsset(ctx context.Context, assetID string) (*Asset, error) {
	if s == nil || s.repo == nil {
		return nil, ErrNotFound
	}

	record, err := s.repo.GetAssetByID(ctx, strings.TrimSpace(assetID))
	if err != nil {
		return nil, err
	}

	return s.assetResponse(record), nil
}

func (s *Service) OpenAsset(ctx context.Context, assetID string) (*Asset, io.ReadCloser, error) {
	if s == nil || s.repo == nil {
		return nil, nil, ErrNotFound
	}
	if s.storage == nil {
		return nil, nil, ErrStorageNotReady
	}

	record, err := s.repo.GetAssetByID(ctx, strings.TrimSpace(assetID))
	if err != nil {
		return nil, nil, err
	}

	reader, err := s.storage.Open(ctx, record.StorageKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to open article asset content")
	}

	return s.assetResponse(record), reader, nil
}

func (s *Service) assetResponse(record *AssetRecord) *Asset {
	urlPath := fmt.Sprintf("%s/%s/%s", s.publicBasePath, record.ID, url.PathEscape(record.OriginalFilename))
	return &Asset{
		ID:               record.ID,
		ArticleID:        record.ArticleID,
		URL:              urlPath,
		Markdown:         fmt.Sprintf("![%s](%s)", altTextFromFilename(record.OriginalFilename), urlPath),
		OriginalFilename: record.OriginalFilename,
		ContentType:      record.ContentType,
		ByteSize:         record.ByteSize,
		CreatedAt:        record.CreatedAt,
	}
}

func sanitizeFilename(filename string) string {
	filename = filepath.Base(strings.TrimSpace(filename))
	if filename == "." || filename == "" {
		return ""
	}

	filename = safeFilenamePattern.ReplaceAllString(filename, "-")
	filename = strings.Trim(filename, "-.")
	if filename == "" {
		return ""
	}

	return filename
}

func altTextFromFilename(filename string) string {
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	base = strings.Join(strings.Fields(base), " ")
	if base == "" {
		return "Uploaded image"
	}
	return base
}
