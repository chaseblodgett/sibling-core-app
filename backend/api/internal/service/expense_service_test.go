package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
	mockrepository "sibling-core/api/test/mock/mock_repository"
)

func newTestExpenseService(t *testing.T) (service.ExpenseService, *mockrepository.MockExpenseRepository, *mockrepository.MockMemberRepository, *mockrepository.MockParentRepository) {
	t.Helper()
	expenseRepo := &mockrepository.MockExpenseRepository{}
	memberRepo := &mockrepository.MockMemberRepository{}
	parentRepo := &mockrepository.MockParentRepository{}
	activityRepo := &mockrepository.MockActivityRepository{}
	documentRepo := &mockrepository.MockDocumentRepository{}
	return service.NewExpenseService(expenseRepo, memberRepo, parentRepo, documentRepo, activityRepo), expenseRepo, memberRepo, parentRepo
}

func familyMember(id, familyID string) *models.Member {
	return &models.Member{ID: id, FamilyID: familyID}
}

func TestCreateExpenseLogsActivity(t *testing.T) {
	expenseRepo := &mockrepository.MockExpenseRepository{}
	memberRepo := &mockrepository.MockMemberRepository{}
	activityRepo := &mockrepository.MockActivityRepository{}
	svc := service.NewExpenseService(expenseRepo, memberRepo, &mockrepository.MockParentRepository{}, &mockrepository.MockDocumentRepository{}, activityRepo)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return familyMember(id, "family-1"), nil
	}
	expenseRepo.CreateFunc = func(ctx context.Context, e *models.Expense, splits []*models.ExpenseSplit) error {
		e.ID = "expense-1"
		return nil
	}
	var logged *models.ActivityEntry
	activityRepo.CreateFunc = func(ctx context.Context, entry *models.ActivityEntry) error {
		logged = entry
		return nil
	}

	_, err := svc.CreateExpense(context.Background(), models.CreateExpenseInput{
		FamilyID: "family-1", PaidBy: "member-1", Description: "Groceries", AmountCents: 4200,
		ActorID: "member-1", ActorName: "Jamie",
	})

	require.NoError(t, err)
	require.NotNil(t, logged)
	assert.Equal(t, "expense.created", logged.Verb)
	assert.Equal(t, "expense", logged.SubjectType)
	assert.Equal(t, "Groceries", logged.SubjectLabel)
	assert.Equal(t, "Jamie", logged.ActorName)
}

func TestAttachReceiptSuccess(t *testing.T) {
	expenseRepo := &mockrepository.MockExpenseRepository{}
	documentRepo := &mockrepository.MockDocumentRepository{}
	svc := service.NewExpenseService(expenseRepo, &mockrepository.MockMemberRepository{}, &mockrepository.MockParentRepository{}, documentRepo, &mockrepository.MockActivityRepository{})

	expenseRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Expense, error) {
		return &models.Expense{ID: id, FamilyID: "family-1"}, nil
	}
	documentRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Document, error) {
		return &models.Document{ID: id, FamilyID: "family-1"}, nil
	}
	var setDocID *string
	expenseRepo.SetDocumentIDFunc = func(ctx context.Context, expenseID string, documentID *string) error {
		setDocID = documentID
		return nil
	}

	expense, err := svc.AttachReceipt(context.Background(), "expense-1", "doc-1")

	require.NoError(t, err)
	require.NotNil(t, setDocID)
	assert.Equal(t, "doc-1", *setDocID)
	require.NotNil(t, expense.DocumentID)
	assert.Equal(t, "doc-1", *expense.DocumentID)
}

