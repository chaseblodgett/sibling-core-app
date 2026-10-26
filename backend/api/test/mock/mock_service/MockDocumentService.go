package mockservice

import (
	"context"
	"io"

	"sibling-core/api/internal/models"
)

// MockDocumentService is a test double for service.DocumentService. Assign
// the *Func fields your test exercises; unset ones return zero values.
type MockDocumentService struct {
	UploadDocumentFunc   func(ctx context.Context, input models.UploadDocumentInput) (*models.Document, error)
	GetDocumentFunc      func(ctx context.Context, id string) (*models.Document, error)
	ListDocumentsFunc    func(ctx context.Context, familyID string, parentID *string) ([]*models.Document, error)
	DownloadDocumentFunc func(ctx context.Context, id string) (*models.Document, io.ReadCloser, error)
	DeleteDocumentFunc   func(ctx context.Context, id string) error
}

func (m *MockDocumentService) UploadDocument(ctx context.Context, input models.UploadDocumentInput) (*models.Document, error) {
	if m.UploadDocumentFunc == nil {
		return nil, nil
	}
	return m.UploadDocumentFunc(ctx, input)
}

func (m *MockDocumentService) GetDocument(ctx context.Context, id string) (*models.Document, error) {
	if m.GetDocumentFunc == nil {
		return nil, nil
	}
	return m.GetDocumentFunc(ctx, id)
}

func (m *MockDocumentService) ListDocuments(ctx context.Context, familyID string, parentID *string) ([]*models.Document, error) {
	if m.ListDocumentsFunc == nil {
		return nil, nil
	}
	return m.ListDocumentsFunc(ctx, familyID, parentID)
}

func (m *MockDocumentService) DownloadDocument(ctx context.Context, id string) (*models.Document, io.ReadCloser, error) {
	if m.DownloadDocumentFunc == nil {
		return nil, nil, nil
	}
	return m.DownloadDocumentFunc(ctx, id)
}

func (m *MockDocumentService) DeleteDocument(ctx context.Context, id string) error {
	if m.DeleteDocumentFunc == nil {
		return nil
	}
	return m.DeleteDocumentFunc(ctx, id)
}
