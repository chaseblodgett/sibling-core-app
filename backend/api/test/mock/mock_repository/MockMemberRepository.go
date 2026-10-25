package mockrepository

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockMemberRepository is a test double for repository.MemberRepository.
// Assign the *Func fields your test exercises; unset ones return zero values.
type MockMemberRepository struct {
	CreateFunc             func(ctx context.Context, member *models.Member) error
	GetByIDFunc            func(ctx context.Context, id string) (*models.Member, error)
	GetByUserAndFamilyFunc func(ctx context.Context, userID, familyID string) (*models.Member, error)
	ListByFamilyFunc       func(ctx context.Context, familyID string) ([]*models.Member, error)
	ListByUserFunc         func(ctx context.Context, userID string) ([]*models.Member, error)
	DeleteFunc             func(ctx context.Context, id string) error
}

func (m *MockMemberRepository) Create(ctx context.Context, member *models.Member) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, member)
}

func (m *MockMemberRepository) GetByID(ctx context.Context, id string) (*models.Member, error) {
	if m.GetByIDFunc == nil {
		return nil, nil
	}
	return m.GetByIDFunc(ctx, id)
}

func (m *MockMemberRepository) GetByUserAndFamily(ctx context.Context, userID, familyID string) (*models.Member, error) {
	if m.GetByUserAndFamilyFunc == nil {
		return nil, nil
	}
	return m.GetByUserAndFamilyFunc(ctx, userID, familyID)
}

func (m *MockMemberRepository) ListByFamily(ctx context.Context, familyID string) ([]*models.Member, error) {
	if m.ListByFamilyFunc == nil {
		return nil, nil
	}
	return m.ListByFamilyFunc(ctx, familyID)
}

func (m *MockMemberRepository) ListByUser(ctx context.Context, userID string) ([]*models.Member, error) {
	if m.ListByUserFunc == nil {
		return nil, nil
	}
	return m.ListByUserFunc(ctx, userID)
}

func (m *MockMemberRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc == nil {
		return nil
	}
	return m.DeleteFunc(ctx, id)
}
