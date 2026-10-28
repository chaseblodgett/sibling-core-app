package service

import (
	"context"
	"sort"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// activityScanLimit bounds the single ListByFamily call ContributionService
// makes to tally task.completed entries. There's no pagination here (unlike
// ActivityService.ListActivity) because a contribution total needs every
// entry ever logged, not a recent page of them — 10,000 is a practical
// "effectively unlimited" bound for a small-family-scale app, not a real
// page size.
const activityScanLimit = 10000

// ContributionService computes, read-only, how much caregiving work each
// family member has picked up: tasks they've actually completed (not just
// been assigned) and expenses they've actually fronted (not just logged
// on someone else's behalf). No new write-side model — it's an
// aggregation over Member/Task-completion activity/Expense data that
// already exists, the same way ExpenseService.GetBalances aggregates
// over Expense/ExpenseSplit without its own table.
type ContributionService interface {
	GetContributions(ctx context.Context, familyID string) ([]*models.Contribution, error)
}

type contributionService struct {
	memberRepo   repository.MemberRepository
	expenseRepo  repository.ExpenseRepository
	activityRepo repository.ActivityRepository
}

// NewContributionService constructs a ContributionService.
func NewContributionService(memberRepo repository.MemberRepository, expenseRepo repository.ExpenseRepository, activityRepo repository.ActivityRepository) ContributionService {
	return &contributionService{memberRepo: memberRepo, expenseRepo: expenseRepo, activityRepo: activityRepo}
}

func (s *contributionService) GetContributions(ctx context.Context, familyID string) ([]*models.Contribution, error) {
	members, err := s.memberRepo.ListByFamily(ctx, familyID)
	if err != nil {
		return nil, err
	}

	tasksCompleted, err := s.tasksCompletedByMember(ctx, familyID)
	if err != nil {
		return nil, err
	}

	expenseCounts, expenseCents, err := s.expensesPaidByMember(ctx, familyID)
	if err != nil {
		return nil, err
	}

	contributions := make([]*models.Contribution, 0, len(members))
	for _, m := range members {
		name := m.FullName
		if name == "" {
			name = m.Email
		}
		contributions = append(contributions, &models.Contribution{
			MemberID:          m.ID,
			MemberName:        name,
			TasksCompleted:    tasksCompleted[m.ID],
			ExpensesPaidCount: expenseCounts[m.ID],
			ExpensesPaidCents: expenseCents[m.ID],
		})
	}

	// Everyone shows up, including members with nothing logged yet — the
	// point of a contribution view is surfacing imbalance, which a list
	// that hides low/zero contributors would defeat. Ranked by tasks
	// completed first (the plan's primary contribution signal), then by
	// amount fronted, so the most active member leads.
	sort.Slice(contributions, func(i, j int) bool {
		if contributions[i].TasksCompleted != contributions[j].TasksCompleted {
			return contributions[i].TasksCompleted > contributions[j].TasksCompleted
		}
		if contributions[i].ExpensesPaidCents != contributions[j].ExpensesPaidCents {
			return contributions[i].ExpensesPaidCents > contributions[j].ExpensesPaidCents
		}
		return contributions[i].MemberName < contributions[j].MemberName
	})
	return contributions, nil
}

// tasksCompletedByMember counts task.completed activity entries per actor.
// This is the only accurate source for "who completed it" — the tasks
// table itself has no completed_by column, only created_by/assignee_id,
// neither of which reliably says who actually marked it done.
func (s *contributionService) tasksCompletedByMember(ctx context.Context, familyID string) (map[string]int, error) {
	entries, err := s.activityRepo.ListByFamily(ctx, familyID, activityScanLimit, 0)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, e := range entries {
		if e.Verb != "task.completed" || e.ActorID == nil {
			continue
		}
		counts[*e.ActorID]++
	}
	return counts, nil
}

// expensesPaidByMember counts/sums expenses by who actually paid_by —
// not who logged them (activity_log's expense.created actor is the
// requester, which can differ from the payer when logging a sibling's
// receipt on their behalf; attribution for a fairness-tracking feature
// needs to be the payer).
func (s *contributionService) expensesPaidByMember(ctx context.Context, familyID string) (counts map[string]int, cents map[string]int64, err error) {
	expenses, err := s.expenseRepo.ListByFamily(ctx, familyID, nil)
	if err != nil {
		return nil, nil, err
	}
	counts = make(map[string]int)
	cents = make(map[string]int64)
	for _, e := range expenses {
		counts[e.PaidBy]++
		cents[e.PaidBy] += e.AmountCents
	}
	return counts, cents, nil
}
