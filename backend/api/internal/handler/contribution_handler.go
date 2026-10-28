package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/middleware"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// ContributionHandler exposes a family's contribution breakdown
// (read-only — it's a computed aggregation, not a resource with its own
// writes).
type ContributionHandler struct {
	contributionService service.ContributionService
	familyService       service.FamilyService
	authService         service.AuthService
}

// NewContributionHandler constructs a ContributionHandler.
func NewContributionHandler(contributionService service.ContributionService, familyService service.FamilyService, authService service.AuthService) *ContributionHandler {
	return &ContributionHandler{contributionService: contributionService, familyService: familyService, authService: authService}
}

// RegisterRoutes wires the contributions endpoint onto the given group,
// behind middleware.RequireAuth.
func (h *ContributionHandler) RegisterRoutes(g *echo.Group) {
	auth := middleware.RequireAuth(h.authService)
	g.GET("/families/:familyId/contributions", h.GetContributions, auth)
}

// GetContributions returns a family's per-member contribution breakdown.
func (h *ContributionHandler) GetContributions(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	contributions, err := h.contributionService.GetContributions(ctx, familyID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: "internal server error"})
	}
	return c.JSON(http.StatusOK, contributions)
}
