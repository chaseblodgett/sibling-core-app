package mockservice

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockFamilyService is a test double for service.FamilyService. Assign the
// *Func fields your test exercises; unset ones return zero values.
type MockFamilyService struct {
	CreateFamilyFunc func(ctx context.Context, input models.CreateFamilyInput) (*models.Family, error)
	GetFamilyFunc    func(ctx context.Context, id string) (*models.Family, error)
	JoinFamilyFunc   func(ctx context.Context, inviteCode, userID string) (*models.Family, error)
	InviteMemberFunc func(ctx context.Context, input models.InviteMemberInput) (*models.FamilyInvite, error)
	ListMembersFunc  func(ctx context.Context, familyID string) ([]*models.Member, error)
	RemoveMemberFunc func(ctx context.Context, familyID, memberID string) error

	ListFamiliesForUserFunc func(ctx context.Context, userID string) ([]*models.Family, error)
	RequireMembershipFunc   func(ctx context.Context, userID, familyID string) (*models.Member, error)
}

func (m *MockFamilyService) CreateFamily(ctx context.Context, input models.CreateFamilyInput) (*models.Family, error) {
	if m.CreateFamilyFunc == nil {
		return nil, nil
	}
	return m.CreateFamilyFunc(ctx, input)
}

func (m *MockFamilyService) GetFamily(ctx context.Context, id string) (*models.Family, error) {
	if m.GetFamilyFunc == nil {
		return nil, nil
	}
	return m.GetFamilyFunc(ctx, id)
}

func (m *MockFamilyService) JoinFamily(ctx context.Context, inviteCode, userID string) (*models.Family, error) {
	if m.JoinFamilyFunc == nil {
		return nil, nil
	}
	return m.JoinFamilyFunc(ctx, inviteCode, userID)
}

func (m *MockFamilyService) InviteMember(ctx context.Context, input models.InviteMemberInput) (*models.FamilyInvite, error) {
	if m.InviteMemberFunc == nil {
		return nil, nil
	}
	return m.InviteMemberFunc(ctx, input)
}

func (m *MockFamilyService) ListMembers(ctx context.Context, familyID string) ([]*models.Member, error) {
	if m.ListMembersFunc == nil {
		return nil, nil
	}
	return m.ListMembersFunc(ctx, familyID)
}

func (m *MockFamilyService) RemoveMember(ctx context.Context, familyID, memberID string) error {
	if m.RemoveMemberFunc == nil {
		return nil
	}
	return m.RemoveMemberFunc(ctx, familyID, memberID)
}

func (m *MockFamilyService) ListFamiliesForUser(ctx context.Context, userID string) ([]*models.Family, error) {
	if m.ListFamiliesForUserFunc == nil {
		return nil, nil
	}
	return m.ListFamiliesForUserFunc(ctx, userID)
}

func (m *MockFamilyService) RequireMembership(ctx context.Context, userID, familyID string) (*models.Member, error) {
	if m.RequireMembershipFunc == nil {
		return nil, nil
	}
	return m.RequireMembershipFunc(ctx, userID, familyID)
}
