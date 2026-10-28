package models

import "time"

// NotificationPreferences controls which emails a user receives. A
// missing row (ErrNotificationPreferencesNotFound) means "use the
// defaults" — rows are created lazily on first write, not pre-created
// for every signup. DigestCadence is reserved for a future digest-email
// phase (that phase needs a scheduler this app doesn't have yet) —
// stored now, acted on later.
type NotificationPreferences struct {
	UserID            string    `json:"user_id" db:"user_id"`
	TaskAssignedEmail bool      `json:"task_assigned_email" db:"task_assigned_email"`
	DigestCadence     string    `json:"digest_cadence" db:"digest_cadence"` // none, daily, weekly — unused this pass
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateNotificationPreferencesRequest is the JSON body for
// PATCH /me/notification-preferences.
type UpdateNotificationPreferencesRequest struct {
	TaskAssignedEmail *bool `json:"task_assigned_email,omitempty"`
}
