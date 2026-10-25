package mockservice

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockDocumentService is a test double for service.DocumentService. Assign
// the *Func fields your test exercises; unset ones return zero values.
type MockDocumentService struct {
	UploadFunc func(ctx context.Context, input models.UploadDocumentInput) (*models.Document, error)
	GetFunc    func(ctx context.Context, id string) (*models.Document, error)
	ListFunc   func(ctx context.Context, familyID string) ([]*models.Document, error)
	DeleteFunc func(ctx context.Context, id string) error
}

func (m *MockDocumentService) Upload(ctx context.Context, input models.UploadDocumentInput) (*models.Document, error) {
	if m.UploadFunc == nil {
		return nil, nil
	}
	return m.UploadFunc(ctx, input)
}

func (m *MockDocumentService) Get(ctx context.Context, id string) (*models.Document, error) {
	if m.GetFunc == nil {
		return nil, nil
	}
	return m.GetFunc(ctx, id)
}

func (m *MockDocumentService) List(ctx context.Context, familyID string) ([]*models.Document, error) {
	if m.ListFunc == nil {
		return nil, nil
	}
	return m.ListFunc(ctx, familyID)
}

func (m *MockDocumentService) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc == nil {
		return nil
	}
	return m.DeleteFunc(ctx, id)
}
