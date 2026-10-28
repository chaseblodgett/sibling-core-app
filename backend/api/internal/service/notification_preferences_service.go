package service

import (
	"context"
	"errors"

	"sibling-core/api/internal/models"
	"sibling-core/api/internal/repository"
)

// defaultNotificationPreferences is what a user who has never saved
// preferences gets — email notifications on, digest off (digest isn't
// built yet regardless).
func defaultNotificationPreferences(userID string) *models.NotificationPreferences {
	return &models.NotificationPreferences{UserID: userID, TaskAssignedEmail: true, DigestCadence: "none"}
}

// NotificationPreferencesService manages a user's notification settings.
type NotificationPreferencesService interface {
	// GetPreferences returns userID's saved preferences, or the defaults
	// (not persisted) if they've never saved any.
	GetPreferences(ctx context.Context, userID string) (*models.NotificationPreferences, error)
	UpdatePreferences(ctx context.Context, userID string, input models.UpdateNotificationPreferencesRequest) (*models.NotificationPreferences, error)
}

type notificationPreferencesService struct {
	repo repository.NotificationPreferencesRepository
}

// NewNotificationPreferencesService constructs a NotificationPreferencesService.
func NewNotificationPreferencesService(repo repository.NotificationPreferencesRepository) NotificationPreferencesService {
	return &notificationPreferencesService{repo: repo}
}

func (s *notificationPreferencesService) GetPreferences(ctx context.Context, userID string) (*models.NotificationPreferences, error) {
	prefs, err := s.repo.GetByUserID(ctx, userID)
	if errors.Is(err, repository.ErrNotificationPreferencesNotFound) {
		return defaultNotificationPreferences(userID), nil
	}
	if err != nil {
		return nil, err
	}
	return prefs, nil
}

func (s *notificationPreferencesService) UpdatePreferences(ctx context.Context, userID string, input models.UpdateNotificationPreferencesRequest) (*models.NotificationPreferences, error) {
	prefs, err := s.GetPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}
	if input.TaskAssignedEmail != nil {
		prefs.TaskAssignedEmail = *input.TaskAssignedEmail
	}
	if err := s.repo.Upsert(ctx, prefs); err != nil {
		return nil, err
	}
	return prefs, nil
}
