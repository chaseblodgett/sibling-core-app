package mockrepository

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockActivityRepository is a test double for repository.ActivityRepository.
// Assign the *Func fields your test exercises; unset ones return zero values.
type MockActivityRepository struct {
	CreateFunc       func(ctx context.Context, entry *models.ActivityEntry) error
	ListByFamilyFunc func(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error)
}

func (m *MockActivityRepository) Create(ctx context.Context, entry *models.ActivityEntry) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, entry)
}

func (m *MockActivityRepository) ListByFamily(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error) {
	if m.ListByFamilyFunc == nil {
		return nil, nil
	}
	return m.ListByFamilyFunc(ctx, familyID, limit, offset)
}
