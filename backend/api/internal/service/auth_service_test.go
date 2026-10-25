package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
	mockrepository "sibling-core/api/test/mock/mock_repository"
)

const testSecret = "test-secret-at-least-32-characters-long"

func newTestAuthService(t *testing.T) (service.AuthService, *mockrepository.MockUserRepository) {
	t.Helper()
	repo := &mockrepository.MockUserRepository{}
	return service.NewAuthService(repo, testSecret), repo
}

func mustHash(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	return string(hash)
}

func TestSignUpSuccess(t *testing.T) {
	svc, repo := newTestAuthService(t)
	ctx := context.Background()

	repo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		assert.Equal(t, "jane@example.com", email)
		return nil, repository.ErrUserNotFound
	}
	repo.CreateFunc = func(ctx context.Context, u *models.User) error {
		u.ID = "user-1"
		u.CreatedAt = time.Now()
		return nil
	}

	resp, err := svc.SignUp(ctx, models.SignUpInput{
		Email:    "Jane@Example.com",
		Password: "supersecret123",
		Name:     "Jane Doe",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "jane@example.com", resp.User.Email)
	assert.NotEqual(t, "supersecret123", resp.User.PasswordHash)
	assert.NotEmpty(t, resp.User.PasswordHash)
}

func TestSignUpDuplicateEmail(t *testing.T) {
	svc, repo := newTestAuthService(t)
	ctx := context.Background()

	existing := &models.User{ID: "user-1", Email: "dup@example.com"}
	repo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		assert.Equal(t, "dup@example.com", email)
		return existing, nil
	}

	_, err := svc.SignUp(ctx, models.SignUpInput{Email: "dup@example.com", Password: "supersecret123", Name: "A"})

	require.Error(t, err)
	assert.ErrorIs(t, err, service.ErrEmailTaken)
}

func TestSignUpInvalidInput(t *testing.T) {
	svc, _ := newTestAuthService(t)
	ctx := context.Background()

	cases := []models.SignUpInput{
		{Email: "not-an-email", Password: "supersecret123", Name: "A"},
		{Email: "a@example.com", Password: "short", Name: "A"},
		{Email: "a@example.com", Password: "supersecret123", Name: "  "},
	}
	for _, in := range cases {
		_, err := svc.SignUp(ctx, in)
		assert.ErrorIsf(t, err, service.ErrInvalidInput, "input %+v", in)
	}
}

func TestLoginSuccess(t *testing.T) {
	svc, repo := newTestAuthService(t)
	ctx := context.Background()

	user := &models.User{ID: "user-1", Email: "login@example.com", PasswordHash: mustHash(t, "supersecret123")}
	repo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		assert.Equal(t, "login@example.com", email)
		return user, nil
	}

	resp, err := svc.Login(ctx, models.LoginInput{Email: "LOGIN@example.com", Password: "supersecret123"})

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, user.ID, resp.User.ID)
}

func TestLoginWrongPassword(t *testing.T) {
	svc, repo := newTestAuthService(t)
	ctx := context.Background()

	user := &models.User{ID: "user-1", Email: "wrong@example.com", PasswordHash: mustHash(t, "therealpassword")}
	repo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		return user, nil
	}

	_, err := svc.Login(ctx, models.LoginInput{Email: "wrong@example.com", Password: "nope12345"})

	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestLoginUnknownEmail(t *testing.T) {
	svc, repo := newTestAuthService(t)
	ctx := context.Background()

	repo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		return nil, repository.ErrUserNotFound
	}

	_, err := svc.Login(ctx, models.LoginInput{Email: "ghost@example.com", Password: "whatever1"})

	assert.ErrorIs(t, err, service.ErrInvalidCredentials)
}

func TestValidateTokenRoundTrip(t *testing.T) {
	svc, repo := newTestAuthService(t)
	ctx := context.Background()

	repo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		return nil, repository.ErrUserNotFound
	}
	repo.CreateFunc = func(ctx context.Context, u *models.User) error {
		u.ID = "user-1"
		u.CreatedAt = time.Now()
		return nil
	}

	resp, err := svc.SignUp(ctx, models.SignUpInput{Email: "me@example.com", Password: "supersecret123", Name: "Me"})
	require.NoError(t, err)

	repo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
		assert.Equal(t, resp.User.ID, id)
		return resp.User, nil
	}

	user, err := svc.ValidateToken(ctx, resp.Token)

	require.NoError(t, err)
	assert.Equal(t, resp.User.ID, user.ID)
}

func TestValidateTokenGarbage(t *testing.T) {
	svc, _ := newTestAuthService(t)

	_, err := svc.ValidateToken(context.Background(), "not-a-real-token")

	assert.ErrorIs(t, err, service.ErrInvalidToken)
}

func TestValidateTokenWrongSecret(t *testing.T) {
	ctx := context.Background()

	repoA := &mockrepository.MockUserRepository{
		FindByEmailFunc: func(ctx context.Context, email string) (*models.User, error) {
			return nil, repository.ErrUserNotFound
		},
		CreateFunc: func(ctx context.Context, u *models.User) error {
			u.ID = "user-1"
			u.CreatedAt = time.Now()
			return nil
		},
	}
	svcA := service.NewAuthService(repoA, "secret-a-aaaaaaaaaaaaaaaaaaaaaaaa")
	resp, err := svcA.SignUp(ctx, models.SignUpInput{Email: "x@example.com", Password: "supersecret123", Name: "X"})
	require.NoError(t, err)

	repoB := &mockrepository.MockUserRepository{}
	svcB := service.NewAuthService(repoB, "secret-b-bbbbbbbbbbbbbbbbbbbbbbbb")

	_, err = svcB.ValidateToken(ctx, resp.Token)

	assert.ErrorIs(t, err, service.ErrInvalidToken)
}
