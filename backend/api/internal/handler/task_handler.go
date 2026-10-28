package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/middleware"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
)

// TaskHandler exposes HTTP endpoints for the shared family task board.
// Every route requires the requester to belong to the task's family;
// there's no further role restriction — any member can manage any task.
type TaskHandler struct {
	taskService   service.TaskService
	familyService service.FamilyService
	authService   service.AuthService
}

// NewTaskHandler constructs a TaskHandler.
func NewTaskHandler(taskService service.TaskService, familyService service.FamilyService, authService service.AuthService) *TaskHandler {
	return &TaskHandler{taskService: taskService, familyService: familyService, authService: authService}
}

// RegisterRoutes wires task endpoints onto the given group, all behind
// middleware.RequireAuth.
func (h *TaskHandler) RegisterRoutes(g *echo.Group) {
	auth := middleware.RequireAuth(h.authService)
	g.POST("/families/:familyId/tasks", h.CreateTask, auth)
	g.GET("/families/:familyId/tasks", h.ListTasks, auth)
	g.GET("/tasks/:id", h.GetTask, auth)
	g.PATCH("/tasks/:id", h.UpdateTask, auth)
	g.POST("/tasks/:id/assign", h.AssignTask, auth)
	g.DELETE("/tasks/:id", h.DeleteTask, auth)
}

// CreateTask creates a new task on a family's board, optionally tied to a
// parent (?parent_id in the body).
func (h *TaskHandler) CreateTask(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	// tasks.created_by references members(id), not users(id) — unlike
	// families.created_by. Need the requester's membership row for this
	// family, not their bare user id.
	requester, err := h.familyService.RequireMembership(ctx, user.ID, familyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if err := requireNotViewer(requester); err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Message: err.Error()})
	}

	var req models.CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	task, err := h.taskService.CreateTask(ctx, models.CreateTaskInput{
		FamilyID:       familyID,
		Title:          req.Title,
		Description:    req.Description,
		AssigneeID:     req.AssigneeID,
		ParentID:       req.ParentID,
		CreatedBy:      requester.ID,
		DueDate:        req.DueDate,
		RecurrenceRule: req.RecurrenceRule,
		ActorName:      memberDisplayName(requester),
	})
	if err != nil {
		status, msg := mapTaskError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusCreated, task)
}

// ListTasks returns a family's tasks. An optional ?parent_id query param
// filters to tasks tied to that parent.
func (h *TaskHandler) ListTasks(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	var parentID *string
	if v := c.QueryParam("parent_id"); v != "" {
		parentID = &v
	}

	tasks, err := h.taskService.ListTasks(ctx, familyID, parentID)
	if err != nil {
		status, msg := mapTaskError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, tasks)
}

// GetTask returns a single task by id, if the requester belongs to its family.
func (h *TaskHandler) GetTask(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	task, err := h.taskService.GetTask(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapTaskError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if _, err := h.familyService.RequireMembership(ctx, user.ID, task.FamilyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, task)
}

// UpdateTask updates fields on an existing task.
func (h *TaskHandler) UpdateTask(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	task, err := h.taskService.GetTask(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapTaskError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	requester, err := h.familyService.RequireMembership(ctx, user.ID, task.FamilyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if err := requireNotViewer(requester); err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Message: err.Error()})
	}

	var req models.UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	updated, err := h.taskService.UpdateTask(ctx, task.ID, models.UpdateTaskInput{
		Title:          req.Title,
		Description:    req.Description,
		Status:         req.Status,
		AssigneeID:     req.AssigneeID,
		ParentID:       req.ParentID,
		DueDate:        req.DueDate,
		RecurrenceRule: req.RecurrenceRule,
		ActorID:        requester.ID,
		ActorName:      memberDisplayName(requester),
	})
	if err != nil {
		status, msg := mapTaskError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, updated)
}

// AssignTask assigns a task to a family member.
func (h *TaskHandler) AssignTask(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	task, err := h.taskService.GetTask(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapTaskError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	requester, err := h.familyService.RequireMembership(ctx, user.ID, task.FamilyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if err := requireNotViewer(requester); err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Message: err.Error()})
	}

	var req models.AssignTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	updated, err := h.taskService.AssignTask(ctx, task.ID, req.MemberID, requester.ID, memberDisplayName(requester))
	if err != nil {
		status, msg := mapTaskError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, updated)
}

// DeleteTask removes a task from the board.
func (h *TaskHandler) DeleteTask(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	task, err := h.taskService.GetTask(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapTaskError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	requester, err := h.familyService.RequireMembership(ctx, user.ID, task.FamilyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if err := requireNotViewer(requester); err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Message: err.Error()})
	}

	if err := h.taskService.DeleteTask(ctx, task.ID); err != nil {
		status, msg := mapTaskError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.NoContent(http.StatusNoContent)
}

func mapTaskError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, repository.ErrTaskNotFound):
		return http.StatusNotFound, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
