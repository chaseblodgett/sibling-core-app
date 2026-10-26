package mockrepository

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockDocumentRepository is a test double for repository.DocumentRepository.
// Assign the *Func fields your test exercises; unset ones return zero values.
type MockDocumentRepository struct {
	CreateFunc       func(ctx context.Context, document *models.Document) error
	GetByIDFunc      func(ctx context.Context, id string) (*models.Document, error)
	ListByFamilyFunc func(ctx context.Context, familyID string, parentID *string) ([]*models.Document, error)
	DeleteFunc       func(ctx context.Context, id string) error
}

func (m *MockDocumentRepository) Create(ctx context.Context, document *models.Document) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, document)
}

func (m *MockDocumentRepository) GetByID(ctx context.Context, id string) (*models.Document, error) {
	if m.GetByIDFunc == nil {
		return nil, nil
	}
	return m.GetByIDFunc(ctx, id)
}

func (m *MockDocumentRepository) ListByFamily(ctx context.Context, familyID string, parentID *string) ([]*models.Document, error) {
	if m.ListByFamilyFunc == nil {
		return nil, nil
	}
	return m.ListByFamilyFunc(ctx, familyID, parentID)
}

func (m *MockDocumentRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc == nil {
		return nil
	}
	return m.DeleteFunc(ctx, id)
}
