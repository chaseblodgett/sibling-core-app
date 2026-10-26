package service

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// ExpenseService contains business logic for expense tracking, splitting, and balances.
type ExpenseService interface {
	CreateExpense(ctx context.Context, input models.CreateExpenseInput) (*models.Expense, error)
	GetExpense(ctx context.Context, id string) (*models.Expense, error)
	// ListExpenses returns a family's expenses, optionally filtered to a parent.
	ListExpenses(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error)
	DeleteExpense(ctx context.Context, id string) error

	ListSplits(ctx context.Context, expenseID string) ([]*models.ExpenseSplit, error)
	// GetSplit looks up a single split by id — used to resolve its family
	// (via its expense) for authorization before SettleSplit mutates it.
	GetSplit(ctx context.Context, id string) (*models.ExpenseSplit, error)
	SettleSplit(ctx context.Context, splitID string) (*models.ExpenseSplit, error)

	// GetBalances returns each member's running net balance for a family,
	// derived from unsettled splits — positive means owed to them.
	GetBalances(ctx context.Context, familyID string) ([]*models.Balance, error)
}

type expenseService struct {
	expenseRepo repository.ExpenseRepository
	memberRepo  repository.MemberRepository
	parentRepo  repository.ParentRepository
}

// NewExpenseService constructs an ExpenseService.
func NewExpenseService(expenseRepo repository.ExpenseRepository, memberRepo repository.MemberRepository, parentRepo repository.ParentRepository) ExpenseService {
	return &expenseService{expenseRepo: expenseRepo, memberRepo: memberRepo, parentRepo: parentRepo}
}

func (s *expenseService) CreateExpense(ctx context.Context, input models.CreateExpenseInput) (*models.Expense, error) {
	description := strings.TrimSpace(input.Description)
	if description == "" {
		return nil, fmt.Errorf("%w: description is required", ErrInvalidInput)
	}
	if input.AmountCents <= 0 {
		return nil, fmt.Errorf("%w: amount must be greater than zero", ErrInvalidInput)
	}

	if err := requireMemberInFamily(ctx, s.memberRepo, input.PaidBy, input.FamilyID, "paid_by"); err != nil {
		return nil, err
	}
	if input.ParentID != nil {
		if err := requireParentInFamily(ctx, s.parentRepo, *input.ParentID, input.FamilyID); err != nil {
			return nil, err
		}
	}

	splits := make([]*models.ExpenseSplit, 0, len(input.Splits))
	if len(input.Splits) > 0 {
		seen := make(map[string]bool, len(input.Splits))
		var sum int64
		for _, split := range input.Splits {
			if seen[split.MemberID] {
				return nil, fmt.Errorf("%w: a member can only appear once in splits", ErrInvalidInput)
			}
			seen[split.MemberID] = true

			if split.ShareCents <= 0 {
				return nil, fmt.Errorf("%w: each split must be greater than zero", ErrInvalidInput)
			}
			if err := requireMemberInFamily(ctx, s.memberRepo, split.MemberID, input.FamilyID, "split member"); err != nil {
				return nil, err
			}
			sum += split.ShareCents
			splits = append(splits, &models.ExpenseSplit{MemberID: split.MemberID, ShareCents: split.ShareCents})
		}
		if sum != input.AmountCents {
			return nil, fmt.Errorf("%w: splits must add up to the total amount", ErrInvalidInput)
		}
	}

	currency := strings.TrimSpace(input.Currency)
	if currency == "" {
		currency = "USD"
	}

	expense := &models.Expense{
		FamilyID:    input.FamilyID,
		PaidBy:      input.PaidBy,
		Description: description,
		AmountCents: input.AmountCents,
		Currency:    currency,
		Category:    strings.TrimSpace(input.Category),
		ParentID:    input.ParentID,
	}
	if err := s.expenseRepo.Create(ctx, expense, splits); err != nil {
		return nil, err
	}
	return expense, nil
}

func (s *expenseService) GetExpense(ctx context.Context, id string) (*models.Expense, error) {
	return s.expenseRepo.GetByID(ctx, id)
}

func (s *expenseService) ListExpenses(ctx context.Context, familyID string, parentID *string) ([]*models.Expense, error) {
	return s.expenseRepo.ListByFamily(ctx, familyID, parentID)
}

func (s *expenseService) DeleteExpense(ctx context.Context, id string) error {
	return s.expenseRepo.Delete(ctx, id)
}

func (s *expenseService) ListSplits(ctx context.Context, expenseID string) ([]*models.ExpenseSplit, error) {
	return s.expenseRepo.ListSplitsByExpense(ctx, expenseID)
}

func (s *expenseService) GetSplit(ctx context.Context, id string) (*models.ExpenseSplit, error) {
	return s.expenseRepo.GetSplitByID(ctx, id)
}

func (s *expenseService) SettleSplit(ctx context.Context, splitID string) (*models.ExpenseSplit, error) {
	split, err := s.expenseRepo.GetSplitByID(ctx, splitID)
	if err != nil {
		return nil, err
	}
	if err := s.expenseRepo.SetSplitSettled(ctx, splitID, true); err != nil {
		return nil, err
	}
	split.Settled = true
	return split, nil
}

func (s *expenseService) GetBalances(ctx context.Context, familyID string) ([]*models.Balance, error) {
	expenses, err := s.expenseRepo.ListByFamily(ctx, familyID, nil)
	if err != nil {
		return nil, err
	}
	paidByExpense := make(map[string]string, len(expenses))
	for _, e := range expenses {
		paidByExpense[e.ID] = e.PaidBy
	}

	splits, err := s.expenseRepo.ListSplitsByFamily(ctx, familyID)
	if err != nil {
		return nil, err
	}

	net := make(map[string]int64)
	for _, split := range splits {
		if split.Settled {
			continue
		}
		payer, ok := paidByExpense[split.ExpenseID]
		if !ok || payer == split.MemberID {
			// Unknown expense (shouldn't happen) or the payer's own
			// share of their own expense — nets to zero either way.
			continue
		}
		net[split.MemberID] -= split.ShareCents
		net[payer] += split.ShareCents
	}

	balances := make([]*models.Balance, 0, len(net))
	for memberID, n := range net {
		if n == 0 {
			continue
		}
		balances = append(balances, &models.Balance{MemberID: memberID, NetCents: n})
	}
	sort.Slice(balances, func(i, j int) bool { return balances[i].MemberID < balances[j].MemberID })
	return balances, nil
}
