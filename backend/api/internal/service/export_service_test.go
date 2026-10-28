package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
	mockrepository "sibling-core/api/test/mock/mock_repository"
)

func newTestExportService(t *testing.T) (service.ExportService, *mockrepository.MockExpenseRepository, *mockrepository.MockTaskRepository, *mockrepository.MockMemberRepository, *mockrepository.MockParentRepository) {
	t.Helper()
	expenseRepo := &mockrepository.MockExpenseRepository{}
	taskRepo := &mockrepository.MockTaskRepository{}
	memberRepo := &mockrepository.MockMemberRepository{}
	parentRepo := &mockrepository.MockParentRepository{}
	return service.NewExportService(expenseRepo, taskRepo, memberRepo, parentRepo), expenseRepo, taskRepo, memberRepo, parentRepo
}

func TestExportExpensesIncludesHeaderAndRows(t *testing.T) {
	svc, expenseRepo, _, memberRepo, parentRepo := newTestExportService(t)

	incurred := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	expenseRepo.ListByFamilyFunc = func(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error) {
		return []*models.Expense{
			{ID: "expense-1", PaidBy: "member-1", Description: "Groceries", Category: "Food", AmountCents: 4200, Currency: "USD", IncurredAt: incurred},
		}, nil
	}
	expenseRepo.ListSplitsByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.ExpenseSplit, error) {
		return []*models.ExpenseSplit{
			{ID: "split-1", ExpenseID: "expense-1", MemberID: "member-1", ShareCents: 2100, Settled: true},
			{ID: "split-2", ExpenseID: "expense-1", MemberID: "member-2", ShareCents: 2100, Settled: false},
		}, nil
	}
	memberRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Member, error) {
		return []*models.Member{
			{ID: "member-1", FullName: "Jamie"},
			{ID: "member-2", FullName: "Sarah"},
		}, nil
	}
	parentRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Parent, error) {
		return nil, nil
	}

	csvBytes, err := svc.ExportExpenses(context.Background(), "family-1", nil, nil, nil)

	require.NoError(t, err)
	csvStr := string(csvBytes)
	assert.Contains(t, csvStr, "date,description,category,amount,currency,paid_by,parent,splits")
	assert.Contains(t, csvStr, "Groceries")
	assert.Contains(t, csvStr, "$42.00")
	assert.Contains(t, csvStr, "Jamie")
	assert.Contains(t, csvStr, "settled")
	assert.Contains(t, csvStr, "unsettled")
}

func TestExportExpensesFiltersByDateRange(t *testing.T) {
	svc, expenseRepo, _, memberRepo, parentRepo := newTestExportService(t)

	inRange := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	outOfRange := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	expenseRepo.ListByFamilyFunc = func(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error) {
		return []*models.Expense{
			{ID: "e1", PaidBy: "member-1", Description: "In range", AmountCents: 100, Currency: "USD", IncurredAt: inRange},
			{ID: "e2", PaidBy: "member-1", Description: "Out of range", AmountCents: 100, Currency: "USD", IncurredAt: outOfRange},
		}, nil
	}
	expenseRepo.ListSplitsByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.ExpenseSplit, error) { return nil, nil }
	memberRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Member, error) { return nil, nil }
	parentRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Parent, error) { return nil, nil }

	from := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC)
	csvBytes, err := svc.ExportExpenses(context.Background(), "family-1", nil, &from, &to)

	require.NoError(t, err)
	csvStr := string(csvBytes)
	assert.Contains(t, csvStr, "In range")
	assert.NotContains(t, csvStr, "Out of range")
}

func TestExportTasksIncludesHeaderAndRows(t *testing.T) {
	svc, _, taskRepo, memberRepo, parentRepo := newTestExportService(t)

	created := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	assigneeID := "member-1"
	taskRepo.ListByFamilyFunc = func(ctx context.Context, familyID string, parentID *string) ([]*models.Task, error) {
		return []*models.Task{
			{ID: "task-1", Title: "Refill prescription", Status: "done", AssigneeID: &assigneeID, RecurrenceRule: "weekly", CreatedAt: created},
		}, nil
	}
	memberRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Member, error) {
		return []*models.Member{{ID: "member-1", FullName: "Jamie"}}, nil
	}
	parentRepo.ListByFamilyFunc = func(ctx context.Context, familyID string) ([]*models.Parent, error) { return nil, nil }

	csvBytes, err := svc.ExportTasks(context.Background(), "family-1", nil, nil, nil)

	require.NoError(t, err)
	csvStr := string(csvBytes)
	assert.Contains(t, csvStr, "created,title,description,status,assignee,parent,due_date,repeats")
	assert.Contains(t, csvStr, "Refill prescription")
	assert.Contains(t, csvStr, "Jamie")
	assert.Contains(t, csvStr, "weekly")
}
