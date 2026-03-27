package articleassets

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
)

type Storage interface {
	Save(ctx context.Context, key string, content []byte) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
}

type LocalDiskStorage struct {
	root string
}

func NewLocalDiskStorage(root string) *LocalDiskStorage {
	return &LocalDiskStorage{root: root}
}

func (s *LocalDiskStorage) Save(_ context.Context, key string, content []byte) error {
	path := filepath.Join(s.root, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, content, 0o644)
}

func (s *LocalDiskStorage) Open(_ context.Context, key string) (io.ReadCloser, error) {
	path := filepath.Join(s.root, filepath.FromSlash(key))
	return os.Open(path)
}

type MemoryStorage struct {
	files map[string][]byte
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{files: map[string][]byte{}}
}

func (s *MemoryStorage) Save(_ context.Context, key string, content []byte) error {
	s.files[key] = append([]byte(nil), content...)
	return nil
}

func (s *MemoryStorage) Open(_ context.Context, key string) (io.ReadCloser, error) {
	content, ok := s.files[key]
	if !ok {
		return nil, os.ErrNotExist
	}

	return io.NopCloser(bytes.NewReader(content)), nil
}
