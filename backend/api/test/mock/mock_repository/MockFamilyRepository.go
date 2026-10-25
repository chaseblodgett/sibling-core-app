package mockrepository

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockFamilyRepository is a test double for repository.FamilyRepository.
// Assign the *Func fields your test exercises; unset ones return zero values.
type MockFamilyRepository struct {
	CreateFunc              func(ctx context.Context, family *models.Family) error
	GetByIDFunc             func(ctx context.Context, id string) (*models.Family, error)
	GetByInviteCodeFunc     func(ctx context.Context, code string) (*models.Family, error)
	UpdateFunc              func(ctx context.Context, family *models.Family) error
	DeleteFunc              func(ctx context.Context, id string) error
	CreateInviteFunc        func(ctx context.Context, invite *models.FamilyInvite) error
	GetInviteByIDFunc       func(ctx context.Context, id string) (*models.FamilyInvite, error)
	ListInvitesByFamilyFunc func(ctx context.Context, familyID string) ([]*models.FamilyInvite, error)
	UpdateInviteStatusFunc  func(ctx context.Context, id string, status string) error
}

func (m *MockFamilyRepository) Create(ctx context.Context, family *models.Family) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, family)
}

func (m *MockFamilyRepository) GetByID(ctx context.Context, id string) (*models.Family, error) {
	if m.GetByIDFunc == nil {
		return nil, nil
	}
	return m.GetByIDFunc(ctx, id)
}

func (m *MockFamilyRepository) GetByInviteCode(ctx context.Context, code string) (*models.Family, error) {
	if m.GetByInviteCodeFunc == nil {
		return nil, nil
	}
	return m.GetByInviteCodeFunc(ctx, code)
}

func (m *MockFamilyRepository) Update(ctx context.Context, family *models.Family) error {
	if m.UpdateFunc == nil {
		return nil
	}
	return m.UpdateFunc(ctx, family)
}

func (m *MockFamilyRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc == nil {
		return nil
	}
	return m.DeleteFunc(ctx, id)
}

func (m *MockFamilyRepository) CreateInvite(ctx context.Context, invite *models.FamilyInvite) error {
	if m.CreateInviteFunc == nil {
		return nil
	}
	return m.CreateInviteFunc(ctx, invite)
}

func (m *MockFamilyRepository) GetInviteByID(ctx context.Context, id string) (*models.FamilyInvite, error) {
	if m.GetInviteByIDFunc == nil {
		return nil, nil
	}
	return m.GetInviteByIDFunc(ctx, id)
}

func (m *MockFamilyRepository) ListInvitesByFamily(ctx context.Context, familyID string) ([]*models.FamilyInvite, error) {
	if m.ListInvitesByFamilyFunc == nil {
		return nil, nil
	}
	return m.ListInvitesByFamilyFunc(ctx, familyID)
}

func (m *MockFamilyRepository) UpdateInviteStatus(ctx context.Context, id string, status string) error {
	if m.UpdateInviteStatusFunc == nil {
		return nil
	}
	return m.UpdateInviteStatusFunc(ctx, id, status)
}
