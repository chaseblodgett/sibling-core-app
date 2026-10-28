package models

import (
	"io"
	"time"
)

// Document represents a file uploaded to Supabase Storage and tied to a
// family. StoragePath is an internal pointer into the storage bucket —
// never exposed to the client, which fetches content via the download
// route instead of talking to storage directly.
type Document struct {
	ID          string    `json:"id" db:"id"`
	FamilyID    string    `json:"family_id" db:"family_id"`
	UploadedBy  string    `json:"uploaded_by" db:"uploaded_by"`
	FileName    string    `json:"file_name" db:"file_name"`
	StoragePath string    `json:"-" db:"storage_path"`
	ContentType string    `json:"content_type" db:"content_type"`
	SizeBytes   int64     `json:"size_bytes" db:"size_bytes"`
	ParentID    *string   `json:"parent_id,omitempty" db:"parent_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UploadDocumentInput carries the fields needed to upload a document.
// UploadedBy is a member id (documents.uploaded_by references members,
// not users) — the handler must resolve the requester's member row for
// the family, same as Task/Expense.
type UploadDocumentInput struct {
	FamilyID    string
	UploadedBy  string
	FileName    string
	ContentType string
	SizeBytes   int64
	ParentID    *string
	Content     io.Reader
	// ActorName is UploadedBy's display name, for the activity-feed entry
	// this upload produces.
	ActorName string
}
