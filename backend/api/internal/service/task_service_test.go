package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sibling-core/api/internal/mailer"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
	mockmailer "sibling-core/api/test/mock/mock_mailer"
	mockrepository "sibling-core/api/test/mock/mock_repository"
)

func newTestTaskService(t *testing.T) (service.TaskService, *mockrepository.MockTaskRepository, *mockrepository.MockMemberRepository, *mockrepository.MockParentRepository) {
	t.Helper()
	taskRepo := &mockrepository.MockTaskRepository{}
	memberRepo := &mockrepository.MockMemberRepository{}
	parentRepo := &mockrepository.MockParentRepository{}
	activityRepo := &mockrepository.MockActivityRepository{}
	return service.NewTaskService(taskRepo, memberRepo, parentRepo, activityRepo, &mockrepository.MockNotificationPreferencesRepository{}, &mockmailer.MockMailer{}), taskRepo, memberRepo, parentRepo
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

func TestCreateTaskWithDueDate(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	var created *models.Task
	taskRepo.CreateFunc = func(ctx context.Context, task *models.Task) error {
		created = task
		return nil
	}

	dueDate := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	_, err := svc.CreateTask(context.Background(), models.CreateTaskInput{
		FamilyID:  "family-1",
		Title:     "Refill prescriptions",
		CreatedBy: "user-1",
		DueDate:   &dueDate,
	})

	require.NoError(t, err)
	require.NotNil(t, created.DueDate)
	assert.True(t, dueDate.Equal(*created.DueDate))
}

func TestCreateTaskLogsActivity(t *testing.T) {
	taskRepo := &mockrepository.MockTaskRepository{}
	activityRepo := &mockrepository.MockActivityRepository{}
	svc := service.NewTaskService(taskRepo, &mockrepository.MockMemberRepository{}, &mockrepository.MockParentRepository{}, activityRepo, &mockrepository.MockNotificationPreferencesRepository{}, &mockmailer.MockMailer{})

	taskRepo.CreateFunc = func(ctx context.Context, task *models.Task) error {
		task.ID = "task-1"
		return nil
	}
	var logged *models.ActivityEntry
	activityRepo.CreateFunc = func(ctx context.Context, entry *models.ActivityEntry) error {
		logged = entry
		return nil
	}

	_, err := svc.CreateTask(context.Background(), models.CreateTaskInput{
		FamilyID: "family-1", Title: "Refill prescriptions", CreatedBy: "member-1", ActorName: "Jamie",
	})

	require.NoError(t, err)
	require.NotNil(t, logged)
	assert.Equal(t, "family-1", logged.FamilyID)
	assert.Equal(t, "task.created", logged.Verb)
	assert.Equal(t, "task", logged.SubjectType)
	require.NotNil(t, logged.SubjectID)
	assert.Equal(t, "task-1", *logged.SubjectID)
	assert.Equal(t, "Refill prescriptions", logged.SubjectLabel)
	assert.Equal(t, "Jamie", logged.ActorName)
	require.NotNil(t, logged.ActorID)
	assert.Equal(t, "member-1", *logged.ActorID)
}

func TestUpdateTaskCompletingLogsActivity(t *testing.T) {
	taskRepo := &mockrepository.MockTaskRepository{}
	activityRepo := &mockrepository.MockActivityRepository{}
	svc := service.NewTaskService(taskRepo, &mockrepository.MockMemberRepository{}, &mockrepository.MockParentRepository{}, activityRepo, &mockrepository.MockNotificationPreferencesRepository{}, &mockmailer.MockMailer{})

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1", Title: "Doctor visit", Status: "todo"}, nil
	}
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error { return nil }
	var logged *models.ActivityEntry
	activityRepo.CreateFunc = func(ctx context.Context, entry *models.ActivityEntry) error {
		logged = entry
		return nil
	}

	done := "done"
	_, err := svc.UpdateTask(context.Background(), "task-1", models.UpdateTaskInput{
		Status: &done, ActorID: "member-2", ActorName: "Sarah",
	})

	require.NoError(t, err)
	require.NotNil(t, logged)
	assert.Equal(t, "task.completed", logged.Verb)
	assert.Equal(t, "Doctor visit", logged.SubjectLabel)
	assert.Equal(t, "Sarah", logged.ActorName)
}

