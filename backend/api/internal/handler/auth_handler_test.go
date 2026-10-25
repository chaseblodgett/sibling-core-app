package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sibling-core/api/internal/handler"
	"sibling-core/api/internal/middleware"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/service"
	mockservice "sibling-core/api/test/mock/mock_service"
)

func newTestContext(method, target, body string) (*echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	return echo.NewContext(req, rec, e), rec
}

func TestAuthHandlerSignUpSuccess(t *testing.T) {
	authSvc := &mockservice.MockAuthService{
		SignUpFunc: func(ctx context.Context, input models.SignUpInput) (*models.AuthResponse, error) {
			assert.Equal(t, models.SignUpInput{Email: "jane@example.com", Password: "supersecret123", Name: "Jane"}, input)
			return &models.AuthResponse{Token: "tok", User: &models.User{ID: "u1", Email: "jane@example.com", Name: "Jane"}}, nil
		},
	}
	h := handler.NewAuthHandler(authSvc)

	c, rec := newTestContext(http.MethodPost, "/auth/signup", `{"email":"jane@example.com","password":"supersecret123","name":"Jane"}`)

	err := h.SignUp(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var body models.AuthResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "tok", body.Token)
	assert.Equal(t, "jane@example.com", body.User.Email)
}

func TestAuthHandlerSignUpBadJSON(t *testing.T) {
	authSvc := &mockservice.MockAuthService{
		SignUpFunc: func(ctx context.Context, input models.SignUpInput) (*models.AuthResponse, error) {
			t.Fatal("SignUp should not be called for malformed input")
			return nil, nil
		},
	}
	h := handler.NewAuthHandler(authSvc)

	c, rec := newTestContext(http.MethodPost, "/auth/signup", `{"email":`)

	err := h.SignUp(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandlerLoginInvalidCredentials(t *testing.T) {
	authSvc := &mockservice.MockAuthService{
		LoginFunc: func(ctx context.Context, input models.LoginInput) (*models.AuthResponse, error) {
			assert.Equal(t, models.LoginInput{Email: "jane@example.com", Password: "wrong"}, input)
			return nil, service.ErrInvalidCredentials
		},
	}
	h := handler.NewAuthHandler(authSvc)

	c, rec := newTestContext(http.MethodPost, "/auth/login", `{"email":"jane@example.com","password":"wrong"}`)

	err := h.Login(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthHandlerMeWithMiddleware(t *testing.T) {
	user := &models.User{ID: "u1", Email: "jane@example.com", Name: "Jane"}
	authSvc := &mockservice.MockAuthService{
		ValidateTokenFunc: func(ctx context.Context, token string) (*models.User, error) {
			assert.Equal(t, "valid-token", token)
			return user, nil
		},
	}
	h := handler.NewAuthHandler(authSvc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()
	c := echo.NewContext(req, rec, e)

	wrapped := middleware.RequireAuth(authSvc)(h.Me)
	err := wrapped(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var body models.User
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, user.ID, body.ID)
}

func TestAuthHandlerMeMissingToken(t *testing.T) {
	authSvc := &mockservice.MockAuthService{
		ValidateTokenFunc: func(ctx context.Context, token string) (*models.User, error) {
			t.Fatal("ValidateToken should not be called without a token")
			return nil, nil
		},
	}
	h := handler.NewAuthHandler(authSvc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()
	c := echo.NewContext(req, rec, e)

	wrapped := middleware.RequireAuth(authSvc)(h.Me)
	err := wrapped(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
