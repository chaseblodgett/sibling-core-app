package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/middleware"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// NotificationPreferencesHandler exposes the current user's own
// notification settings — not scoped to a family, since preferences are
// per-user.
type NotificationPreferencesHandler struct {
	prefsService service.NotificationPreferencesService
	authService  service.AuthService
}

// NewNotificationPreferencesHandler constructs a NotificationPreferencesHandler.
func NewNotificationPreferencesHandler(prefsService service.NotificationPreferencesService, authService service.AuthService) *NotificationPreferencesHandler {
	return &NotificationPreferencesHandler{prefsService: prefsService, authService: authService}
}

// RegisterRoutes wires the preferences endpoints onto the given group,
// behind middleware.RequireAuth.
func (h *NotificationPreferencesHandler) RegisterRoutes(g *echo.Group) {
	auth := middleware.RequireAuth(h.authService)
	g.GET("/me/notification-preferences", h.GetPreferences, auth)
	g.PATCH("/me/notification-preferences", h.UpdatePreferences, auth)
}

// GetPreferences returns the current user's notification preferences
// (defaults if they've never saved any).
func (h *NotificationPreferencesHandler) GetPreferences(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	prefs, err := h.prefsService.GetPreferences(ctx, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "internal server error"})
	}
	return c.JSON(http.StatusOK, prefs)
}

// UpdatePreferences updates the current user's notification preferences.
func (h *NotificationPreferencesHandler) UpdatePreferences(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	var req models.UpdateNotificationPreferencesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	prefs, err := h.prefsService.UpdatePreferences(ctx, user.ID, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "internal server error"})
	}
	return c.JSON(http.StatusOK, prefs)
}
