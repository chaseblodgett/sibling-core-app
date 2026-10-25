package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/middleware"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

// AuthHandler exposes HTTP endpoints for sign up, login, and the current
// authenticated user.
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
	g.GET("/auth/me", h.Me, middleware.RequireAuth(h.authService))
}

// SignUp registers a new user account and returns a session token.
func (h *AuthHandler) SignUp(c *echo.Context) error {
	var req models.SignUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	result, err := h.authService.SignUp(c.Request().Context(), models.SignUpInput{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		status, msg := mapAuthError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	return c.JSON(http.StatusCreated, result)
}

// Login authenticates a user and returns a session token.
func (h *AuthHandler) Login(c *echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	result, err := h.authService.Login(c.Request().Context(), models.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		status, msg := mapAuthError(err)
		return c.JSON(status, models.ErrorResponse{Message: msg})
	}

	return c.JSON(http.StatusOK, result)
}

// Me returns the currently authenticated user, as attached by RequireAuth.
func (h *AuthHandler) Me(c *echo.Context) error {
	return c.JSON(http.StatusOK, middleware.CurrentUser(c))
}

func mapAuthError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, service.ErrEmailTaken):
		return http.StatusConflict, err.Error()
	case errors.Is(err, service.ErrInvalidCredentials):
		return http.StatusUnauthorized, err.Error()
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
