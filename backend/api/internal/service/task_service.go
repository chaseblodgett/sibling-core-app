package service

import (
	"context"
	"fmt"
	"strings"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

var validTaskStatuses = map[string]bool{"todo": true, "in_progress": true, "done": true}

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

	task := &models.Task{
		FamilyID:    input.FamilyID,
		Title:       title,
		Description: input.Description,
		AssigneeID:  input.AssigneeID,
		ParentID:    input.ParentID,
		CreatedBy:   input.CreatedBy,
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

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
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
