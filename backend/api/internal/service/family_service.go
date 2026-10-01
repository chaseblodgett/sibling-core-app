package service

import (
	"context"
	"errors"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// FamilyService contains business logic for family groups, membership, and invites.
type FamilyService interface {
	CreateFamily(ctx context.Context, input models.CreateFamilyInput) (*models.Family, error)
	GetFamily(ctx context.Context, id string) (*models.Family, error)
	JoinFamily(ctx context.Context, inviteCode, userID string) (*models.Family, error)
	InviteMember(ctx context.Context, input models.InviteMemberInput) (*models.FamilyInvite, error)
	ListMembers(ctx context.Context, familyID string) ([]*models.Member, error)
	RemoveMember(ctx context.Context, familyID, memberID string) error
}

type familyService struct {
	familyRepo repository.FamilyRepository
	memberRepo repository.MemberRepository
}

// NewFamilyService constructs a FamilyService.
func NewFamilyService(familyRepo repository.FamilyRepository, memberRepo repository.MemberRepository) FamilyService {
	return &familyService{familyRepo: familyRepo, memberRepo: memberRepo}
}

func (s *familyService) CreateFamily(ctx context.Context, input models.CreateFamilyInput) (*models.Family, error) {
	return nil, errors.New("not implemented")
}

func (s *familyService) GetFamily(ctx context.Context, id string) (*models.Family, error) {
	return nil, errors.New("not implemented")
}

func (s *familyService) JoinFamily(ctx context.Context, inviteCode, userID string) (*models.Family, error) {
	return nil, errors.New("not implemented")
}

func (s *familyService) InviteMember(ctx context.Context, input models.InviteMemberInput) (*models.FamilyInvite, error) {
	return nil, errors.New("not implemented")
}

func (s *familyService) ListMembers(ctx context.Context, familyID string) ([]*models.Member, error) {
	return nil, errors.New("not implemented")
}

func (s *familyService) RemoveMember(ctx context.Context, familyID, memberID string) error {
	return errors.New("not implemented")
}
