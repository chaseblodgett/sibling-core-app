package mockservice

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockExpenseService is a test double for service.ExpenseService. Assign
// the *Func fields your test exercises; unset ones return zero values.
type MockExpenseService struct {
	CreateExpenseFunc func(ctx context.Context, input models.CreateExpenseInput) (*models.Expense, error)
	GetExpenseFunc    func(ctx context.Context, id string) (*models.Expense, error)
	ListExpensesFunc  func(ctx context.Context, familyID string) ([]*models.Expense, error)
	DeleteExpenseFunc func(ctx context.Context, id string) error
	SettleSplitFunc   func(ctx context.Context, splitID string) error
	GetBalancesFunc   func(ctx context.Context, familyID string) ([]*models.Balance, error)
}

func (m *MockExpenseService) CreateExpense(ctx context.Context, input models.CreateExpenseInput) (*models.Expense, error) {
	if m.CreateExpenseFunc == nil {
		return nil, nil
	}
	return m.CreateExpenseFunc(ctx, input)
}

func (m *MockExpenseService) GetExpense(ctx context.Context, id string) (*models.Expense, error) {
	if m.GetExpenseFunc == nil {
		return nil, nil
	}
	return m.GetExpenseFunc(ctx, id)
}

func (m *MockExpenseService) ListExpenses(ctx context.Context, familyID string) ([]*models.Expense, error) {
	if m.ListExpensesFunc == nil {
		return nil, nil
	}
	return m.ListExpensesFunc(ctx, familyID)
}

func (m *MockExpenseService) DeleteExpense(ctx context.Context, id string) error {
	if m.DeleteExpenseFunc == nil {
		return nil
	}
	return m.DeleteExpenseFunc(ctx, id)
}

func (m *MockExpenseService) SettleSplit(ctx context.Context, splitID string) error {
	if m.SettleSplitFunc == nil {
		return nil
	}
	return m.SettleSplitFunc(ctx, splitID)
}

func (m *MockExpenseService) GetBalances(ctx context.Context, familyID string) ([]*models.Balance, error) {
	if m.GetBalancesFunc == nil {
		return nil, nil
	}
	return m.GetBalancesFunc(ctx, familyID)
}
