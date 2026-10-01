package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

// Handlers bundles every resource handler so routes can be registered in
// one place from the DI container.
type Handlers struct {
	Auth     *AuthHandler
	Family   *FamilyHandler
	Task     *TaskHandler
	Expense  *ExpenseHandler
	Document *DocumentHandler
}

// RegisterRoutes mounts all API routes under /api/v1, plus a top-level
// health check.
func RegisterRoutes(e *echo.Echo, h *Handlers) {
	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	v1 := e.Group("/api/v1")

	h.Auth.RegisterRoutes(v1)
	h.Family.RegisterRoutes(v1)
	h.Task.RegisterRoutes(v1)
	h.Expense.RegisterRoutes(v1)
	h.Document.RegisterRoutes(v1)
}
