package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/middleware"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// ExportHandler exposes CSV exports of a family's expenses/tasks —
// read-only, like Activity/Contribution.
type ExportHandler struct {
	exportService service.ExportService
	familyService service.FamilyService
	authService   service.AuthService
}

// NewExportHandler constructs an ExportHandler.
func NewExportHandler(exportService service.ExportService, familyService service.FamilyService, authService service.AuthService) *ExportHandler {
	return &ExportHandler{exportService: exportService, familyService: familyService, authService: authService}
}

// RegisterRoutes wires the export endpoints onto the given group, behind
// middleware.RequireAuth.
func (h *ExportHandler) RegisterRoutes(g *echo.Group) {
	auth := middleware.RequireAuth(h.authService)
	g.GET("/families/:familyId/export/expenses.csv", h.ExportExpenses, auth)
	g.GET("/families/:familyId/export/tasks.csv", h.ExportTasks, auth)
}

// exportFilters parses the ?parent_id=&from=YYYY-MM-DD&to=YYYY-MM-DD
// query params shared by both export routes. A malformed date is
// silently ignored (treated as unset) rather than erroring the whole
// export over a filter mistake.
func exportFilters(c *echo.Context) (parentID *string, from, to *time.Time) {
	if v := c.QueryParam("parent_id"); v != "" {
		parentID = &v
	}
	if v := c.QueryParam("from"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			from = &t
		}
	}
	if v := c.QueryParam("to"); v != "" {
		if t, err := time.Parse("2006-01-02", v); err == nil {
			// End-of-day, so a "to" date includes everything that happened
			// on it rather than cutting off at midnight.
			t = t.Add(24*time.Hour - time.Nanosecond)
			to = &t
		}
	}
	return parentID, from, to
}

// ExportExpenses returns a family's expenses as a CSV download.
func (h *ExportHandler) ExportExpenses(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	parentID, from, to := exportFilters(c)
	csvBytes, err := h.exportService.ExportExpenses(ctx, familyID, parentID, from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "internal server error"})
	}
	c.Response().Header().Set("Content-Disposition", `attachment; filename="expenses.csv"`)
	return c.Blob(http.StatusOK, "text/csv", csvBytes)
}

// ExportTasks returns a family's tasks as a CSV download.
func (h *ExportHandler) ExportTasks(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	parentID, from, to := exportFilters(c)
	csvBytes, err := h.exportService.ExportTasks(ctx, familyID, parentID, from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "internal server error"})
	}
	c.Response().Header().Set("Content-Disposition", `attachment; filename="tasks.csv"`)
	return c.Blob(http.StatusOK, "text/csv", csvBytes)
}
