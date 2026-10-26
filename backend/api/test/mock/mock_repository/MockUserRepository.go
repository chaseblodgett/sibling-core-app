package mockrepository

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockUserRepository is a test double for repository.UserRepository. Assign
// the *Func fields your test exercises; unset ones return zero values.
type MockUserRepository struct {
	CreateFunc            func(ctx context.Context, user *models.User) error
	FindByEmailFunc       func(ctx context.Context, email string) (*models.User, error)
	FindByIDFunc          func(ctx context.Context, id string) (*models.User, error)
	LinkGoogleAccountFunc func(ctx context.Context, userID, googleSub string) error
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	if m.CreateFunc == nil {
		return nil
	}
	return m.CreateFunc(ctx, user)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.FindByEmailFunc == nil {
		return nil, nil
	}
	return m.FindByEmailFunc(ctx, email)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	if m.FindByIDFunc == nil {
		return nil, nil
	}
	return m.FindByIDFunc(ctx, id)
}

func (m *MockUserRepository) LinkGoogleAccount(ctx context.Context, userID, googleSub string) error {
	if m.LinkGoogleAccountFunc == nil {
		return nil
	}
	return m.LinkGoogleAccountFunc(ctx, userID, googleSub)
}
