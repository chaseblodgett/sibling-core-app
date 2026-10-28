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
	ListByFamilyFunc        func(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error)
	DeleteFunc              func(ctx context.Context, id string) error
	SetDocumentIDFunc       func(ctx context.Context, expenseID string, documentID *string) error
	ListSplitsByExpenseFunc func(ctx context.Context, expenseID string) ([]*models.ExpenseSplit, error)
	ListSplitsByFamilyFunc  func(ctx context.Context, familyID string) ([]*models.ExpenseSplit, error)
	GetSplitByIDFunc        func(ctx context.Context, id string) (*models.ExpenseSplit, error)
	SetSplitSettledFunc     func(ctx context.Context, splitID string, settled bool) error
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

func (m *MockExpenseRepository) ListByFamily(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error) {
	if m.ListByFamilyFunc == nil {
		return nil, nil
	}
	return m.ListByFamilyFunc(ctx, familyID, parentID)
}

func (m *MockExpenseRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc == nil {
		return nil
	}
	return m.DeleteFunc(ctx, id)
}

func (m *MockExpenseRepository) SetDocumentID(ctx context.Context, expenseID string, documentID *string) error {
	if m.SetDocumentIDFunc == nil {
		return nil
	}
	return m.SetDocumentIDFunc(ctx, expenseID, documentID)
}

func (m *MockExpenseRepository) ListSplitsByExpense(ctx context.Context, expenseID string) ([]*models.ExpenseSplit, error) {
	if m.ListSplitsByExpenseFunc == nil {
		return nil, nil
	}
	return m.ListSplitsByExpenseFunc(ctx, expenseID)
}

func (m *MockExpenseRepository) ListSplitsByFamily(ctx context.Context, familyID string) ([]*models.ExpenseSplit, error) {
	if m.ListSplitsByFamilyFunc == nil {
		return nil, nil
	}
	return m.ListSplitsByFamilyFunc(ctx, familyID)
}

func (m *MockExpenseRepository) GetSplitByID(ctx context.Context, id string) (*models.ExpenseSplit, error) {
	if m.GetSplitByIDFunc == nil {
		return nil, nil
	}
	return m.GetSplitByIDFunc(ctx, id)
}

func (m *MockExpenseRepository) SetSplitSettled(ctx context.Context, splitID string, settled bool) error {
	if m.SetSplitSettledFunc == nil {
		return nil
	}
	return m.SetSplitSettledFunc(ctx, splitID, settled)
}
