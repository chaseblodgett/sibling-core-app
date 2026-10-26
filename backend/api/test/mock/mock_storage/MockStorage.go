package mockstorage

import (
	"context"
	"io"
)

// MockStorage is a test double for storage.Storage. Assign the *Func
// fields your test exercises; unset ones return zero values.
type MockStorage struct {
	UploadFunc   func(ctx context.Context, path, contentType string, content io.Reader) error
	DownloadFunc func(ctx context.Context, path string) (io.ReadCloser, error)
	DeleteFunc   func(ctx context.Context, path string) error
}

func (m *MockStorage) Upload(ctx context.Context, path, contentType string, content io.Reader) error {
	if m.UploadFunc == nil {
		return nil
	}
	return m.UploadFunc(ctx, path, contentType, content)
}

func (m *MockStorage) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	if m.DownloadFunc == nil {
		return nil, nil
	}
	return m.DownloadFunc(ctx, path)
}

func (m *MockStorage) Delete(ctx context.Context, path string) error {
	if m.DeleteFunc == nil {
		return nil
	}
	return m.DeleteFunc(ctx, path)
}
