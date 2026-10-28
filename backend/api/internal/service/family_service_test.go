package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sibling-core/api/internal/mailer"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
	mockmailer "sibling-core/api/test/mock/mock_mailer"
	mockrepository "sibling-core/api/test/mock/mock_repository"
)

func newTestFamilyService(t *testing.T) (service.FamilyService, *mockrepository.MockFamilyRepository, *mockrepository.MockMemberRepository, *mockrepository.MockUserRepository) {
	t.Helper()
	svc, familyRepo, memberRepo, userRepo, _ := newTestFamilyServiceWithMailer(t)
	return svc, familyRepo, memberRepo, userRepo
}

func newTestFamilyServiceWithMailer(t *testing.T) (service.FamilyService, *mockrepository.MockFamilyRepository, *mockrepository.MockMemberRepository, *mockrepository.MockUserRepository, *mockmailer.MockMailer) {
	t.Helper()
	familyRepo := &mockrepository.MockFamilyRepository{}
	memberRepo := &mockrepository.MockMemberRepository{}
	userRepo := &mockrepository.MockUserRepository{}
	mail := &mockmailer.MockMailer{}
	activityRepo := &mockrepository.MockActivityRepository{}
	return service.NewFamilyService(familyRepo, memberRepo, userRepo, mail, activityRepo), familyRepo, memberRepo, userRepo, mail
}

func TestCreateFamilySuccess(t *testing.T) {
	svc, familyRepo, memberRepo, userRepo := newTestFamilyService(t)
	ctx := context.Background()

	userRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
		return &models.User{ID: "user-1", Email: "owner@example.com", Name: "Owner"}, nil
	}

	var createdFamily *models.Family
	familyRepo.CreateFunc = func(ctx context.Context, family *models.Family) error {
		family.ID = "family-1"
		createdFamily = family
		return nil
	}

	var createdMember *models.Member
	memberRepo.CreateFunc = func(ctx context.Context, member *models.Member) error {
		member.ID = "member-1"
		createdMember = member
		return nil
	}

	family, err := svc.CreateFamily(ctx, models.CreateFamilyInput{Name: "The Does", CreatedBy: "user-1"})

	require.NoError(t, err)
	assert.Equal(t, "The Does", family.Name)
	assert.Len(t, family.InviteCode, 8)
	require.NotNil(t, createdFamily)
	require.NotNil(t, createdMember)
	assert.Equal(t, "family-1", createdMember.FamilyID)
	assert.Equal(t, "owner@example.com", createdMember.Email)
	assert.Equal(t, "owner", createdMember.Role)
}

func TestCreateFamilyBlankName(t *testing.T) {
	svc, _, _, _ := newTestFamilyService(t)

	_, err := svc.CreateFamily(context.Background(), models.CreateFamilyInput{Name: "   ", CreatedBy: "user-1"})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestCreateFamilyRetriesOnInviteCodeCollision(t *testing.T) {
	svc, familyRepo, memberRepo, userRepo := newTestFamilyService(t)
	ctx := context.Background()

	userRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
		return &models.User{ID: "user-1", Email: "owner@example.com", Name: "Owner"}, nil
	}
	memberRepo.CreateFunc = func(ctx context.Context, member *models.Member) error {
		member.ID = "member-1"
		return nil
	}

	attempts := 0
	familyRepo.CreateFunc = func(ctx context.Context, family *models.Family) error {
		attempts++
		if attempts < 3 {
			return repository.ErrDuplicateInviteCode
		}
		family.ID = "family-1"
		return nil
	}

	family, err := svc.CreateFamily(ctx, models.CreateFamilyInput{Name: "The Does", CreatedBy: "user-1"})

	require.NoError(t, err)
	assert.Equal(t, 3, attempts)
	assert.Equal(t, "family-1", family.ID)
}

