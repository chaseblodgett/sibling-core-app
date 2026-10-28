package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
	"sibling-core/api/internal/service"
	mockrepository "sibling-core/api/test/mock/mock_repository"
)

func TestGetPreferencesReturnsDefaultsWhenNoneSaved(t *testing.T) {
	repo := &mockrepository.MockNotificationPreferencesRepository{}
	svc := service.NewNotificationPreferencesService(repo)

	repo.GetByUserIDFunc = func(ctx context.Context, userID string) (*models.NotificationPreferences, error) {
		return nil, repository.ErrNotificationPreferencesNotFound
	}

	prefs, err := svc.GetPreferences(context.Background(), "user-1")

	require.NoError(t, err)
	assert.True(t, prefs.TaskAssignedEmail)
	assert.Equal(t, "none", prefs.DigestCadence)
}

func TestUpdatePreferencesPersistsOptOut(t *testing.T) {
	repo := &mockrepository.MockNotificationPreferencesRepository{}
	svc := service.NewNotificationPreferencesService(repo)

	repo.GetByUserIDFunc = func(ctx context.Context, userID string) (*models.NotificationPreferences, error) {
		return nil, repository.ErrNotificationPreferencesNotFound
	}
	var upserted *models.NotificationPreferences
	repo.UpsertFunc = func(ctx context.Context, prefs *models.NotificationPreferences) error {
		upserted = prefs
		return nil
	}

	disabled := false
	prefs, err := svc.UpdatePreferences(context.Background(), "user-1", models.UpdateNotificationPreferencesRequest{TaskAssignedEmail: &disabled})

	require.NoError(t, err)
	require.NotNil(t, upserted)
	assert.False(t, prefs.TaskAssignedEmail)
	assert.False(t, upserted.TaskAssignedEmail)
}
