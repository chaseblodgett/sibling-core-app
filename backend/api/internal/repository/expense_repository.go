package repository

import (
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

var (
	ErrExpenseNotFound = errors.New("expense not found")
	ErrSplitNotFound   = errors.New("expense split not found")
)

const (
	expensesTable      = "expenses"
	expenseSplitsTable = "expense_splits"
)

// ExpenseRepository defines data access for expenses and their splits.
type ExpenseRepository interface {
	// Create inserts the expense and all of its splits as one transaction
	// — either everything lands, or nothing does.
	Create(ctx context.Context, expense *models.Expense, splits []*models.ExpenseSplit) error
	GetByID(ctx context.Context, id string) (*models.Expense, error)
	// ListByFamily returns a family's expenses, optionally filtered to
	// those tied to a specific parent (nil parentID returns all of them).
	ListByFamily(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error)
	Delete(ctx context.Context, id string) error
	// SetDocumentID attaches (or, with a nil documentID, detaches) a
	// receipt document to an expense.
	SetDocumentID(ctx context.Context, expenseID string, documentID *string) error

	ListSplitsByExpense(ctx context.Context, expenseID string) ([]*models.ExpenseSplit, error)
	// ListSplitsByFamily returns every split across every expense in a
	// family — the raw data GetBalances aggregates in the service layer.
	ListSplitsByFamily(ctx context.Context, familyID string) ([]*models.ExpenseSplit, error)
	GetSplitByID(ctx context.Context, id string) (*models.ExpenseSplit, error)
	SetSplitSettled(ctx context.Context, splitID string, settled bool) error
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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck // no-op once committed

	expenseQuery, expenseArgs, err := dialect.Insert(expensesTable).
		Rows(goqu.Record{
			"family_id":    expense.FamilyID,
			"paid_by":      expense.PaidBy,
			"description":  expense.Description,
			"amount_cents": expense.AmountCents,
			"currency":     expense.Currency,
			"category":     expense.Category,
			"parent_id":    expense.ParentID,
			"incurred_at":  goqu.L("now()"),
		}).
		Returning("id", "incurred_at", "created_at", "updated_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}
	if err := tx.QueryRow(ctx, expenseQuery, expenseArgs...).Scan(
		&expense.ID, &expense.IncurredAt, &expense.CreatedAt, &expense.UpdatedAt,
	); err != nil {
		return err
	}

	for _, split := range splits {
		split.ExpenseID = expense.ID
		splitQuery, splitArgs, err := dialect.Insert(expenseSplitsTable).
			Rows(goqu.Record{
				"expense_id":  split.ExpenseID,
				"member_id":   split.MemberID,
				"share_cents": split.ShareCents,
			}).
			Returning("id", "settled").
			Prepared(true).
			ToSQL()
		if err != nil {
			return err
		}
		if err := tx.QueryRow(ctx, splitQuery, splitArgs...).Scan(&split.ID, &split.Settled); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *postgresExpenseRepository) GetByID(ctx context.Context, id string) (*models.Expense, error) {
	query, args, err := expenseSelectDataset().Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanExpense(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresExpenseRepository) ListByFamily(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error) {
	ds := expenseSelectDataset().Where(goqu.Ex{"family_id": familyID})
	if parentID != nil {
		ds = ds.Where(goqu.Ex{"parent_id": *parentID})
	}
	query, args, err := ds.Order(goqu.I("incurred_at").Desc()).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]*models.Expense, 0)
	for rows.Next() {
		var e models.Expense
		if err := scanExpenseRow(rows, &e); err != nil {
			return nil, err
		}
		expenses = append(expenses, &e)
	}
	return expenses, rows.Err()
}

func (r *postgresExpenseRepository) Delete(ctx context.Context, id string) error {
	query, args, err := dialect.Delete(expensesTable).Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrExpenseNotFound
	}
	return nil
}

func (r *postgresExpenseRepository) SetDocumentID(ctx context.Context, expenseID string, documentID *string) error {
	query, args, err := dialect.Update(expensesTable).
		Set(goqu.Record{"document_id": documentID, "updated_at": goqu.L("now()")}).
		Where(goqu.Ex{"id": expenseID}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrExpenseNotFound
	}
	return nil
}

func (r *postgresExpenseRepository) ListSplitsByExpense(ctx context.Context, expenseID string) ([]*models.ExpenseSplit, error) {
	query, args, err := splitSelectDataset().Where(goqu.Ex{"expense_id": expenseID}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return querySplits(ctx, r.db, query, args...)
}

func (r *postgresExpenseRepository) ListSplitsByFamily(ctx context.Context, familyID string) ([]*models.ExpenseSplit, error) {
	query, args, err := dialect.From(expenseSplitsTable).
		Select(
			goqu.I("expense_splits.id"),
			goqu.I("expense_splits.expense_id"),
			goqu.I("expense_splits.member_id"),
			goqu.I("expense_splits.share_cents"),
			goqu.I("expense_splits.settled"),
		).
		InnerJoin(goqu.I(expensesTable), goqu.On(goqu.I("expense_splits.expense_id").Eq(goqu.I("expenses.id")))).
		Where(goqu.I("expenses.family_id").Eq(familyID)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}
	return querySplits(ctx, r.db, query, args...)
}

func (r *postgresExpenseRepository) GetSplitByID(ctx context.Context, id string) (*models.ExpenseSplit, error) {
	query, args, err := splitSelectDataset().Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanSplit(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresExpenseRepository) SetSplitSettled(ctx context.Context, splitID string, settled bool) error {
	query, args, err := dialect.Update(expenseSplitsTable).
		Set(goqu.Record{"settled": settled}).
		Where(goqu.Ex{"id": splitID}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrSplitNotFound
	}
	return nil
}

func expenseSelectDataset() *goqu.SelectDataset {
	return dialect.From(expensesTable).
		Select("id", "family_id", "paid_by", "description", "amount_cents", "currency", "category", "parent_id", "document_id", "incurred_at", "created_at", "updated_at")
}

func scanExpenseRow(row taskScanner, e *models.Expense) error {
	return row.Scan(&e.ID, &e.FamilyID, &e.PaidBy, &e.Description, &e.AmountCents, &e.Currency, &e.Category, &e.ParentID, &e.DocumentID, &e.IncurredAt, &e.CreatedAt, &e.UpdatedAt)
}

func scanExpense(row pgx.Row) (*models.Expense, error) {
	var e models.Expense
	if err := scanExpenseRow(row, &e); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrExpenseNotFound
		}
		return nil, err
	}
	return &e, nil
}

func splitSelectDataset() *goqu.SelectDataset {
	return dialect.From(expenseSplitsTable).Select("id", "expense_id", "member_id", "share_cents", "settled")
}

func scanSplitRow(row taskScanner, s *models.ExpenseSplit) error {
	return row.Scan(&s.ID, &s.ExpenseID, &s.MemberID, &s.ShareCents, &s.Settled)
}

func scanSplit(row pgx.Row) (*models.ExpenseSplit, error) {
	var s models.ExpenseSplit
	if err := scanSplitRow(row, &s); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSplitNotFound
		}
		return nil, err
	}
	return &s, nil
}

func querySplits(ctx context.Context, db *pgxpool.Pool, query string, args ...any) ([]*models.ExpenseSplit, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	splits := make([]*models.ExpenseSplit, 0)
	for rows.Next() {
		var s models.ExpenseSplit
		if err := scanSplitRow(rows, &s); err != nil {
			return nil, err
		}
		splits = append(splits, &s)
	}
	return splits, rows.Err()
}