func TestJoinFamilySuccess(t *testing.T) {
	svc, familyRepo, memberRepo, userRepo := newTestFamilyService(t)
	ctx := context.Background()

	family := &models.Family{ID: "family-1", InviteCode: "ABCD2345"}
	familyRepo.GetByInviteCodeFunc = func(ctx context.Context, code string) (*models.Family, error) {
		assert.Equal(t, "ABCD2345", code)
		return family, nil
	}
	memberRepo.GetByUserAndFamilyFunc = func(ctx context.Context, userID, familyID string) (*models.Member, error) {
		return nil, repository.ErrMemberNotFound
	}
	userRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
		return &models.User{ID: "user-2", Email: "sibling@example.com", Name: "Sibling"}, nil
	}

	var createdMember *models.Member
	memberRepo.CreateFunc = func(ctx context.Context, member *models.Member) error {
		member.ID = "member-2"
		createdMember = member
		return nil
	}
	familyRepo.ListInvitesByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.FamilyInvite, error) {
		return []*models.FamilyInvite{
			{ID: "invite-1", Email: "sibling@example.com", Status: "pending"},
		}, nil
	}
	var acceptedInviteID string
	familyRepo.UpdateInviteStatusFunc = func(ctx context.Context, id string, status string) error {
		acceptedInviteID = id
		assert.Equal(t, "accepted", status)
		return nil
	}

	got, err := svc.JoinFamily(ctx, " abcd2345 ", "user-2")

	require.NoError(t, err)
	assert.Equal(t, family, got)
	require.NotNil(t, createdMember)
	assert.Equal(t, "member", createdMember.Role)
	assert.Equal(t, "invite-1", acceptedInviteID)
}

func TestJoinFamilyLogsActivity(t *testing.T) {
	familyRepo := &mockrepository.MockFamilyRepository{}
	memberRepo := &mockrepository.MockMemberRepository{}
	userRepo := &mockrepository.MockUserRepository{}
	activityRepo := &mockrepository.MockActivityRepository{}
	svc := service.NewFamilyService(familyRepo, memberRepo, userRepo, &mockmailer.MockMailer{}, activityRepo)
	ctx := context.Background()

	family := &models.Family{ID: "family-1", InviteCode: "ABCD2345"}
	familyRepo.GetByInviteCodeFunc = func(ctx context.Context, code string) (*models.Family, error) { return family, nil }
	memberRepo.GetByUserAndFamilyFunc = func(ctx context.Context, userID, familyID string) (*models.Member, error) {
		return nil, repository.ErrMemberNotFound
	}
	userRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
		return &models.User{ID: "user-2", Email: "sibling@example.com", Name: "Sibling"}, nil
	}
	memberRepo.CreateFunc = func(ctx context.Context, member *models.Member) error {
		member.ID = "member-2"
		return nil
	}
	familyRepo.ListInvitesByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.FamilyInvite, error) {
		return nil, nil
	}
	var logged *models.ActivityEntry
	activityRepo.CreateFunc = func(ctx context.Context, entry *models.ActivityEntry) error {
		logged = entry
		return nil
	}

	_, err := svc.JoinFamily(ctx, "abcd2345", "user-2")

	require.NoError(t, err)
	require.NotNil(t, logged)
	assert.Equal(t, "family-1", logged.FamilyID)
	assert.Equal(t, "member.joined", logged.Verb)
	assert.Equal(t, "Sibling", logged.ActorName)
	require.NotNil(t, logged.ActorID)
	assert.Equal(t, "member-2", *logged.ActorID)
}

func TestJoinFamilyAlreadyMemberIsIdempotent(t *testing.T) {
	svc, familyRepo, memberRepo, _ := newTestFamilyService(t)
	ctx := context.Background()

	family := &models.Family{ID: "family-1", InviteCode: "ABCD2345"}
	familyRepo.GetByInviteCodeFunc = func(ctx context.Context, code string) (*models.Family, error) {
		return family, nil
	}
	memberRepo.GetByUserAndFamilyFunc = func(ctx context.Context, userID, familyID string) (*models.Member, error) {
		return &models.Member{ID: "member-2"}, nil
	}

	createCalled := false
	memberRepo.CreateFunc = func(ctx context.Context, member *models.Member) error {
		createCalled = true
		return nil
	}

	got, err := svc.JoinFamily(ctx, "ABCD2345", "user-2")

	require.NoError(t, err)
	assert.Equal(t, family, got)
	assert.False(t, createCalled, "should not create a duplicate membership")
}

