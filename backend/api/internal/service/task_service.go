package service

import (
	"context"
	"errors"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// TaskService contains business logic for the shared family task board.
type TaskService interface {
	CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error)
	GetTask(ctx context.Context, id string) (*models.Task, error)
	ListTasks(ctx context.Context, familyID string) ([]*models.Task, error)
	UpdateTask(ctx context.Context, id string, input models.UpdateTaskInput) (*models.Task, error)
	AssignTask(ctx context.Context, id, memberID string) (*models.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type taskService struct {
	taskRepo repository.TaskRepository
}

// NewTaskService constructs a TaskService.
func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (s *taskService) GetTask(ctx context.Context, id string) (*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (s *taskService) ListTasks(ctx context.Context, familyID string) ([]*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (s *taskService) UpdateTask(ctx context.Context, id string, input models.UpdateTaskInput) (*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (s *taskService) AssignTask(ctx context.Context, id, memberID string) (*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
