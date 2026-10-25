package mockrepository

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockExpenseRepository is a test double for repository.ExpenseRepository.
// Assign the *Func fields your test exercises; unset ones return zero values.
type MockExpenseRepository struct {
	CreateFunc              func(ctx context.Context, expense *models.Expense, splits []*models.ExpenseSplit) error
	GetByIDFunc             func(ctx context.Context, id string) (*models.Expense, error)
	ListByFamilyFunc        func(ctx context.Context, familyID string) ([]*models.Expense, error)
	UpdateFunc              func(ctx context.Context, expense *models.Expense) error
	DeleteFunc              func(ctx context.Context, id string) error
	ListSplitsByExpenseFunc func(ctx context.Context, expenseID string) ([]*models.ExpenseSplit, error)
	SetSplitSettledFunc     func(ctx context.Context, splitID string, settled bool) error
	GetBalancesFunc         func(ctx context.Context, familyID string) ([]*models.Balance, error)
}

func (m *MockExpenseRepository) Create(ctx context.Context, expense *models.Expense, splits []*models.ExpenseSplit) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, expense, splits)
}

func (m *MockExpenseRepository) GetByID(ctx context.Context, id string) (*models.Expense, error) {
	if m.GetByIDFunc == nil {
		return nil, nil
	}
	return m.GetByIDFunc(ctx, id)
}

func (m *MockExpenseRepository) ListByFamily(ctx context.Context, familyID string) ([]*models.Expense, error) {
	if m.ListByFamilyFunc == nil {
		return nil, nil
	}
	return m.ListByFamilyFunc(ctx, familyID)
}

func (m *MockExpenseRepository) Update(ctx context.Context, expense *models.Expense) error {
	if m.UpdateFunc == nil {
		return nil
	}
	return m.UpdateFunc(ctx, expense)
}

func (m *MockExpenseRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc == nil {
		return nil
	}
	return m.DeleteFunc(ctx, id)
}

func (m *MockExpenseRepository) ListSplitsByExpense(ctx context.Context, expenseID string) ([]*models.ExpenseSplit, error) {
	if m.ListSplitsByExpenseFunc == nil {
		return nil, nil
	}
	return m.ListSplitsByExpenseFunc(ctx, expenseID)
}

func (m *MockExpenseRepository) SetSplitSettled(ctx context.Context, splitID string, settled bool) error {
	if m.SetSplitSettledFunc == nil {
		return nil
	}
	return m.SetSplitSettledFunc(ctx, splitID, settled)
}

func (m *MockExpenseRepository) GetBalances(ctx context.Context, familyID string) ([]*models.Balance, error) {
	if m.GetBalancesFunc == nil {
		return nil, nil
	}
	return m.GetBalancesFunc(ctx, familyID)
}
