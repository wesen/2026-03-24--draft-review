package articleassets

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	draftauth "github.com/go-go-golems/draft-review/pkg/auth"
)

type fakeRepository struct {
	record *AssetRecord
}

func (f *fakeRepository) CreateAssetRecord(_ context.Context, ownerUserID, articleID string, input AssetRecordInput) (*AssetRecord, error) {
	f.record = &AssetRecord{
		ID:               input.ID,
		ArticleID:        articleID,
		StorageKey:       input.StorageKey,
		OriginalFilename: input.OriginalFilename,
		ContentType:      input.ContentType,
		ByteSize:         input.ByteSize,
		CreatedAt:        time.Now().UTC(),
	}
	return f.record, nil
}

func (f *fakeRepository) GetAssetByID(_ context.Context, assetID string) (*AssetRecord, error) {
	if f.record == nil || f.record.ID != assetID {
		return nil, ErrNotFound
	}
	return f.record, nil
}

func TestUploadImageReturnsMarkdownSnippet(t *testing.T) {
	t.Parallel()

	repo := &fakeRepository{}
	storage := NewMemoryStorage()
	service := NewService(repo, storage, 1024*1024, "/media/article-assets")

	pngBytes, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+a4i8AAAAASUVORK5CYII=")
	if err != nil {
		t.Fatalf("failed to decode png fixture: %v", err)
	}

	asset, err := service.UploadImage(context.Background(), &draftauth.User{ID: "user-1"}, "article-1", UploadInput{
		Filename: "diagram.png",
		Content:  pngBytes,
	})
	if err != nil {
		t.Fatalf("UploadImage returned error: %v", err)
	}

	if asset.URL == "" {
		t.Fatalf("expected URL to be populated")
	}
	if asset.Markdown != "![diagram]("+asset.URL+")" {
		t.Fatalf("unexpected markdown snippet: %q", asset.Markdown)
	}

	reader, err := storage.Open(context.Background(), repo.record.StorageKey)
	if err != nil {
		t.Fatalf("expected stored content to exist: %v", err)
	}
	_ = reader.Close()
}

func TestUploadImageRejectsUnsupportedMime(t *testing.T) {
	t.Parallel()

	service := NewService(&fakeRepository{}, NewMemoryStorage(), 1024, "/media/article-assets")
	_, err := service.UploadImage(context.Background(), &draftauth.User{ID: "user-1"}, "article-1", UploadInput{
		Filename: "vector.svg",
		Content:  []byte("<svg xmlns=\"http://www.w3.org/2000/svg\"></svg>"),
	})
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !IsValidationError(err) {
		t.Fatalf("expected validation error, got %T", err)
	}
}
