package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

var validTaskStatuses = map[string]bool{"todo": true, "in_progress": true, "done": true}

var validRecurrenceRules = map[string]bool{"none": true, "daily": true, "weekly": true, "monthly": true}

// nextDueDate advances due by one cycle of rule. Callers only invoke this
// for non-"none" rules, so the default case (returning due unchanged)
// never actually triggers in practice — it exists so the function is
// total rather than panicking on an unexpected value.
func nextDueDate(due time.Time, rule string) time.Time {
	switch rule {
	case "daily":
		return due.AddDate(0, 0, 1)
	case "weekly":
		return due.AddDate(0, 0, 7)
	case "monthly":
		return due.AddDate(0, 1, 0)
	default:
		return due
	}
}

// TaskService contains business logic for the shared family task board.
type TaskService interface {
	CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error)
	GetTask(ctx context.Context, id string) (*models.Task, error)
	// ListTasks returns a family's tasks, optionally filtered to a parent.
	ListTasks(ctx context.Context, familyID string, parentID *string) ([]*models.Task, error)
	UpdateTask(ctx context.Context, id string, input models.UpdateTaskInput) (*models.Task, error)
	AssignTask(ctx context.Context, id, memberID string) (*models.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type taskService struct {
	taskRepo   repository.TaskRepository
	memberRepo repository.MemberRepository
	parentRepo repository.ParentRepository
}

// NewTaskService constructs a TaskService.
func NewTaskService(taskRepo repository.TaskRepository, memberRepo repository.MemberRepository, parentRepo repository.ParentRepository) TaskService {
	return &taskService{taskRepo: taskRepo, memberRepo: memberRepo, parentRepo: parentRepo}
}

func (s *taskService) CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.AssigneeID != nil {
		if err := requireMemberInFamily(ctx, s.memberRepo, *input.AssigneeID, input.FamilyID, "assignee"); err != nil {
			return nil, err
		}
	}
	if input.ParentID != nil {
		if err := requireParentInFamily(ctx, s.parentRepo, *input.ParentID, input.FamilyID); err != nil {
			return nil, err
		}
	}

	recurrenceRule := input.RecurrenceRule
	if recurrenceRule == "" {
		recurrenceRule = "none"
	} else if !validRecurrenceRules[recurrenceRule] {
		return nil, fmt.Errorf("%w: recurrence_rule must be none, daily, weekly, or monthly", ErrInvalidInput)
	}

	task := &models.Task{
		FamilyID:       input.FamilyID,
		Title:          title,
		Description:    input.Description,
		AssigneeID:     input.AssigneeID,
		ParentID:       input.ParentID,
		CreatedBy:      input.CreatedBy,
		DueDate:        input.DueDate,
		RecurrenceRule: recurrenceRule,
	}
	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetTask(ctx context.Context, id string) (*models.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *taskService) ListTasks(ctx context.Context, familyID string, parentID *string) ([]*models.Task, error) {
	return s.taskRepo.ListByFamily(ctx, familyID, parentID)
}

func (s *taskService) UpdateTask(ctx context.Context, id string, input models.UpdateTaskInput) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	wasDone := task.Status == "done"

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return nil, fmt.Errorf("%w: title is required", ErrInvalidInput)
		}
		task.Title = title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Status != nil {
		if !validTaskStatuses[*input.Status] {
			return nil, fmt.Errorf("%w: status must be todo, in_progress, or done", ErrInvalidInput)
		}
		task.Status = *input.Status
	}
	if input.AssigneeID != nil {
		if err := requireMemberInFamily(ctx, s.memberRepo, *input.AssigneeID, task.FamilyID, "assignee"); err != nil {
			return nil, err
		}
		task.AssigneeID = input.AssigneeID
	}
	if input.ParentID != nil {
		if err := requireParentInFamily(ctx, s.parentRepo, *input.ParentID, task.FamilyID); err != nil {
			return nil, err
		}
		task.ParentID = input.ParentID
	}
	if input.DueDate != nil {
		task.DueDate = input.DueDate
	}
	if input.RecurrenceRule != nil {
		if !validRecurrenceRules[*input.RecurrenceRule] {
			return nil, fmt.Errorf("%w: recurrence_rule must be none, daily, weekly, or monthly", ErrInvalidInput)
		}
		task.RecurrenceRule = *input.RecurrenceRule
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	// Completing a recurring task spawns its next occurrence rather than
	// just closing it out — that's the whole point of marking a task
	// recurring. Only on a true not-done -> done transition, so replaying
	// the same PATCH (or editing an already-done task) doesn't spawn
	// duplicates.
	if !wasDone && task.Status == "done" && task.RecurrenceRule != "none" && task.DueDate != nil {
		nextDue := nextDueDate(*task.DueDate, task.RecurrenceRule)
		next := &models.Task{
			FamilyID:       task.FamilyID,
			Title:          task.Title,
			Description:    task.Description,
			AssigneeID:     task.AssigneeID,
			ParentID:       task.ParentID,
			CreatedBy:      task.CreatedBy,
			DueDate:        &nextDue,
			RecurrenceRule: task.RecurrenceRule,
		}
		if err := s.taskRepo.Create(ctx, next); err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (s *taskService) AssignTask(ctx context.Context, id, memberID string) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := requireMemberInFamily(ctx, s.memberRepo, memberID, task.FamilyID, "assignee"); err != nil {
		return nil, err
	}

	task.AssigneeID = &memberID
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	return s.taskRepo.Delete(ctx, id)
}
