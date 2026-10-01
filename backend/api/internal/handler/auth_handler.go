package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// AuthHandler exposes HTTP endpoints for sign up, login, logout, and the
// current authenticated user.
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRoutes wires auth endpoints onto the given group.
func (h *AuthHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/auth/signup", h.SignUp)
	g.POST("/auth/login", h.Login)
	g.POST("/auth/logout", h.Logout)
	g.GET("/auth/me", h.Me)
}

// SignUp registers a new user account.
func (h *AuthHandler) SignUp(c *echo.Context) error {
	var req models.SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// Login authenticates a user and returns session tokens.
func (h *AuthHandler) Login(c *echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// Logout invalidates the current session.
func (h *AuthHandler) Logout(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}

// Me returns the currently authenticated user's profile/membership info.
func (h *AuthHandler) Me(c *echo.Context) error {
	return c.JSON(http.StatusNotImplemented, models.ErrorResponse{Message: "not implemented"})
}
