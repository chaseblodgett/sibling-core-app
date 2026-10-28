package mockrepository

import (
	"context"

	"sibling-core/api/internal/models"
)

// MockNotificationPreferencesRepository is a test double for
// repository.NotificationPreferencesRepository. Assign the *Func fields
// your test exercises; unset ones return zero values.
type MockNotificationPreferencesRepository struct {
	GetByUserIDFunc func(ctx context.Context, userID string) (*models.NotificationPreferences, error)
	UpsertFunc      func(ctx context.Context, prefs *models.NotificationPreferences) error
}

func (m *MockNotificationPreferencesRepository) GetByUserID(ctx context.Context, userID string) (*models.NotificationPreferences, error) {
	if m.GetByUserIDFunc == nil {
		return nil, nil
	}
	return m.GetByUserIDFunc(ctx, userID)
}

func (m *MockNotificationPreferencesRepository) Upsert(ctx context.Context, prefs *models.NotificationPreferences) error {
	if m.UpsertFunc == nil {
		return nil
	}
	return m.UpsertFunc(ctx, prefs)
}
