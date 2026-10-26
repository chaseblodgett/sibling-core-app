package mockservice

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockParentService is a test double for service.ParentService. Assign
// the *Func fields your test exercises; unset ones return zero values.
type MockParentService struct {
	CreateParentFunc func(ctx context.Context, input models.CreateParentInput) (*models.Parent, error)
	GetParentFunc    func(ctx context.Context, id string) (*models.Parent, error)
	ListParentsFunc  func(ctx context.Context, familyID string) ([]*models.Parent, error)
	UpdateParentFunc func(ctx context.Context, id string, input models.UpdateParentInput) (*models.Parent, error)
	DeleteParentFunc func(ctx context.Context, id string) error
}

func (m *MockParentService) CreateParent(ctx context.Context, input models.CreateParentInput) (*models.Parent, error) {
	if m.CreateParentFunc == nil {
		return nil, nil
	}
	return m.CreateParentFunc(ctx, input)
}

func (m *MockParentService) GetParent(ctx context.Context, id string) (*models.Parent, error) {
	if m.GetParentFunc == nil {
		return nil, nil
	}
	return m.GetParentFunc(ctx, id)
}

func (m *MockParentService) ListParents(ctx context.Context, familyID string) ([]*models.Parent, error) {
	if m.ListParentsFunc == nil {
		return nil, nil
	}
	return m.ListParentsFunc(ctx, familyID)
}

func (m *MockParentService) UpdateParent(ctx context.Context, id string, input models.UpdateParentInput) (*models.Parent, error) {
	if m.UpdateParentFunc == nil {
		return nil, nil
	}
	return m.UpdateParentFunc(ctx, id, input)
}

func (m *MockParentService) DeleteParent(ctx context.Context, id string) error {
	if m.DeleteParentFunc == nil {
		return nil
	}
	return m.DeleteParentFunc(ctx, id)
}
