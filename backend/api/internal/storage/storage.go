// Package storage stores and retrieves document file bytes in Supabase
// Storage over its plain HTTP object API — no SDK dependency, same
// approach as internal/mailer for Resend. The bucket is private; only
// this backend (using the service role key) ever talks to it directly,
// so there's no storage-level RLS to reason about.
package storage

import (
	"context"
	"errors"
	"io"
)

// ErrObjectNotFound means the object doesn't exist in storage.
var ErrObjectNotFound = errors.New("storage object not found")

// Storage stores and retrieves file bytes for documents. path identifies
// an object within the storage backend's single "documents" bucket.
type Storage interface {
	Upload(ctx context.Context, path, contentType string, content io.Reader) error
	// Download returns the object's bytes; the caller must Close it.
	Download(ctx context.Context, path string) (io.ReadCloser, error)
	// Delete is a no-op (not an error) if the object doesn't exist.
	Delete(ctx context.Context, path string) error
}
