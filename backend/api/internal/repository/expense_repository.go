package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

// ExpenseRepository defines data access for expenses and their splits.
type ExpenseRepository interface {
	Create(ctx context.Context, expense *models.Expense, splits []*models.ExpenseSplit) error
	GetByID(ctx context.Context, id string) (*models.Expense, error)
	ListByFamily(ctx context.Context, familyID string) ([]*models.Expense, error)
	Update(ctx context.Context, expense *models.Expense) error
	Delete(ctx context.Context, id string) error

	ListSplitsByExpense(ctx context.Context, expenseID string) ([]*models.ExpenseSplit, error)
	SetSplitSettled(ctx context.Context, splitID string, settled bool) error

	// GetBalances computes each member's running net balance for a family.
	GetBalances(ctx context.Context, familyID string) ([]*models.Balance, error)
}

// postgresExpenseRepository is the Postgres-backed implementation of ExpenseRepository.
type postgresExpenseRepository struct {
	db *pgxpool.Pool
}

// NewExpenseRepository constructs a Postgres-backed ExpenseRepository.
func NewExpenseRepository(db *pgxpool.Pool) ExpenseRepository {
	return &postgresExpenseRepository{db: db}
}

func (r *postgresExpenseRepository) Create(ctx context.Context, expense *models.Expense, splits []*models.ExpenseSplit) error {
	return errors.New("not implemented")
}

func (r *postgresExpenseRepository) GetByID(ctx context.Context, id string) (*models.Expense, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresExpenseRepository) ListByFamily(ctx context.Context, familyID string) ([]*models.Expense, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresExpenseRepository) Update(ctx context.Context, expense *models.Expense) error {
	return errors.New("not implemented")
}

func (r *postgresExpenseRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (r *postgresExpenseRepository) ListSplitsByExpense(ctx context.Context, expenseID string) ([]*models.ExpenseSplit, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresExpenseRepository) SetSplitSettled(ctx context.Context, splitID string, settled bool) error {
	return errors.New("not implemented")
}

func (r *postgresExpenseRepository) GetBalances(ctx context.Context, familyID string) ([]*models.Balance, error) {
	return nil, errors.New("not implemented")
}
