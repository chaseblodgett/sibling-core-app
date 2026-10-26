package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
	mockrepository "sibling-core/api/test/mock/mock_repository"
)

func newTestTaskService(t *testing.T) (service.TaskService, *mockrepository.MockTaskRepository, *mockrepository.MockMemberRepository, *mockrepository.MockParentRepository) {
	t.Helper()
	taskRepo := &mockrepository.MockTaskRepository{}
	memberRepo := &mockrepository.MockMemberRepository{}
	parentRepo := &mockrepository.MockParentRepository{}
	return service.NewTaskService(taskRepo, memberRepo, parentRepo), taskRepo, memberRepo, parentRepo
}

func TestCreateTaskSuccess(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	var created *models.Task
	taskRepo.CreateFunc = func(ctx context.Context, task *models.Task) error {
		task.ID = "task-1"
		task.Status = "todo"
		created = task
		return nil
	}

	task, err := svc.CreateTask(context.Background(), models.CreateTaskInput{
		FamilyID:  "family-1",
		Title:     "  Refill prescriptions ",
		CreatedBy: "user-1",
	})

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, "Refill prescriptions", task.Title)
	assert.Equal(t, "todo", task.Status)
}

func TestCreateTaskBlankTitle(t *testing.T) {
	svc, _, _, _ := newTestTaskService(t)

	_, err := svc.CreateTask(context.Background(), models.CreateTaskInput{FamilyID: "family-1", Title: "  "})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestCreateTaskWithParentValidatesFamily(t *testing.T) {
	svc, _, _, parentRepo := newTestTaskService(t)

	parentRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Parent, error) {
		return &models.Parent{ID: id, FamilyID: "other-family"}, nil
	}

	parentID := "parent-1"
	_, err := svc.CreateTask(context.Background(), models.CreateTaskInput{
		FamilyID: "family-1",
		Title:    "Doctor visit",
		ParentID: &parentID,
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestCreateTaskWithParentInSameFamilySucceeds(t *testing.T) {
	svc, taskRepo, _, parentRepo := newTestTaskService(t)

	parentRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Parent, error) {
		return &models.Parent{ID: id, FamilyID: "family-1"}, nil
	}
	var created *models.Task
	taskRepo.CreateFunc = func(ctx context.Context, task *models.Task) error {
		created = task
		return nil
	}

	parentID := "parent-1"
	_, err := svc.CreateTask(context.Background(), models.CreateTaskInput{
		FamilyID: "family-1",
		Title:    "Doctor visit",
		ParentID: &parentID,
	})

	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotNil(t, created.ParentID)
	assert.Equal(t, "parent-1", *created.ParentID)
}

func TestCreateTaskRejectsAssigneeFromAnotherFamily(t *testing.T) {
	svc, _, memberRepo, _ := newTestTaskService(t)

	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "other-family"}, nil
	}

	assigneeID := "member-1"
	_, err := svc.CreateTask(context.Background(), models.CreateTaskInput{
		FamilyID:   "family-1",
		Title:      "Doctor visit",
		AssigneeID: &assigneeID,
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestUpdateTaskValidatesStatus(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1", Title: "Doctor visit", Status: "todo"}, nil
	}

	badStatus := "done-ish"
	_, err := svc.UpdateTask(context.Background(), "task-1", models.UpdateTaskInput{Status: &badStatus})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestUpdateTaskPartial(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1", Title: "Doctor visit", Status: "todo"}, nil
	}
	var updated *models.Task
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error {
		updated = task
		return nil
	}

	done := "done"
	task, err := svc.UpdateTask(context.Background(), "task-1", models.UpdateTaskInput{Status: &done})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "Doctor visit", task.Title, "title should be unchanged when not provided")
	assert.Equal(t, "done", task.Status)
}

func TestAssignTaskRejectsMemberFromAnotherFamily(t *testing.T) {
	svc, taskRepo, memberRepo, _ := newTestTaskService(t)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1"}, nil
	}
	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "other-family"}, nil
	}

	_, err := svc.AssignTask(context.Background(), "task-1", "member-1")

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestAssignTaskSuccess(t *testing.T) {
	svc, taskRepo, memberRepo, _ := newTestTaskService(t)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1"}, nil
	}
	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "family-1"}, nil
	}
	var updated *models.Task
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error {
		updated = task
		return nil
	}

	task, err := svc.AssignTask(context.Background(), "task-1", "member-1")

	require.NoError(t, err)
	require.NotNil(t, updated)
	require.NotNil(t, task.AssigneeID)
	assert.Equal(t, "member-1", *task.AssigneeID)
}

func TestListTasksFiltersByParent(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	taskRepo.ListByFamilyFunc = func(ctx context.Context, familyID string, parentID *string) ([]*models.Task, error) {
		assert.Equal(t, "family-1", familyID)
		require.NotNil(t, parentID)
		assert.Equal(t, "parent-1", *parentID)
		return []*models.Task{{ID: "task-1"}}, nil
	}

	parentID := "parent-1"
	tasks, err := svc.ListTasks(context.Background(), "family-1", &parentID)

	require.NoError(t, err)
	assert.Len(t, tasks, 1)
}

func TestDeleteTaskNotFound(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	taskRepo.DeleteFunc = func(ctx context.Context, id string) error {
		return repository.ErrTaskNotFound
	}

	err := svc.DeleteTask(context.Background(), "task-1")

	assert.ErrorIs(t, err, repository.ErrTaskNotFound)
}
