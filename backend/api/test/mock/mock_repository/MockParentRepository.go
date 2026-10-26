package mockrepository

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockParentRepository is a test double for repository.ParentRepository.
// Assign the *Func fields your test exercises; unset ones return zero values.
type MockParentRepository struct {
	CreateFunc       func(ctx context.Context, parent *models.Parent) error
	GetByIDFunc      func(ctx context.Context, id string) (*models.Parent, error)
	ListByFamilyFunc func(ctx context.Context, familyID string) ([]*models.Parent, error)
	UpdateFunc       func(ctx context.Context, parent *models.Parent) error
	DeleteFunc       func(ctx context.Context, id string) error
}

func (m *MockParentRepository) Create(ctx context.Context, parent *models.Parent) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, parent)
}

func (m *MockParentRepository) GetByID(ctx context.Context, id string) (*models.Parent, error) {
	if m.GetByIDFunc == nil {
		return nil, nil
	}
	return m.GetByIDFunc(ctx, id)
}

func (m *MockParentRepository) ListByFamily(ctx context.Context, familyID string) ([]*models.Parent, error) {
	if m.ListByFamilyFunc == nil {
		return nil, nil
	}
	return m.ListByFamilyFunc(ctx, familyID)
}

func (m *MockParentRepository) Update(ctx context.Context, parent *models.Parent) error {
	if m.UpdateFunc == nil {
		return nil
	}
	return m.UpdateFunc(ctx, parent)
}

func (m *MockParentRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc == nil {
		return nil
	}
	return m.DeleteFunc(ctx, id)
}