func TestJoinFamilyInvalidCode(t *testing.T) {
	svc, familyRepo, _, _ := newTestFamilyService(t)

	familyRepo.GetByInviteCodeFunc = func(ctx context.Context, code string) (*models.Family, error) {
		return nil, repository.ErrFamilyNotFound
	}

	_, err := svc.JoinFamily(context.Background(), "NOPE0000", "user-2")

	assert.ErrorIs(t, err, service.ErrInvalidInviteCode)
}

func TestInviteMemberSuccess(t *testing.T) {
	svc, familyRepo, _, userRepo, mail := newTestFamilyServiceWithMailer(t)

	var created *models.FamilyInvite
	familyRepo.CreateInviteFunc = func(ctx context.Context, invite *models.FamilyInvite) error {
		invite.ID = "invite-1"
		created = invite
		return nil
	}
	familyRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Family, error) {
		return &models.Family{ID: id, Name: "The Does", InviteCode: "ABCD2345"}, nil
	}
	userRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
		return &models.User{ID: id, Name: "Owner Person"}, nil
	}

	var sent mailer.FamilyInviteEmail
	sendCalled := false
	mail.SendFamilyInviteFunc = func(ctx context.Context, email mailer.FamilyInviteEmail) error {
		sendCalled = true
		sent = email
		return nil
	}

	invite, err := svc.InviteMember(context.Background(), models.InviteMemberInput{
		FamilyID:  "family-1",
		Email:     "  Sibling@Example.com ",
		InvitedBy: "user-1",
	})

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "sibling@example.com", invite.Email)
	assert.Equal(t, "pending", invite.Status)

	require.True(t, sendCalled)
	assert.Equal(t, "sibling@example.com", sent.To)
	assert.Equal(t, "The Does", sent.FamilyName)
	assert.Equal(t, "ABCD2345", sent.InviteCode)
	assert.Equal(t, "Owner Person", sent.InviterName)
}

func TestInviteMemberSucceedsEvenIfEmailSendFails(t *testing.T) {
	svc, familyRepo, _, userRepo, mail := newTestFamilyServiceWithMailer(t)

	familyRepo.CreateInviteFunc = func(ctx context.Context, invite *models.FamilyInvite) error {
		invite.ID = "invite-1"
		return nil
	}
	familyRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Family, error) {
		return &models.Family{ID: id, Name: "The Does", InviteCode: "ABCD2345"}, nil
	}
	userRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
		return &models.User{ID: id, Name: "Owner Person"}, nil
	}
	mail.SendFamilyInviteFunc = func(ctx context.Context, email mailer.FamilyInviteEmail) error {
		return errors.New("resend is down")
	}

	invite, err := svc.InviteMember(context.Background(), models.InviteMemberInput{
		FamilyID:  "family-1",
		Email:     "sibling@example.com",
		InvitedBy: "user-1",
	})

	require.NoError(t, err)
	assert.Equal(t, "sibling@example.com", invite.Email)
}

func TestInviteMemberInvalidEmail(t *testing.T) {
	svc, _, _, _ := newTestFamilyService(t)

	_, err := svc.InviteMember(context.Background(), models.InviteMemberInput{
		FamilyID:  "family-1",
		Email:     "not-an-email",
		InvitedBy: "user-1",
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestRemoveMemberRejectsCrossFamily(t *testing.T) {
	svc, _, memberRepo, _ := newTestFamilyService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "other-family"}, nil
	}
	deleteCalled := false
	memberRepo.DeleteFunc = func(ctx context.Context, id string) error {
		deleteCalled = true
		return nil
	}

	err := svc.RemoveMember(context.Background(), "family-1", "member-9")

	assert.ErrorIs(t, err, repository.ErrMemberNotFound)
	assert.False(t, deleteCalled)
}

func TestUpdateMemberRoleSuccess(t *testing.T) {
	svc, _, memberRepo, _ := newTestFamilyService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "family-1", Role: "member"}, nil
	}
	var setRole string
	memberRepo.SetRoleFunc = func(ctx context.Context, id, role string) error {
		setRole = role
		return nil
	}

	member, err := svc.UpdateMemberRole(context.Background(), "family-1", "member-1", "viewer")

	require.NoError(t, err)
	assert.Equal(t, "viewer", setRole)
	assert.Equal(t, "viewer", member.Role)
}

