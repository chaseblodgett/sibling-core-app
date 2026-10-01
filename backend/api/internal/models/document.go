package models

import (
	"io"
	"time"
)

// Document represents a file uploaded to Supabase Storage and tied to a family.
type Document struct {
	ID          string    `json:"id" db:"id"`
	FamilyID    string    `json:"family_id" db:"family_id"`
	UploadedBy  string    `json:"uploaded_by" db:"uploaded_by"`
	FileName    string    `json:"file_name" db:"file_name"`
	StoragePath string    `json:"storage_path" db:"storage_path"` // path within Supabase Storage bucket
	ContentType string    `json:"content_type" db:"content_type"`
	SizeBytes   int64     `json:"size_bytes" db:"size_bytes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UploadDocumentInput carries the fields needed to upload a document.
type UploadDocumentInput struct {
	FamilyID    string
	UploadedBy  string
	FileName    string
	ContentType string
	Content     io.Reader
}
