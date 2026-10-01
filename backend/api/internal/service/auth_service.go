package service

import (
	"context"
	"errors"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// AuthService handles account creation and authentication.
// Actual user identity/credentials are delegated to Supabase Auth; this
// service is responsible for wiring that up and syncing local Member rows.
type AuthService interface {
	SignUp(ctx context.Context, input models.SignUpInput) (*models.AuthResult, error)
	Login(ctx context.Context, input models.LoginInput) (*models.AuthResult, error)
	Logout(ctx context.Context, userID string) error
	Me(ctx context.Context, userID string) (*models.Member, error)
}

type authService struct {
	memberRepo repository.MemberRepository
}

// NewAuthService constructs an AuthService.
func NewAuthService(memberRepo repository.MemberRepository) AuthService {
	return &authService{memberRepo: memberRepo}
}

func (s *authService) SignUp(ctx context.Context, input models.SignUpInput) (*models.AuthResult, error) {
	return nil, errors.New("not implemented")
}

func (s *authService) Login(ctx context.Context, input models.LoginInput) (*models.AuthResult, error) {
	return nil, errors.New("not implemented")
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	return errors.New("not implemented")
}

func (s *authService) Me(ctx context.Context, userID string) (*models.Member, error) {
	return nil, errors.New("not implemented")
}
