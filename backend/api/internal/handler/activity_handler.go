package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/middleware"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// ActivityHandler exposes a family's activity feed (read-only — entries
// are written as a side effect of the other handlers' own writes, not
// through any endpoint here).
type ActivityHandler struct {
	activityService service.ActivityService
	familyService   service.FamilyService
	authService     service.AuthService
}

// NewActivityHandler constructs an ActivityHandler.
func NewActivityHandler(activityService service.ActivityService, familyService service.FamilyService, authService service.AuthService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService, familyService: familyService, authService: authService}
}

// RegisterRoutes wires the activity endpoint onto the given group, behind
// middleware.RequireAuth.
func (h *ActivityHandler) RegisterRoutes(g *echo.Group) {
	auth := middleware.RequireAuth(h.authService)
	g.GET("/families/:familyId/activity", h.ListActivity, auth)
}

// ListActivity returns a family's activity feed, newest first. Optional
// ?limit= (default 20, max 100) and ?offset= query params page through it.
func (h *ActivityHandler) ListActivity(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	entries, err := h.activityService.ListActivity(ctx, familyID, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "internal server error"})
	}
	return c.JSON(http.StatusOK, entries)
}