func TestAttachReceiptRejectsDocumentFromAnotherFamily(t *testing.T) {
	expenseRepo := &mockrepository.MockExpenseRepository{}
	documentRepo := &mockrepository.MockDocumentRepository{}
	svc := service.NewExpenseService(expenseRepo, &mockrepository.MockMemberRepository{}, &mockrepository.MockParentRepository{}, documentRepo, &mockrepository.MockActivityRepository{})

	expenseRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Expense, error) {
		return &models.Expense{ID: id, FamilyID: "family-1"}, nil
	}
	documentRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Document, error) {
		return &models.Document{ID: id, FamilyID: "other-family"}, nil
	}

	_, err := svc.AttachReceipt(context.Background(), "expense-1", "doc-1")

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestRemoveReceiptClearsDocumentID(t *testing.T) {
	expenseRepo := &mockrepository.MockExpenseRepository{}
	svc := service.NewExpenseService(expenseRepo, &mockrepository.MockMemberRepository{}, &mockrepository.MockParentRepository{}, &mockrepository.MockDocumentRepository{}, &mockrepository.MockActivityRepository{})

	docID := "doc-1"
	expenseRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Expense, error) {
		return &models.Expense{ID: id, FamilyID: "family-1", DocumentID: &docID}, nil
	}
	var setDocID *string
	setCalled := false
	expenseRepo.SetDocumentIDFunc = func(ctx context.Context, expenseID string, documentID *string) error {
		setCalled = true
		setDocID = documentID
		return nil
	}

	expense, err := svc.RemoveReceipt(context.Background(), "expense-1")

	require.NoError(t, err)
	assert.True(t, setCalled)
	assert.Nil(t, setDocID)
	assert.Nil(t, expense.DocumentID)
}