func TestUpdateTaskNotCompletingDoesNotLogActivity(t *testing.T) {
	taskRepo := &mockrepository.MockTaskRepository{}
	activityRepo := &mockrepository.MockActivityRepository{}
	svc := service.NewTaskService(taskRepo, &mockrepository.MockMemberRepository{}, &mockrepository.MockParentRepository{}, activityRepo, &mockrepository.MockNotificationPreferencesRepository{}, &mockmailer.MockMailer{})

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1", Title: "Doctor visit", Status: "todo"}, nil
	}
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error { return nil }
	logged := false
	activityRepo.CreateFunc = func(ctx context.Context, entry *models.ActivityEntry) error {
		logged = true
		return nil
	}

	inProgress := "in_progress"
	_, err := svc.UpdateTask(context.Background(), "task-1", models.UpdateTaskInput{Status: &inProgress})

	require.NoError(t, err)
	assert.False(t, logged, "only a transition to done should log activity")
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

func TestUpdateTaskDueDate(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1", Title: "Doctor visit", Status: "todo"}, nil
	}
	var updated *models.Task
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error {
		updated = task
		return nil
	}

	dueDate := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	task, err := svc.UpdateTask(context.Background(), "task-1", models.UpdateTaskInput{DueDate: &dueDate})

	require.NoError(t, err)
	require.NotNil(t, updated)
	require.NotNil(t, task.DueDate)
	assert.True(t, dueDate.Equal(*task.DueDate))
}

func TestCreateTaskDefaultsRecurrenceRuleToNone(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	var created *models.Task
	taskRepo.CreateFunc = func(ctx context.Context, task *models.Task) error {
		created = task
		return nil
	}

	_, err := svc.CreateTask(context.Background(), models.CreateTaskInput{FamilyID: "family-1", Title: "Doctor visit"})

	require.NoError(t, err)
	assert.Equal(t, "none", created.RecurrenceRule)
}

func TestCreateTaskRejectsInvalidRecurrenceRule(t *testing.T) {
	svc, _, _, _ := newTestTaskService(t)

	_, err := svc.CreateTask(context.Background(), models.CreateTaskInput{
		FamilyID:       "family-1",
		Title:          "Doctor visit",
		RecurrenceRule: "fortnightly",
	})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestUpdateTaskCompletingRecurringTaskSpawnsNextOccurrence(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	dueDate := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{
			ID: id, FamilyID: "family-1", Title: "Refill prescription", Status: "todo",
			RecurrenceRule: "weekly", DueDate: &dueDate,
		}, nil
	}
	var updated *models.Task
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error {
		updated = task
		return nil
	}
	var created *models.Task
	taskRepo.CreateFunc = func(ctx context.Context, task *models.Task) error {
		created = task
		return nil
	}

	done := "done"
	task, err := svc.UpdateTask(context.Background(), "task-1", models.UpdateTaskInput{Status: &done})

	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, "done", task.Status, "the completed task itself is marked done, not replaced")

	require.NotNil(t, created, "a next occurrence should have been created")
	assert.Equal(t, "Refill prescription", created.Title)
	assert.Equal(t, "weekly", created.RecurrenceRule)
	require.NotNil(t, created.DueDate)
	assert.True(t, dueDate.AddDate(0, 0, 7).Equal(*created.DueDate))
}

func TestUpdateTaskCompletingNonRecurringTaskDoesNotSpawnNextOccurrence(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1", Title: "One-off errand", Status: "todo", RecurrenceRule: "none"}, nil
	}
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error { return nil }
	created := false
	taskRepo.CreateFunc = func(ctx context.Context, task *models.Task) error {
		created = true
		return nil
	}

	done := "done"
	_, err := svc.UpdateTask(context.Background(), "task-1", models.UpdateTaskInput{Status: &done})

	require.NoError(t, err)
	assert.False(t, created)
}

func TestUpdateTaskRecompletingAlreadyDoneTaskDoesNotSpawnDuplicate(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	dueDate := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{
			ID: id, FamilyID: "family-1", Title: "Refill prescription", Status: "done",
			RecurrenceRule: "weekly", DueDate: &dueDate,
		}, nil
	}
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error { return nil }
	created := false
	taskRepo.CreateFunc = func(ctx context.Context, task *models.Task) error {
		created = true
		return nil
	}

	done := "done"
	_, err := svc.UpdateTask(context.Background(), "task-1", models.UpdateTaskInput{Status: &done})

	require.NoError(t, err)
	assert.False(t, created, "already-done -> done should not re-spawn a next occurrence")
}

