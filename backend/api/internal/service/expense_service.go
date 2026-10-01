package service

import (
	"context"
	"errors"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// ExpenseService contains business logic for expense tracking, splitting, and balances.
type ExpenseService interface {
	CreateExpense(ctx context.Context, input models.CreateExpenseInput) (*models.Expense, error)
	GetExpense(ctx context.Context, id string) (*models.Expense, error)
	ListExpenses(ctx context.Context, familyID string) ([]*models.Expense, error)
	DeleteExpense(ctx context.Context, id string) error

	SettleSplit(ctx context.Context, splitID string) error

	// GetBalances returns each member's running net balance for a family.
	GetBalances(ctx context.Context, familyID string) ([]*models.Balance, error)
}

type expenseService struct {
	expenseRepo repository.ExpenseRepository
}

// NewExpenseService constructs an ExpenseService.
func NewExpenseService(expenseRepo repository.ExpenseRepository) ExpenseService {
	return &expenseService{expenseRepo: expenseRepo}
}

func (s *expenseService) CreateExpense(ctx context.Context, input models.CreateExpenseInput) (*models.Expense, error) {
	return nil, errors.New("not implemented")
}

func (s *expenseService) GetExpense(ctx context.Context, id string) (*models.Expense, error) {
	return nil, errors.New("not implemented")
}

func (s *expenseService) ListExpenses(ctx context.Context, familyID string) ([]*models.Expense, error) {
	return nil, errors.New("not implemented")
}

func (s *expenseService) DeleteExpense(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (s *expenseService) SettleSplit(ctx context.Context, splitID string) error {
	return errors.New("not implemented")
}

func (s *expenseService) GetBalances(ctx context.Context, familyID string) ([]*models.Balance, error) {
	return nil, errors.New("not implemented")
}
