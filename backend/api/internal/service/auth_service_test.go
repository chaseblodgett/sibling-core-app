package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"sibling-core/api/internal/googleauth"
	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
	mockgoogleauth "sibling-core/api/test/mock/mock_googleauth"
	mockrepository "sibling-core/api/test/mock/mock_repository"
)

const testSecret = "test-secret-at-least-32-characters-long"

func newTestAuthService(t *testing.T) (service.AuthService, *mockrepository.MockUserRepository) {
	t.Helper()
	svc, repo, _ := newTestAuthServiceWithVerifier(t)
	return svc, repo
}

func newTestAuthServiceWithVerifier(t *testing.T) (service.AuthService, *mockrepository.MockUserRepository, *mockgoogleauth.MockVerifier) {
	t.Helper()
	repo := &mockrepository.MockUserRepository{}
	verifier := &mockgoogleauth.MockVerifier{}
	return service.NewAuthService(repo, testSecret, verifier), repo, verifier
}

func mustHash(t *testing.T, password string) *string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	hashStr := string(hash)
	return &hashStr
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
	require.NotNil(t, resp.User.PasswordHash)
	assert.NotEqual(t, "supersecret123", *resp.User.PasswordHash)
	assert.NotEmpty(t, *resp.User.PasswordHash)
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

func TestSignInWithGoogleCreatesNewUser(t *testing.T) {
	svc, repo, verifier := newTestAuthServiceWithVerifier(t)
	ctx := context.Background()

	verifier.VerifyFunc = func(ctx context.Context, idToken string) (*googleauth.Claims, error) {
		assert.Equal(t, "google-id-token", idToken)
		return &googleauth.Claims{Subject: "sub-1", Email: "New@Example.com", EmailVerified: true, Name: "New Person"}, nil
	}
	repo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		assert.Equal(t, "new@example.com", email)
		return nil, repository.ErrUserNotFound
	}
	var created *models.User
	repo.CreateFunc = func(ctx context.Context, u *models.User) error {
		u.ID = "user-1"
		u.CreatedAt = time.Now()
		created = u
		return nil
	}

	resp, err := svc.SignInWithGoogle(ctx, "google-id-token")

	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Nil(t, created.PasswordHash)
	require.NotNil(t, created.GoogleSub)
	assert.Equal(t, "sub-1", *created.GoogleSub)
	assert.Equal(t, "new@example.com", resp.User.Email)
	assert.NotEmpty(t, resp.Token)
}

func TestSignInWithGoogleLinksExistingPasswordAccount(t *testing.T) {
	svc, repo, verifier := newTestAuthServiceWithVerifier(t)
	ctx := context.Background()

	existing := &models.User{ID: "user-1", Email: "both@example.com", PasswordHash: mustHash(t, "supersecret123")}
	verifier.VerifyFunc = func(ctx context.Context, idToken string) (*googleauth.Claims, error) {
		return &googleauth.Claims{Subject: "sub-2", Email: "both@example.com", EmailVerified: true, Name: "Both"}, nil
	}
	repo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		return existing, nil
	}
	var linkedUserID, linkedSub string
	repo.LinkGoogleAccountFunc = func(ctx context.Context, userID, googleSub string) error {
		linkedUserID, linkedSub = userID, googleSub
		return nil
	}

	resp, err := svc.SignInWithGoogle(ctx, "google-id-token")

	require.NoError(t, err)
	assert.Equal(t, "user-1", linkedUserID)
	assert.Equal(t, "sub-2", linkedSub)
	assert.Equal(t, "user-1", resp.User.ID)
}

func TestSignInWithGoogleSkipsRelinkWhenAlreadyLinked(t *testing.T) {
	svc, repo, verifier := newTestAuthServiceWithVerifier(t)
	ctx := context.Background()

	sub := "sub-3"
	existing := &models.User{ID: "user-1", Email: "linked@example.com", GoogleSub: &sub}
	verifier.VerifyFunc = func(ctx context.Context, idToken string) (*googleauth.Claims, error) {
		return &googleauth.Claims{Subject: "sub-3", Email: "linked@example.com", EmailVerified: true}, nil
	}
	repo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
		return existing, nil
	}
	relinkCalled := false
	repo.LinkGoogleAccountFunc = func(ctx context.Context, userID, googleSub string) error {
		relinkCalled = true
		return nil
	}

	_, err := svc.SignInWithGoogle(ctx, "google-id-token")

	require.NoError(t, err)
	assert.False(t, relinkCalled, "should not re-link when google_sub already matches")
}

func TestSignInWithGoogleRejectsUnverifiedEmail(t *testing.T) {
	svc, _, verifier := newTestAuthServiceWithVerifier(t)

	verifier.VerifyFunc = func(ctx context.Context, idToken string) (*googleauth.Claims, error) {
		return &googleauth.Claims{Subject: "sub-4", Email: "unverified@example.com", EmailVerified: false}, nil
	}

	_, err := svc.SignInWithGoogle(context.Background(), "google-id-token")

	assert.ErrorIs(t, err, service.ErrGoogleEmailUnverified)
}

func TestSignInWithGoogleRejectsInvalidToken(t *testing.T) {
	svc, _, verifier := newTestAuthServiceWithVerifier(t)

	verifier.VerifyFunc = func(ctx context.Context, idToken string) (*googleauth.Claims, error) {
		return nil, googleauth.ErrInvalidToken
	}

	_, err := svc.SignInWithGoogle(context.Background(), "bad-token")

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
	svcA := service.NewAuthService(repoA, "secret-a-aaaaaaaaaaaaaaaaaaaaaaaa", &mockgoogleauth.MockVerifier{})
	resp, err := svcA.SignUp(ctx, models.SignUpInput{Email: "x@example.com", Password: "supersecret123", Name: "X"})
	require.NoError(t, err)

	repoB := &mockrepository.MockUserRepository{}
	svcB := service.NewAuthService(repoB, "secret-b-bbbbbbbbbbbbbbbbbbbbbbbb", &mockgoogleauth.MockVerifier{})

	_, err = svcB.ValidateToken(ctx, resp.Token)

	assert.ErrorIs(t, err, service.ErrInvalidToken)
}
