package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// TaskHandler exposes HTTP endpoints for the shared family task board.
type TaskHandler struct {
	taskService service.TaskService
}

// NewTaskHandler constructs a TaskHandler.
func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// RegisterRoutes wires task endpoints onto the given group.
func (h *TaskHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/families/:familyId/tasks", h.CreateTask)
	g.GET("/families/:familyId/tasks", h.ListTasks)
	g.GET("/tasks/:id", h.GetTask)
	g.PATCH("/tasks/:id", h.UpdateTask)
	g.POST("/tasks/:id/assign", h.AssignTask)
	g.DELETE("/tasks/:id", h.DeleteTask)
}

// CreateTask creates a new task on a family's board.
func (h *TaskHandler) CreateTask(c *echo.Context) error {
	var req models.CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// ListTasks returns all tasks for a family.
func (h *TaskHandler) ListTasks(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// GetTask returns a single task by id.
func (h *TaskHandler) GetTask(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// UpdateTask updates fields on an existing task.
func (h *TaskHandler) UpdateTask(c *echo.Context) error {
	var req models.UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// AssignTask assigns a task to a family member.
func (h *TaskHandler) AssignTask(c *echo.Context) error {
	var req models.AssignTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// DeleteTask removes a task from the board.
func (h *TaskHandler) DeleteTask(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}
