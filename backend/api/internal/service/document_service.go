package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/storage"
)

var unsafeFileNameChars = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// DocumentService contains business logic for uploading, downloading, and
// listing family documents. File bytes live in Supabase Storage; this
// coordinates that with the metadata row persisted via DocumentRepository.
type DocumentService interface {
	UploadDocument(ctx context.Context, input models.UploadDocumentInput) (*models.Document, error)
	GetDocument(ctx context.Context, id string) (*models.Document, error)
	// ListDocuments returns a family's documents, optionally filtered to a parent.
	ListDocuments(ctx context.Context, familyID string, parentID *string) ([]*models.Document, error)
	// DownloadDocument returns the document's metadata and its content;
	// the caller must Close the content.
	DownloadDocument(ctx context.Context, id string) (*models.Document, io.ReadCloser, error)
	DeleteDocument(ctx context.Context, id string) error
}

type documentService struct {
	documentRepo repository.DocumentRepository
	parentRepo   repository.ParentRepository
	storage      storage.Storage
	activityRepo repository.ActivityRepository
}

// NewDocumentService constructs a DocumentService.
func NewDocumentService(documentRepo repository.DocumentRepository, parentRepo repository.ParentRepository, store storage.Storage, activityRepo repository.ActivityRepository) DocumentService {
	return &documentService{documentRepo: documentRepo, parentRepo: parentRepo, storage: store, activityRepo: activityRepo}
}

func (s *documentService) UploadDocument(ctx context.Context, input models.UploadDocumentInput) (*models.Document, error) {
	fileName := strings.TrimSpace(input.FileName)
	if fileName == "" {
		return nil, fmt.Errorf("%w: file name is required", ErrInvalidInput)
	}
	if input.SizeBytes <= 0 {
		return nil, fmt.Errorf("%w: file is empty", ErrInvalidInput)
	}
	if input.ParentID != nil {
		if err := requireParentInFamily(ctx, s.parentRepo, *input.ParentID, input.FamilyID); err != nil {
			return nil, err
		}
	}

	contentType := input.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	token, err := randomToken()
	if err != nil {
		return nil, err
	}
	storagePath := fmt.Sprintf("%s/%s-%s", input.FamilyID, token, sanitizeFileName(fileName))

	if err := s.storage.Upload(ctx, storagePath, contentType, input.Content); err != nil {
		return nil, err
	}

	doc := &models.Document{
		FamilyID:    input.FamilyID,
		UploadedBy:  input.UploadedBy,
		FileName:    fileName,
		StoragePath: storagePath,
		ContentType: contentType,
		SizeBytes:   input.SizeBytes,
		ParentID:    input.ParentID,
	}
	if err := s.documentRepo.Create(ctx, doc); err != nil {
		// The DB row is the source of truth; a dangling storage object
		// from a failed metadata insert is harmless clutter, so clean it
		// up best-effort rather than leave the upload "half done" either way.
		if cleanupErr := s.storage.Delete(ctx, storagePath); cleanupErr != nil {
			log.Printf("document upload: storage cleanup failed for %s: %v", storagePath, cleanupErr)
		}
		return nil, err
	}
	logActivity(ctx, s.activityRepo, doc.FamilyID, &doc.UploadedBy, input.ActorName, "document.uploaded", "document", &doc.ID, doc.FileName)
	return doc, nil
}

func (s *documentService) GetDocument(ctx context.Context, id string) (*models.Document, error) {
	return s.documentRepo.GetByID(ctx, id)
}

func (s *documentService) ListDocuments(ctx context.Context, familyID string, parentID *string) ([]*models.Document, error) {
	return s.documentRepo.ListByFamily(ctx, familyID, parentID)
}

func (s *documentService) DownloadDocument(ctx context.Context, id string) (*models.Document, io.ReadCloser, error) {
	doc, err := s.documentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	content, err := s.storage.Download(ctx, doc.StoragePath)
	if err != nil {
		return nil, nil, err
	}
	return doc, content, nil
}

func (s *documentService) DeleteDocument(ctx context.Context, id string) error {
	doc, err := s.documentRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.documentRepo.Delete(ctx, id); err != nil {
		return err
	}
	// The metadata row is already gone — the document "doesn't exist" to
	// the app from here regardless of whether this cleanup succeeds.
	if err := s.storage.Delete(ctx, doc.StoragePath); err != nil {
		log.Printf("document delete: storage cleanup failed for %s: %v", doc.StoragePath, err)
	}
	return nil
}

func sanitizeFileName(name string) string {
	return unsafeFileNameChars.ReplaceAllString(name, "_")
}

func randomToken() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