func TestUpdateMemberRoleRejectsInvalidRole(t *testing.T) {
	svc, _, memberRepo, _ := newTestFamilyService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "family-1", Role: "member"}, nil
	}

	_, err := svc.UpdateMemberRole(context.Background(), "family-1", "member-1", "admin")

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestUpdateMemberRoleRejectsCrossFamily(t *testing.T) {
	svc, _, memberRepo, _ := newTestFamilyService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "other-family", Role: "member"}, nil
	}

	_, err := svc.UpdateMemberRole(context.Background(), "family-1", "member-1", "viewer")

	assert.ErrorIs(t, err, repository.ErrMemberNotFound)
}

func TestUpdateMemberRoleRejectsDemotingLastOwner(t *testing.T) {
	svc, _, memberRepo, _ := newTestFamilyService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "family-1", Role: "owner"}, nil
	}
	memberRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Member, error) {
		return []*models.Member{
			{ID: "member-1", FamilyID: familyID, Role: "owner"},
			{ID: "member-2", FamilyID: familyID, Role: "member"},
		}, nil
	}
	setRoleCalled := false
	memberRepo.SetRoleFunc = func(ctx context.Context, id, role string) error {
		setRoleCalled = true
		return nil
	}

	_, err := svc.UpdateMemberRole(context.Background(), "family-1", "member-1", "member")

	assert.ErrorIs(t, err, service.ErrInvalidInput)
	assert.False(t, setRoleCalled)
}

func TestUpdateMemberRoleAllowsDemotingOwnerWithAnotherOwnerPresent(t *testing.T) {
	svc, _, memberRepo, _ := newTestFamilyService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "family-1", Role: "owner"}, nil
	}
	memberRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Member, error) {
		return []*models.Member{
			{ID: "member-1", FamilyID: familyID, Role: "owner"},
			{ID: "member-2", FamilyID: familyID, Role: "owner"},
		}, nil
	}
	setRoleCalled := false
	memberRepo.SetRoleFunc = func(ctx context.Context, id, role string) error {
		setRoleCalled = true
		return nil
	}

	_, err := svc.UpdateMemberRole(context.Background(), "family-1", "member-1", "member")

	require.NoError(t, err)
	assert.True(t, setRoleCalled)
}

func TestRequireMembership(t *testing.T) {
	svc, _, memberRepo, _ := newTestFamilyService(t)

	memberRepo.GetByUserAndFamilyFunc = func(ctx context.Context, userID, familyID string) (*models.Member, error) {
		return nil, repository.ErrMemberNotFound
	}

	_, err := svc.RequireMembership(context.Background(), "user-1", "family-1")

	assert.ErrorIs(t, err, service.ErrNotMember)
}

func TestListFamiliesForUser(t *testing.T) {
	svc, familyRepo, memberRepo, _ := newTestFamilyService(t)

	memberRepo.ListByUserFunc = func(ctx context.Context, userID string) ([]*models.Member, error) {
		return []*models.Member{{FamilyID: "family-1"}, {FamilyID: "family-2"}}, nil
	}
	familyRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Family, error) {
		return &models.Family{ID: id}, nil
	}

	families, err := svc.ListFamiliesForUser(context.Background(), "user-1")

	require.NoError(t, err)
	require.Len(t, families, 2)
	assert.Equal(t, "family-1", families[0].ID)
	assert.Equal(t, "family-2", families[1].ID)
}
