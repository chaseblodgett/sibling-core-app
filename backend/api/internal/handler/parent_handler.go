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

// ParentHandler exposes HTTP endpoints for a family's tracked parents
// (care-recipients). Every route requires the requester to belong to the
// family the parent is (or will be) in; deleting additionally requires
// the owner role, matching RemoveMember's safety bar for shared records.
type ParentHandler struct {
	parentService service.ParentService
	familyService service.FamilyService
	authService   service.AuthService
}

// NewParentHandler constructs a ParentHandler.
func NewParentHandler(parentService service.ParentService, familyService service.FamilyService, authService service.AuthService) *ParentHandler {
	return &ParentHandler{parentService: parentService, familyService: familyService, authService: authService}
}

// RegisterRoutes wires parent endpoints onto the given group, all behind
// middleware.RequireAuth.
func (h *ParentHandler) RegisterRoutes(g *echo.Group) {
	auth := middleware.RequireAuth(h.authService)
	g.POST("/families/:familyId/parents", h.CreateParent, auth)
	g.GET("/families/:familyId/parents", h.ListParents, auth)
	g.GET("/parents/:id", h.GetParent, auth)
	g.PATCH("/parents/:id", h.UpdateParent, auth)
	g.DELETE("/parents/:id", h.DeleteParent, auth)
}

// CreateParent adds a new parent to a family.
func (h *ParentHandler) CreateParent(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	var req models.CreateParentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	parent, err := h.parentService.CreateParent(ctx, models.CreateParentInput{
		FamilyID:  familyID,
		Name:      req.Name,
		Relation:  req.Relation,
		CreatedBy: user.ID,
	})
	if err != nil {
		status, msg := mapParentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusCreated, parent)
}

// ListParents returns every parent tracked by a family.
func (h *ParentHandler) ListParents(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("familyId")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	parents, err := h.parentService.ListParents(ctx, familyID)
	if err != nil {
		status, msg := mapParentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, parents)
}

// GetParent returns a single parent by id, if the requester belongs to
// their family.
func (h *ParentHandler) GetParent(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	parent, err := h.parentService.GetParent(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapParentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if _, err := h.familyService.RequireMembership(ctx, user.ID, parent.FamilyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, parent)
}

// UpdateParent updates a parent's name/relation.
func (h *ParentHandler) UpdateParent(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	parent, err := h.parentService.GetParent(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapParentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if _, err := h.familyService.RequireMembership(ctx, user.ID, parent.FamilyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	var req models.UpdateParentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	updated, err := h.parentService.UpdateParent(ctx, parent.ID, models.UpdateParentInput{
		Name:                  req.Name,
		Relation:              req.Relation,
		DoctorName:            req.DoctorName,
		DoctorPhone:           req.DoctorPhone,
		Allergies:             req.Allergies,
		Medications:           req.Medications,
		InsuranceProvider:     req.InsuranceProvider,
		InsuranceMemberID:     req.InsuranceMemberID,
		EmergencyContactName:  req.EmergencyContactName,
		EmergencyContactPhone: req.EmergencyContactPhone,
	})
	if err != nil {
		status, msg := mapParentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, updated)
}

// DeleteParent removes a parent. Only the family's owner may do this.
func (h *ParentHandler) DeleteParent(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)

	parent, err := h.parentService.GetParent(ctx, c.Param("id"))
	if err != nil {
		status, msg := mapParentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	requester, err := h.familyService.RequireMembership(ctx, user.ID, parent.FamilyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if requester.Role != "owner" {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Message: service.ErrNotAuthorized.Error()})
	}

	if err := h.parentService.DeleteParent(ctx, parent.ID); err != nil {
		status, msg := mapParentError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.NoContent(http.StatusNoContent)
}

func mapParentError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, repository.ErrParentNotFound):
		return http.StatusNotFound, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
