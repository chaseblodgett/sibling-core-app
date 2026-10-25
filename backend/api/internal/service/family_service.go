package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// inviteCodeAlphabet excludes visually ambiguous characters (0/O, 1/I) since
// invite codes are meant to be read and typed by hand.
const inviteCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const inviteCodeLength = 8

// maxInviteCodeAttempts bounds retries if a generated code collides with an
// existing one; at 8 random characters from a 33-char alphabet collisions
// are vanishingly rare, so this only guards against pathological bad luck.
const maxInviteCodeAttempts = 5

var (
	ErrNotAuthorized     = errors.New("not authorized to perform this action")
	ErrNotMember         = errors.New("not a member of this family")
	ErrInvalidInviteCode = errors.New("invalid invite code")
)

// FamilyService contains business logic for family groups, membership, and invites.
type FamilyService interface {
	CreateFamily(ctx context.Context, input models.CreateFamilyInput) (*models.Family, error)
	GetFamily(ctx context.Context, id string) (*models.Family, error)
	JoinFamily(ctx context.Context, inviteCode, userID string) (*models.Family, error)
	InviteMember(ctx context.Context, input models.InviteMemberInput) (*models.FamilyInvite, error)
	ListMembers(ctx context.Context, familyID string) ([]*models.Member, error)
	RemoveMember(ctx context.Context, familyID, memberID string) error

	// ListFamiliesForUser returns every family a user belongs to, letting a
	// logged-in user discover "my family" without already knowing its id.
	ListFamiliesForUser(ctx context.Context, userID string) ([]*models.Family, error)
	// RequireMembership returns the requester's membership row for a family,
	// or ErrNotMember if they don't belong to it. Handlers use this to
	// authorize per-family actions (viewing, inviting, removing members).
	RequireMembership(ctx context.Context, userID, familyID string) (*models.Member, error)
}

type familyService struct {
	familyRepo repository.FamilyRepository
	memberRepo repository.MemberRepository
	userRepo   repository.UserRepository
}

// NewFamilyService constructs a FamilyService.
func NewFamilyService(familyRepo repository.FamilyRepository, memberRepo repository.MemberRepository, userRepo repository.UserRepository) FamilyService {
	return &familyService{familyRepo: familyRepo, memberRepo: memberRepo, userRepo: userRepo}
}

func (s *familyService) CreateFamily(ctx context.Context, input models.CreateFamilyInput) (*models.Family, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalidInput)
	}

	creator, err := s.userRepo.FindByID(ctx, input.CreatedBy)
	if err != nil {
		return nil, err
	}

	var family *models.Family
	for attempt := 0; attempt < maxInviteCodeAttempts; attempt++ {
		code, err := generateInviteCode()
		if err != nil {
			return nil, err
		}

		family = &models.Family{Name: name, InviteCode: code, CreatedBy: input.CreatedBy}
		err = s.familyRepo.Create(ctx, family)
		if err == nil {
			break
		}
		if errors.Is(err, repository.ErrDuplicateInviteCode) && attempt < maxInviteCodeAttempts-1 {
			continue
		}
		return nil, err
	}

	owner := &models.Member{
		FamilyID: family.ID,
		UserID:   creator.ID,
		Email:    creator.Email,
		FullName: creator.Name,
		Role:     "owner",
	}
	if err := s.memberRepo.Create(ctx, owner); err != nil {
		return nil, err
	}

	return family, nil
}

func (s *familyService) GetFamily(ctx context.Context, id string) (*models.Family, error) {
	return s.familyRepo.GetByID(ctx, id)
}

func (s *familyService) JoinFamily(ctx context.Context, inviteCode, userID string) (*models.Family, error) {
	code := strings.ToUpper(strings.TrimSpace(inviteCode))
	family, err := s.familyRepo.GetByInviteCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrFamilyNotFound) {
			return nil, ErrInvalidInviteCode
		}
		return nil, err
	}

	existing, err := s.memberRepo.GetByUserAndFamily(ctx, userID, family.ID)
	if err != nil && !errors.Is(err, repository.ErrMemberNotFound) {
		return nil, err
	}
	if existing != nil {
		return family, nil
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	member := &models.Member{
		FamilyID: family.ID,
		UserID:   user.ID,
		Email:    user.Email,
		FullName: user.Name,
		Role:     "member",
	}
	if err := s.memberRepo.Create(ctx, member); err != nil {
		return nil, err
	}

	s.acceptMatchingInvite(ctx, family.ID, user.Email)

	return family, nil
}

// acceptMatchingInvite marks the newest pending invite for this email as
// accepted, if one exists. Best-effort bookkeeping so the invite list
// reflects reality — joining succeeds regardless of its outcome.
func (s *familyService) acceptMatchingInvite(ctx context.Context, familyID, email string) {
	invites, err := s.familyRepo.ListInvitesByFamily(ctx, familyID)
	if err != nil {
		return
	}
	for _, invite := range invites {
		if invite.Status == "pending" && strings.EqualFold(invite.Email, email) {
			_ = s.familyRepo.UpdateInviteStatus(ctx, invite.ID, "accepted")
			return
		}
	}
}

func (s *familyService) InviteMember(ctx context.Context, input models.InviteMemberInput) (*models.FamilyInvite, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if !strings.Contains(email, "@") {
		return nil, fmt.Errorf("%w: invalid email", ErrInvalidInput)
	}

	invite := &models.FamilyInvite{
		FamilyID:  input.FamilyID,
		Email:     email,
		InvitedBy: input.InvitedBy,
		Status:    "pending",
	}
	if err := s.familyRepo.CreateInvite(ctx, invite); err != nil {
		return nil, err
	}
	return invite, nil
}

func (s *familyService) ListMembers(ctx context.Context, familyID string) ([]*models.Member, error) {
	return s.memberRepo.ListByFamily(ctx, familyID)
}

func (s *familyService) RemoveMember(ctx context.Context, familyID, memberID string) error {
	member, err := s.memberRepo.GetByID(ctx, memberID)
	if err != nil {
		return err
	}
	if member.FamilyID != familyID {
		return repository.ErrMemberNotFound
	}
	return s.memberRepo.Delete(ctx, memberID)
}

func (s *familyService) ListFamiliesForUser(ctx context.Context, userID string) ([]*models.Family, error) {
	memberships, err := s.memberRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	families := make([]*models.Family, 0, len(memberships))
	for _, m := range memberships {
		family, err := s.familyRepo.GetByID(ctx, m.FamilyID)
		if err != nil {
			return nil, err
		}
		families = append(families, family)
	}
	return families, nil
}

func (s *familyService) RequireMembership(ctx context.Context, userID, familyID string) (*models.Member, error) {
	member, err := s.memberRepo.GetByUserAndFamily(ctx, userID, familyID)
	if err != nil {
		if errors.Is(err, repository.ErrMemberNotFound) {
			return nil, ErrNotMember
		}
		return nil, err
	}
	return member, nil
}

func generateInviteCode() (string, error) {
	buf := make([]byte, inviteCodeLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	code := make([]byte, inviteCodeLength)
	for i, b := range buf {
		code[i] = inviteCodeAlphabet[int(b)%len(inviteCodeAlphabet)]
	}
	return string(code), nil
}
