package mockservice

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockAuthService is a test double for service.AuthService. Assign the
// *Func fields your test exercises; unset ones return zero values.
type MockAuthService struct {
	SignUpFunc           func(ctx context.Context, input models.SignUpInput) (*models.AuthResponse, error)
	LoginFunc            func(ctx context.Context, input models.LoginInput) (*models.AuthResponse, error)
	ValidateTokenFunc    func(ctx context.Context, token string) (*models.User, error)
	SignInWithGoogleFunc func(ctx context.Context, idToken string) (*models.AuthResponse, error)
}

func (m *MockAuthService) SignUp(ctx context.Context, input models.SignUpInput) (*models.AuthResponse, error) {
	if m.SignUpFunc == nil {
		return nil, nil
	}
	return m.SignUpFunc(ctx, input)
}

func (m *MockAuthService) Login(ctx context.Context, input models.LoginInput) (*models.AuthResponse, error) {
	if m.LoginFunc == nil {
		return nil, nil
	}
	return m.LoginFunc(ctx, input)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	if m.ValidateTokenFunc == nil {
		return nil, nil
	}
	return m.ValidateTokenFunc(ctx, token)
}

func (m *MockAuthService) SignInWithGoogle(ctx context.Context, idToken string) (*models.AuthResponse, error) {
	if m.SignInWithGoogleFunc == nil {
		return nil, nil
	}
	return m.SignInWithGoogleFunc(ctx, idToken)
}
