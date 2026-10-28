package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"time"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// ExportService renders a family's expenses or tasks as CSV — a real
// caregiving need (compiling a medical expense history for a tax
// deduction, or a task/expense record for an elder-law attorney or a
// sibling dispute), not a nice-to-have. Read-only, like Contribution and
// Activity; no new write-side model.
type ExportService interface {
	// ExportExpenses renders a family's expenses as CSV, one row per
	// expense (not per split — a readable "who owes what" summary column
	// instead, so the file stays one-row-per-transaction). Optional
	// parentID/from/to narrow it; nil means unfiltered.
	ExportExpenses(ctx context.Context, familyID string, parentID *string, from, to *time.Time) ([]byte, error)
	// ExportTasks renders a family's tasks as CSV, one row per task.
	ExportTasks(ctx context.Context, familyID string, parentID *string, from, to *time.Time) ([]byte, error)
}

type exportService struct {
	expenseRepo repository.ExpenseRepository
	taskRepo    repository.TaskRepository
	memberRepo  repository.MemberRepository
	parentRepo  repository.ParentRepository
}

// NewExportService constructs an ExportService.
func NewExportService(expenseRepo repository.ExpenseRepository, taskRepo repository.TaskRepository, memberRepo repository.MemberRepository, parentRepo repository.ParentRepository) ExportService {
	return &exportService{expenseRepo: expenseRepo, taskRepo: taskRepo, memberRepo: memberRepo, parentRepo: parentRepo}
}

func inRange(t time.Time, from, to *time.Time) bool {
	if from != nil && t.Before(*from) {
		return false
	}
	if to != nil && t.After(*to) {
		return false
	}
	return true
}

func (s *exportService) nameMaps(ctx context.Context, familyID string) (members map[string]string, parents map[string]string, err error) {
	memberList, err := s.memberRepo.ListByFamily(ctx, familyID)
	if err != nil {
		return nil, nil, err
	}
	members = make(map[string]string, len(memberList))
	for _, m := range memberList {
		name := m.FullName
		if name == "" {
			name = m.Email
		}
		members[m.ID] = name
	}

	parentList, err := s.parentRepo.ListByFamily(ctx, familyID)
	if err != nil {
		return nil, nil, err
	}
	parents = make(map[string]string, len(parentList))
	for _, p := range parentList {
		parents[p.ID] = p.Name
	}
	return members, parents, nil
}

func (s *exportService) ExportExpenses(ctx context.Context, familyID string, parentID *string, from, to *time.Time) ([]byte, error) {
	expenses, err := s.expenseRepo.ListByFamily(ctx, familyID, parentID)
	if err != nil {
		return nil, err
	}
	splits, err := s.expenseRepo.ListSplitsByFamily(ctx, familyID)
	if err != nil {
		return nil, err
	}
	members, parents, err := s.nameMaps(ctx, familyID)
	if err != nil {
		return nil, err
	}

	splitsByExpense := make(map[string][]*models.ExpenseSplit, len(expenses))
	for _, split := range splits {
		splitsByExpense[split.ExpenseID] = append(splitsByExpense[split.ExpenseID], split)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"date", "description", "category", "amount", "currency", "paid_by", "parent", "splits"})
	for _, e := range expenses {
		if !inRange(e.IncurredAt, from, to) {
			continue
		}
		parentName := ""
		if e.ParentID != nil {
			parentName = parents[*e.ParentID]
		}
		splitsSummary := ""
		for i, split := range splitsByExpense[e.ID] {
			if i > 0 {
				splitsSummary += "; "
			}
			status := "unsettled"
			if split.Settled {
				status = "settled"
			}
			splitsSummary += fmt.Sprintf("%s: %s (%s)", members[split.MemberID], formatCentsCSV(split.ShareCents, e.Currency), status)
		}
		_ = w.Write([]string{
			e.IncurredAt.Format("2006-01-02"),
			e.Description,
			e.Category,
			formatCentsCSV(e.AmountCents, e.Currency),
			e.Currency,
			members[e.PaidBy],
			parentName,
			splitsSummary,
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *exportService) ExportTasks(ctx context.Context, familyID string, parentID *string, from, to *time.Time) ([]byte, error) {
	tasks, err := s.taskRepo.ListByFamily(ctx, familyID, parentID)
	if err != nil {
		return nil, err
	}
	members, parents, err := s.nameMaps(ctx, familyID)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"created", "title", "description", "status", "assignee", "parent", "due_date", "repeats"})
	for _, t := range tasks {
		if !inRange(t.CreatedAt, from, to) {
			continue
		}
		assignee := ""
		if t.AssigneeID != nil {
			assignee = members[*t.AssigneeID]
		}
		parentName := ""
		if t.ParentID != nil {
			parentName = parents[*t.ParentID]
		}
		dueDate := ""
		if t.DueDate != nil {
			dueDate = t.DueDate.Format("2006-01-02")
		}
		_ = w.Write([]string{
			t.CreatedAt.Format("2006-01-02"),
			t.Title,
			t.Description,
			t.Status,
			assignee,
			parentName,
			dueDate,
			t.RecurrenceRule,
		})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func formatCentsCSV(cents int64, currency string) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	return fmt.Sprintf("%s%s%d.%02d", sign, currencySymbol(currency), cents/100, cents%100)
}

func currencySymbol(currency string) string {
	if currency == "USD" || currency == "" {
		return "$"
	}
	return currency + " "
}
