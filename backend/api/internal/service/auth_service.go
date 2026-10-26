package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"sibling-core/api/internal/googleauth"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// tokenTTL controls how long an issued JWT stays valid. There is no refresh
// token or server-side revocation in this scaffold, so this is a straight
// expiry: picking a short TTL would log users out constantly, so this
// trades a slightly longer exposure window for not needing a refresh flow
// yet. Revisit alongside a refresh-token design if that tradeoff stops
// being acceptable.
const tokenTTL = 7 * 24 * time.Hour

var (
	ErrEmailTaken            = errors.New("email already registered")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrInvalidInput          = errors.New("invalid input")
	ErrInvalidToken          = errors.New("invalid or expired token")
	ErrGoogleEmailUnverified = errors.New("google account email is not verified")
)

// AuthService handles account creation, authentication, and JWT validation.
type AuthService interface {
	SignUp(ctx context.Context, input models.SignUpInput) (*models.AuthResponse, error)
	Login(ctx context.Context, input models.LoginInput) (*models.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)
	// SignInWithGoogle verifies a Google ID token and signs the user in,
	// creating an account (or linking an existing password account by
	// verified email) on first sign-in.
	SignInWithGoogle(ctx context.Context, idToken string) (*models.AuthResponse, error)
}

type authService struct {
	userRepo       repository.UserRepository
	jwtSecret      []byte
	googleVerifier googleauth.Verifier
}

// NewAuthService constructs an AuthService.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, googleVerifier googleauth.Verifier) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: []byte(jwtSecret), googleVerifier: googleVerifier}
}

func (s *authService) SignUp(ctx context.Context, input models.SignUpInput) (*models.AuthResponse, error) {
	if err := validateSignUp(input); err != nil {
		return nil, err
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))

	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashStr := string(hash)

	user := &models.User{
		Email:        email,
		PasswordHash: &hashStr,
		Name:         strings.TrimSpace(input.Name),
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}

	return s.authResponse(user)
}

func (s *authService) Login(ctx context.Context, input models.LoginInput) (*models.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if user.PasswordHash == nil {
		// Google-only account — no password to compare against. Same
		// error as a wrong password, so we don't leak how the account
		// was created to an unauthenticated caller.
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.authResponse(user)
}

func (s *authService) SignInWithGoogle(ctx context.Context, idToken string) (*models.AuthResponse, error) {
	claims, err := s.googleVerifier.Verify(ctx, idToken)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if !claims.EmailVerified {
		return nil, ErrGoogleEmailUnverified
	}

	email := strings.ToLower(strings.TrimSpace(claims.Email))

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	if user == nil {
		user = &models.User{
			Email:     email,
			Name:      claims.Name,
			GoogleSub: &claims.Subject,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	} else if user.GoogleSub == nil || *user.GoogleSub != claims.Subject {
		// First Google sign-in on a pre-existing password account: link it.
		// Safe because Google has already told us this email is verified.
		if err := s.userRepo.LinkGoogleAccount(ctx, user.ID, claims.Subject); err != nil {
			return nil, err
		}
		user.GoogleSub = &claims.Subject
	}

	return s.authResponse(user)
}

func (s *authService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	user, err := s.userRepo.FindByID(ctx, claims.Subject)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return user, nil
}

func (s *authService) authResponse(user *models.User) (*models.AuthResponse, error) {
	claims := jwt.RegisteredClaims{
		Subject:   user.ID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}
	return &models.AuthResponse{Token: token, User: user}, nil
}

func validateSignUp(input models.SignUpInput) error {
	if !strings.Contains(input.Email, "@") {
		return fmt.Errorf("%w: invalid email", ErrInvalidInput)
	}
	if len(input.Password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", ErrInvalidInput)
	}
	if strings.TrimSpace(input.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	return nil
}
