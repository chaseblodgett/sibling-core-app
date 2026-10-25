// Package middleware holds cross-cutting Echo middleware shared across handlers.
package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
)

const userContextKey = "user"

func RequireAuth(authService service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "missing or malformed authorization header"})
			}

			user, err := authService.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "invalid or expired token"})
			}

			c.Set(userContextKey, user)
			return next(c)
		}
	}
}

func CurrentUser(c *echo.Context) *models.User {
	user, _ := c.Get(userContextKey).(*models.User)
	return user
}
