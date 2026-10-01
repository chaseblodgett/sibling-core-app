package service

import (
	"context"
	"errors"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// DocumentService contains business logic for uploading and listing family documents.
// File bytes are stored in Supabase Storage; this service coordinates that
// upload with the metadata row persisted via DocumentRepository.
type DocumentService interface {
	Upload(ctx context.Context, input models.UploadDocumentInput) (*models.Document, error)
	Get(ctx context.Context, id string) (*models.Document, error)
	List(ctx context.Context, familyID string) ([]*models.Document, error)
	Delete(ctx context.Context, id string) error
}

type documentService struct {
	documentRepo repository.DocumentRepository
}

// NewDocumentService constructs a DocumentService.
func NewDocumentService(documentRepo repository.DocumentRepository) DocumentService {
	return &documentService{documentRepo: documentRepo}
}

func (s *documentService) Upload(ctx context.Context, input models.UploadDocumentInput) (*models.Document, error) {
	return nil, errors.New("not implemented")
}

func (s *documentService) Get(ctx context.Context, id string) (*models.Document, error) {
	return nil, errors.New("not implemented")
}

func (s *documentService) List(ctx context.Context, familyID string) ([]*models.Document, error) {
	return nil, errors.New("not implemented")
}

func (s *documentService) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