func TestUpdateTaskRejectsInvalidRecurrenceRule(t *testing.T) {
	svc, taskRepo, _, _ := newTestTaskService(t)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1", Title: "Doctor visit", Status: "todo"}, nil
	}

	bad := "fortnightly"
	_, err := svc.UpdateTask(context.Background(), "task-1", models.UpdateTaskInput{RecurrenceRule: &bad})

	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestAssignTaskRejectsMemberFromAnotherFamily(t *testing.T) {
	svc, taskRepo, memberRepo, _ := newTestTaskService(t)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1"}, nil
	}
	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "other-family"}, nil
	}

	_, err := svc.AssignTask(context.Background(), "task-1", "member-1", "actor-1", "Actor")

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

	task, err := svc.AssignTask(context.Background(), "task-1", "member-1", "actor-1", "Actor")

	require.NoError(t, err)
	require.NotNil(t, updated)
	require.NotNil(t, task.AssigneeID)
	assert.Equal(t, "member-1", *task.AssigneeID)
}

func TestAssignTaskSendsEmailOnChange(t *testing.T) {
	taskRepo := &mockrepository.MockTaskRepository{}
	memberRepo := &mockrepository.MockMemberRepository{}
	notificationPrefsRepo := &mockrepository.MockNotificationPreferencesRepository{}
	mail := &mockmailer.MockMailer{}
	svc := service.NewTaskService(taskRepo, memberRepo, &mockrepository.MockParentRepository{}, &mockrepository.MockActivityRepository{}, notificationPrefsRepo, mail)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1", Title: "Doctor visit"}, nil
	}
	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "family-1", UserID: "user-2", Email: "sarah@example.com", FullName: "Sarah"}, nil
	}
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error { return nil }
	var sent *mailer.TaskAssignedEmail
	mail.SendTaskAssignedFunc = func(ctx context.Context, email mailer.TaskAssignedEmail) error {
		sent = &email
		return nil
	}

	_, err := svc.AssignTask(context.Background(), "task-1", "member-2", "member-1", "Jamie")

	require.NoError(t, err)
	require.NotNil(t, sent)
	assert.Equal(t, "sarah@example.com", sent.To)
	assert.Equal(t, "Doctor visit", sent.TaskTitle)
	assert.Equal(t, "Jamie", sent.AssignedBy)
}

func TestAssignTaskDoesNotEmailSelfAssignment(t *testing.T) {
	taskRepo := &mockrepository.MockTaskRepository{}
	memberRepo := &mockrepository.MockMemberRepository{}
	mail := &mockmailer.MockMailer{}
	svc := service.NewTaskService(taskRepo, memberRepo, &mockrepository.MockParentRepository{}, &mockrepository.MockActivityRepository{}, &mockrepository.MockNotificationPreferencesRepository{}, mail)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1", Title: "Doctor visit"}, nil
	}
	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "family-1"}, nil
	}
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error { return nil }
	sent := false
	mail.SendTaskAssignedFunc = func(ctx context.Context, email mailer.TaskAssignedEmail) error {
		sent = true
		return nil
	}

	_, err := svc.AssignTask(context.Background(), "task-1", "member-1", "member-1", "Jamie")

	require.NoError(t, err)
	assert.False(t, sent, "assigning a task to yourself shouldn't send yourself an email")
}

func TestAssignTaskRespectsOptOut(t *testing.T) {
	taskRepo := &mockrepository.MockTaskRepository{}
	memberRepo := &mockrepository.MockMemberRepository{}
	notificationPrefsRepo := &mockrepository.MockNotificationPreferencesRepository{}
	mail := &mockmailer.MockMailer{}
	svc := service.NewTaskService(taskRepo, memberRepo, &mockrepository.MockParentRepository{}, &mockrepository.MockActivityRepository{}, notificationPrefsRepo, mail)

	taskRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Task, error) {
		return &models.Task{ID: id, FamilyID: "family-1", Title: "Doctor visit"}, nil
	}
	memberRepo.GetByIDFunc = func(ctx context.Context, id string) (*models.Member, error) {
		return &models.Member{ID: id, FamilyID: "family-1", UserID: "user-2", Email: "sarah@example.com"}, nil
	}
	taskRepo.UpdateFunc = func(ctx context.Context, task *models.Task) error { return nil }
	notificationPrefsRepo.GetByUserIDFunc = func(ctx context.Context, userID string) (*models.NotificationPreferences, error) {
		return &models.NotificationPreferences{UserID: userID, TaskAssignedEmail: false}, nil
	}
	sent := false
	mail.SendTaskAssignedFunc = func(ctx context.Context, email mailer.TaskAssignedEmail) error {
		sent = true
		return nil
	}

	_, err := svc.AssignTask(context.Background(), "task-1", "member-2", "member-1", "Jamie")

	require.NoError(t, err)
	assert.False(t, sent, "opted-out assignee shouldn't get an email")
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
