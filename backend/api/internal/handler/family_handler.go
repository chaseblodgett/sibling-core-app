package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// FamilyHandler exposes HTTP endpoints for creating/joining a family,
// inviting members by email, and managing membership.
type FamilyHandler struct {
	familyService service.FamilyService
}

// NewFamilyHandler constructs a FamilyHandler.
func NewFamilyHandler(familyService service.FamilyService) *FamilyHandler {
	return &FamilyHandler{familyService: familyService}
}

// RegisterRoutes wires family endpoints onto the given group.
func (h *FamilyHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/families", h.CreateFamily)
	g.GET("/families/:id", h.GetFamily)
	g.POST("/families/join", h.JoinFamily)
	g.POST("/families/:id/invites", h.InviteMember)
	g.GET("/families/:id/members", h.ListMembers)
	g.DELETE("/families/:id/members/:memberId", h.RemoveMember)
}

// CreateFamily creates a new family group with the requester as owner.
func (h *FamilyHandler) CreateFamily(c *echo.Context) error {
	var req models.CreateFamilyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// GetFamily returns a family by id.
func (h *FamilyHandler) GetFamily(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// JoinFamily adds the requester to a family using an invite code.
func (h *FamilyHandler) JoinFamily(c *echo.Context) error {
	var req models.JoinFamilyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// InviteMember invites a new member to the family by email.
func (h *FamilyHandler) InviteMember(c *echo.Context) error {
	var req models.InviteMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// ListMembers returns all members belonging to a family.
func (h *FamilyHandler) ListMembers(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// RemoveMember removes a member from a family.
func (h *FamilyHandler) RemoveMember(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}
