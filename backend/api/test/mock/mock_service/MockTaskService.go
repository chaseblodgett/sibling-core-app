package mockservice

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockTaskService is a test double for service.TaskService. Assign the
// *Func fields your test exercises; unset ones return zero values.
type MockTaskService struct {
	CreateTaskFunc func(ctx context.Context, input models.CreateTaskInput) (*models.Task, error)
	GetTaskFunc    func(ctx context.Context, id string) (*models.Task, error)
	ListTasksFunc  func(ctx context.Context, familyID string) ([]*models.Task, error)
	UpdateTaskFunc func(ctx context.Context, id string, input models.UpdateTaskInput) (*models.Task, error)
	AssignTaskFunc func(ctx context.Context, id, memberID string) (*models.Task, error)
	DeleteTaskFunc func(ctx context.Context, id string) error
}

func (m *MockTaskService) CreateTask(ctx context.Context, input models.CreateTaskInput) (*models.Task, error) {
	if m.CreateTaskFunc == nil {
		return nil, nil
	}
	return m.CreateTaskFunc(ctx, input)
}

func (m *MockTaskService) GetTask(ctx context.Context, id string) (*models.Task, error) {
	if m.GetTaskFunc == nil {
		return nil, nil
	}
	return m.GetTaskFunc(ctx, id)
}

func (m *MockTaskService) ListTasks(ctx context.Context, familyID string) ([]*models.Task, error) {
	if m.ListTasksFunc == nil {
		return nil, nil
	}
	return m.ListTasksFunc(ctx, familyID)
}

func (m *MockTaskService) UpdateTask(ctx context.Context, id string, input models.UpdateTaskInput) (*models.Task, error) {
	if m.UpdateTaskFunc == nil {
		return nil, nil
	}
	return m.UpdateTaskFunc(ctx, id, input)
}

func (m *MockTaskService) AssignTask(ctx context.Context, id, memberID string) (*models.Task, error) {
	if m.AssignTaskFunc == nil {
		return nil, nil
	}
	return m.AssignTaskFunc(ctx, id, memberID)
}

func (m *MockTaskService) DeleteTask(ctx context.Context, id string) error {
	if m.DeleteTaskFunc == nil {
		return nil
	}
	return m.DeleteTaskFunc(ctx, id)
}
