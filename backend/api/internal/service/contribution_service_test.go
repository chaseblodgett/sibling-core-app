package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
	mockrepository "sibling-core/api/test/mock/mock_repository"
)

func newTestContributionService(t *testing.T) (service.ContributionService, *mockrepository.MockMemberRepository, *mockrepository.MockExpenseRepository, *mockrepository.MockActivityRepository) {
	t.Helper()
	memberRepo := &mockrepository.MockMemberRepository{}
	expenseRepo := &mockrepository.MockExpenseRepository{}
	activityRepo := &mockrepository.MockActivityRepository{}
	return service.NewContributionService(memberRepo, expenseRepo, activityRepo), memberRepo, expenseRepo, activityRepo
}

func actorID(id string) *string { return &id }

func TestGetContributionsTalliesTasksAndExpenses(t *testing.T) {
	svc, memberRepo, expenseRepo, activityRepo := newTestContributionService(t)

	memberRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Member, error) {
		return []*models.Member{
			{ID: "member-1", FamilyID: familyID, FullName: "Jamie"},
			{ID: "member-2", FamilyID: familyID, FullName: "Sarah"},
		}, nil
	}
	activityRepo.ListByFamilyFunc = func(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error) {
		return []*models.ActivityEntry{
			{Verb: "task.completed", ActorID: actorID("member-1")},
			{Verb: "task.completed", ActorID: actorID("member-1")},
			{Verb: "task.completed", ActorID: actorID("member-2")},
			{Verb: "task.created", ActorID: actorID("member-1")}, // not a completion, shouldn't count
		}, nil
	}
	expenseRepo.ListByFamilyFunc = func(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error) {
		return []*models.Expense{
			{PaidBy: "member-2", AmountCents: 3000},
			{PaidBy: "member-2", AmountCents: 1500},
		}, nil
	}

	contributions, err := svc.GetContributions(context.Background(), "family-1")

	require.NoError(t, err)
	require.Len(t, contributions, 2)

	// member-1 has more tasks completed (2 vs 1), so ranks first despite
	// member-2 having fronted more money.
	assert.Equal(t, "member-1", contributions[0].MemberID)
	assert.Equal(t, 2, contributions[0].TasksCompleted)
	assert.Equal(t, 0, contributions[0].ExpensesPaidCount)
	assert.Equal(t, int64(0), contributions[0].ExpensesPaidCents)

	assert.Equal(t, "member-2", contributions[1].MemberID)
	assert.Equal(t, 1, contributions[1].TasksCompleted)
	assert.Equal(t, 2, contributions[1].ExpensesPaidCount)
	assert.Equal(t, int64(4500), contributions[1].ExpensesPaidCents)
}

func TestGetContributionsIncludesZeroContributionMembers(t *testing.T) {
	svc, memberRepo, expenseRepo, activityRepo := newTestContributionService(t)

	memberRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Member, error) {
		return []*models.Member{{ID: "member-1", FamilyID: familyID, FullName: "Jamie"}}, nil
	}
	activityRepo.ListByFamilyFunc = func(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error) {
		return nil, nil
	}
	expenseRepo.ListByFamilyFunc = func(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error) {
		return nil, nil
	}

	contributions, err := svc.GetContributions(context.Background(), "family-1")

	require.NoError(t, err)
	require.Len(t, contributions, 1)
	assert.Equal(t, "Jamie", contributions[0].MemberName)
	assert.Equal(t, 0, contributions[0].TasksCompleted)
	assert.Equal(t, 0, contributions[0].ExpensesPaidCount)
}

func TestGetContributionsFallsBackToEmailWhenNoFullName(t *testing.T) {
	svc, memberRepo, _, _ := newTestContributionService(t)

	memberRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Member, error) {
		return []*models.Member{{ID: "member-1", FamilyID: familyID, Email: "jamie@example.com"}}, nil
	}

	contributions, err := svc.GetContributions(context.Background(), "family-1")

	require.NoError(t, err)
	require.Len(t, contributions, 1)
	assert.Equal(t, "jamie@example.com", contributions[0].MemberName)
}

func TestGetContributionsIgnoresNilActorID(t *testing.T) {
	svc, memberRepo, _, activityRepo := newTestContributionService(t)

	memberRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Member, error) {
		return []*models.Member{{ID: "member-1", FamilyID: familyID, FullName: "Jamie"}}, nil
	}
	activityRepo.ListByFamilyFunc = func(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error) {
		return []*models.ActivityEntry{{Verb: "task.completed", ActorID: nil}}, nil
	}

	contributions, err := svc.GetContributions(context.Background(), "family-1")

	require.NoError(t, err)
	assert.Equal(t, 0, contributions[0].TasksCompleted)
}