func TestCreateExpenseSuccess(t *testing.T) {
	svc, expenseRepo, memberRepo, _ := newTestExpenseService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return familyMember(id, "family-1"), nil
	}
	var created *models.Expense
	var createdSplits []*models.ExpenseSplit
	expenseRepo.CreateFunc = func(ctx context.Context, e *models.Expense, splits []*models.ExpenseSplit) error {
		e.ID = "expense-1"
		created = e
		createdSplits = splits
		return nil
	}

	expense, err := svc.CreateExpense(context.Background(), models.CreateExpenseInput{
		FamilyID:    "family-1",
		PaidBy:      "member-1",
		Description: "  Pharmacy run ",
		AmountCents: 3000,
		Splits: []models.SplitInput{
			{MemberID: "member-1", ShareCents: 1500},
			{MemberID: "member-2", ShareCents: 1500},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "Pharmacy run", expense.Description)
	assert.Equal(t, "USD", expense.Currency)
	require.Len(t, createdSplits, 2)
}

func TestCreateExpenseRejectsZeroAmount(t *testing.T) {
	svc, _, _, _ := newTestExpenseService(t)

	_, err := svc.CreateExpense(context.Background(), models.CreateExpenseInput{
		FamilyID: "family-1", PaidBy: "member-1", Description: "Thing", AmountCents: 0,
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestCreateExpenseRejectsSplitsNotSummingToTotal(t *testing.T) {
	svc, _, memberRepo, _ := newTestExpenseService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return familyMember(id, "family-1"), nil
	}

	_, err := svc.CreateExpense(context.Background(), models.CreateExpenseInput{
		FamilyID:    "family-1",
		PaidBy:      "member-1",
		Description: "Thing",
		AmountCents: 1000,
		Splits: []models.SplitInput{
			{MemberID: "member-1", ShareCents: 400},
			{MemberID: "member-2", ShareCents: 400},
		},
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestCreateExpenseRejectsDuplicateSplitMember(t *testing.T) {
	svc, _, memberRepo, _ := newTestExpenseService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return familyMember(id, "family-1"), nil
	}

	_, err := svc.CreateExpense(context.Background(), models.CreateExpenseInput{
		FamilyID:    "family-1",
		PaidBy:      "member-1",
		Description: "Thing",
		AmountCents: 1000,
		Splits: []models.SplitInput{
			{MemberID: "member-1", ShareCents: 500},
			{MemberID: "member-1", ShareCents: 500},
		},
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestCreateExpenseRejectsSplitMemberFromAnotherFamily(t *testing.T) {
	svc, _, memberRepo, _ := newTestExpenseService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		if id == "member-2" {
			return familyMember(id, "other-family"), nil
		}
		return familyMember(id, "family-1"), nil
	}

	_, err := svc.CreateExpense(context.Background(), models.CreateExpenseInput{
		FamilyID:    "family-1",
		PaidBy:      "member-1",
		Description: "Thing",
		AmountCents: 1000,
		Splits: []models.SplitInput{
			{MemberID: "member-1", ShareCents: 500},
			{MemberID: "member-2", ShareCents: 500},
		},
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestCreateExpenseAllowsNoSplits(t *testing.T) {
	svc, expenseRepo, memberRepo, _ := newTestExpenseService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return familyMember(id, "family-1"), nil
	}
	expenseRepo.CreateFunc = func(ctx context.Context, e *models.Expense, splits []*models.ExpenseSplit) error {
		e.ID = "expense-1"
		return nil
	}

	expense, err := svc.CreateExpense(context.Background(), models.CreateExpenseInput{
		FamilyID: "family-1", PaidBy: "member-1", Description: "Just tracking", AmountCents: 500,
	})

	require.NoError(t, err)
	assert.Equal(t, "expense-1", expense.ID)
}

func TestCreateExpenseRejectsPaidByFromAnotherFamily(t *testing.T) {
	svc, _, memberRepo, _ := newTestExpenseService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return familyMember(id, "other-family"), nil
	}

	_, err := svc.CreateExpense(context.Background(), models.CreateExpenseInput{
		FamilyID: "family-1", PaidBy: "member-1", Description: "Thing", AmountCents: 500,
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestCreateExpenseValidatesParentFamily(t *testing.T) {
	svc, _, memberRepo, parentRepo := newTestExpenseService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return familyMember(id, "family-1"), nil
	}
	parentRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Parent, error) {
		return &models.Parent{ID: id, FamilyID: "other-family"}, nil
	}

	parentID := "parent-1"
	_, err := svc.CreateExpense(context.Background(), models.CreateExpenseInput{
		FamilyID: "family-1", PaidBy: "member-1", Description: "Thing", AmountCents: 500, ParentID: &parentID,
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestGetBalancesSkipsSettledAndSelfShares(t *testing.T) {
	svc, expenseRepo, _, _ := newTestExpenseService(t)

	expenseRepo.ListByFamilyFunc = func(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error) {
		return []*models.Expense{
			{ID: "exp-1", PaidBy: "alice"},
			{ID: "exp-2", PaidBy: "bob"},
		}, nil
	}
	expenseRepo.ListSplitsByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.ExpenseSplit, error) {
		return []*models.ExpenseSplit{
			// exp-1 paid by alice, split between alice (self, nets 0) and bob (owes alice).
			{ExpenseID: "exp-1", MemberID: "alice", ShareCents: 500, Settled: false},
			{ExpenseID: "exp-1", MemberID: "bob", ShareCents: 500, Settled: false},
			// exp-2 paid by bob, charlie's share already settled — shouldn't count.
			{ExpenseID: "exp-2", MemberID: "charlie", ShareCents: 1000, Settled: true},
		}, nil
	}

	balances, err := svc.GetBalances(context.Background(), "family-1")

	require.NoError(t, err)
	require.Len(t, balances, 2)
	byMember := make(map[string]int64)
	for _, b := range balances {
		byMember[b.MemberID] = b.NetCents
	}
	assert.Equal(t, int64(500), byMember["alice"])
	assert.Equal(t, int64(-500), byMember["bob"])
	_, charliePresent := byMember["charlie"]
	assert.False(t, charliePresent, "settled splits shouldn't contribute to balance")
}

func TestSettleSplit(t *testing.T) {
	svc, expenseRepo, _, _ := newTestExpenseService(t)

	expenseRepo.GetSplitByIDFunc = func(ctx context.Context, id string) (*models.ExpenseSplit, error) {
		return &models.ExpenseSplit{ID: id, ExpenseID: "exp-1", MemberID: "bob", ShareCents: 500, Settled: false}, nil
	}
	var settledID string
	expenseRepo.SetSplitSettledFunc = func(ctx context.Context, splitID string, settled bool) error {
		settledID = splitID
		assert.True(t, settled)
		return nil
	}

	split, err := svc.SettleSplit(context.Background(), "split-1")

	require.NoError(t, err)
	assert.Equal(t, "split-1", settledID)
	assert.True(t, split.Settled)
}

func TestDeleteExpenseNotFound(t *testing.T) {
	svc, expenseRepo, _, _ := newTestExpenseService(t)

	expenseRepo.DeleteFunc = func(ctx context.Context, id string) error {
		return repository.ErrExpenseNotFound
	}

	err := svc.DeleteExpense(context.Background(), "expense-1")

	assert.ErrorIs(t, err, repository.ErrExpenseNotFound)
}
