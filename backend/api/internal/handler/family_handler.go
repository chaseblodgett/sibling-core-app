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

// FamilyHandler exposes HTTP endpoints for creating/joining a family,
// inviting members by email, and managing membership. Every route requires
// authentication — there is no public family data.
type FamilyHandler struct {
	familyService service.FamilyService
	authService   service.AuthService
}

// NewFamilyHandler constructs a FamilyHandler.
func NewFamilyHandler(familyService service.FamilyService, authService service.AuthService) *FamilyHandler {
	return &FamilyHandler{familyService: familyService, authService: authService}
}

// RegisterRoutes wires family endpoints onto the given group, all behind
// middleware.RequireAuth.
func (h *FamilyHandler) RegisterRoutes(g *echo.Group) {
	fg := g.Group("/families", middleware.RequireAuth(h.authService))
	fg.POST("", h.CreateFamily)
	fg.GET("/mine", h.ListMine)
	fg.POST("/join", h.JoinFamily)
	fg.GET("/:id", h.GetFamily)
	fg.POST("/:id/invites", h.InviteMember)
	fg.GET("/:id/members", h.ListMembers)
	fg.DELETE("/:id/members/:memberId", h.RemoveMember)
	fg.PATCH("/:id/members/:memberId/role", h.UpdateMemberRole)
}

// CreateFamily creates a new family group with the requester as owner.
func (h *FamilyHandler) CreateFamily(c *echo.Context) error {
	var req models.CreateFamilyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	user := middleware.CurrentUser(c)
	family, err := h.familyService.CreateFamily(c.Request().Context(), models.CreateFamilyInput{
		Name:      req.Name,
		CreatedBy: user.ID,
	})
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusCreated, family)
}

// ListMine returns every family the requester belongs to.
func (h *FamilyHandler) ListMine(c *echo.Context) error {
	user := middleware.CurrentUser(c)
	families, err := h.familyService.ListFamiliesForUser(c.Request().Context(), user.ID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, families)
}

// GetFamily returns a family by id, if the requester belongs to it.
func (h *FamilyHandler) GetFamily(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("id")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	family, err := h.familyService.GetFamily(ctx, familyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, family)
}

// JoinFamily adds the requester to a family using an invite code.
func (h *FamilyHandler) JoinFamily(c *echo.Context) error {
	var req models.JoinFamilyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	user := middleware.CurrentUser(c)
	family, err := h.familyService.JoinFamily(c.Request().Context(), req.InviteCode, user.ID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, family)
}

// InviteMember invites a new member to the family by email. Only the
// family's owner may invite.
func (h *FamilyHandler) InviteMember(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("id")

	requester, err := h.familyService.RequireMembership(ctx, user.ID, familyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if requester.Role != "owner" {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Message: service.ErrNotAuthorized.Error()})
	}

	var req models.InviteMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	invite, err := h.familyService.InviteMember(ctx, models.InviteMemberInput{
		FamilyID:  familyID,
		Email:     req.Email,
		InvitedBy: user.ID,
	})
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusCreated, invite)
}

// ListMembers returns all members belonging to a family.
func (h *FamilyHandler) ListMembers(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("id")

	if _, err := h.familyService.RequireMembership(ctx, user.ID, familyID); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	members, err := h.familyService.ListMembers(ctx, familyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, members)
}

// RemoveMember removes a member from a family. Only the family's owner may
// remove members.
func (h *FamilyHandler) RemoveMember(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("id")

	requester, err := h.familyService.RequireMembership(ctx, user.ID, familyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if requester.Role != "owner" {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Message: service.ErrNotAuthorized.Error()})
	}

	if err := h.familyService.RemoveMember(ctx, familyID, c.Param("memberId")); err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.NoContent(http.StatusNoContent)
}

// UpdateMemberRole changes a member's role (owner, member, or viewer).
// Only the family's owner may change roles.
func (h *FamilyHandler) UpdateMemberRole(c *echo.Context) error {
	ctx := c.Request().Context()
	user := middleware.CurrentUser(c)
	familyID := c.Param("id")

	requester, err := h.familyService.RequireMembership(ctx, user.ID, familyID)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	if requester.Role != "owner" {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{Message: service.ErrNotAuthorized.Error()})
	}

	var req models.UpdateMemberRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	updated, err := h.familyService.UpdateMemberRole(ctx, familyID, c.Param("memberId"), req.Role)
	if err != nil {
		status, msg := mapFamilyError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}
	return c.JSON(http.StatusOK, updated)
}

func mapFamilyError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrInvalidInviteCode):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, service.ErrNotMember):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, service.ErrNotAuthorized):
		return http.StatusForbidden, err.Error()
	case errors.Is(err, repository.ErrFamilyNotFound), errors.Is(err, repository.ErrMemberNotFound), errors.Is(err, repository.ErrInviteNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, repository.ErrAlreadyMember):
		return http.StatusConflict, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
