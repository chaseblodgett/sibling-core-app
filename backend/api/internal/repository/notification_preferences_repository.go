package repository

import (
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

var ErrNotificationPreferencesNotFound = errors.New("notification preferences not found")

const notificationPreferencesTable = "notification_preferences"

// NotificationPreferencesRepository defines data access for a user's
// notification preferences. A missing row (ErrNotificationPreferencesNotFound)
// is a normal, expected state — it means "use the defaults" — not an error
// callers should treat as unexpected.
type NotificationPreferencesRepository interface {
	GetByUserID(ctx context.Context, userID string) (*models.NotificationPreferences, error)
	// Upsert creates or updates the row for prefs.UserID.
	Upsert(ctx context.Context, prefs *models.NotificationPreferences) error
}

// postgresNotificationPreferencesRepository is the Postgres-backed
// implementation of NotificationPreferencesRepository.
type postgresNotificationPreferencesRepository struct {
	db *pgxpool.Pool
}

// NewNotificationPreferencesRepository constructs a Postgres-backed
// NotificationPreferencesRepository.
func NewNotificationPreferencesRepository(db *pgxpool.Pool) NotificationPreferencesRepository {
	return &postgresNotificationPreferencesRepository{db: db}
}

func (r *postgresNotificationPreferencesRepository) GetByUserID(ctx context.Context, userID string) (*models.NotificationPreferences, error) {
	query, args, err := dialect.From(notificationPreferencesTable).
		Select("user_id", "task_assigned_email", "digest_cadence", "updated_at").
		Where(goqu.Ex{"user_id": userID}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}

	var p models.NotificationPreferences
	err = r.db.QueryRow(ctx, query, args...).Scan(&p.UserID, &p.TaskAssignedEmail, &p.DigestCadence, &p.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotificationPreferencesNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *postgresNotificationPreferencesRepository) Upsert(ctx context.Context, prefs *models.NotificationPreferences) error {
	query, args, err := dialect.Insert(notificationPreferencesTable).
		Rows(goqu.Record{
			"user_id":             prefs.UserID,
			"task_assigned_email": prefs.TaskAssignedEmail,
			"digest_cadence":      prefs.DigestCadence,
		}).
		OnConflict(goqu.DoUpdate("user_id", goqu.Record{
			"task_assigned_email": prefs.TaskAssignedEmail,
			"digest_cadence":      prefs.DigestCadence,
			"updated_at":          goqu.L("now()"),
		})).
		Returning("updated_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}
	return r.db.QueryRow(ctx, query, args...).Scan(&prefs.UpdatedAt)
}
