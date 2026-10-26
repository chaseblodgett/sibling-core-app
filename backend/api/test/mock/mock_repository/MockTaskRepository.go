package mockrepository

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockTaskRepository is a test double for repository.TaskRepository. Assign
// the *Func fields your test exercises; unset ones return zero values.
type MockTaskRepository struct {
	CreateFunc       func(ctx context.Context, task *models.Task) error
	GetByIDFunc      func(ctx context.Context, id string) (*models.Task, error)
	ListByFamilyFunc func(ctx context.Context, familyID string, parentID *string) ([]*models.Task, error)
	UpdateFunc       func(ctx context.Context, task *models.Task) error
	DeleteFunc       func(ctx context.Context, id string) error
}

func (m *MockTaskRepository) Create(ctx context.Context, task *models.Task) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, task)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	if m.GetByIDFunc == nil {
		return nil, nil
	}
	return m.GetByIDFunc(ctx, id)
}

func (m *MockTaskRepository) ListByFamily(ctx context.Context, familyID string, parentID *string) ([]*models.Task, error) {
	if m.ListByFamilyFunc == nil {
		return nil, nil
	}
	return m.ListByFamilyFunc(ctx, familyID, parentID)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *models.Task) error {
	if m.UpdateFunc == nil {
		return nil
	}
	return m.UpdateFunc(ctx, task)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc == nil {
		return nil
	}
	return m.DeleteFunc(ctx, id)
}
