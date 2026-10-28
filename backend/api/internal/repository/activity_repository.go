package repository

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

const activityLogTable = "activity_log"

// ActivityRepository defines data access for a family's activity feed.
type ActivityRepository interface {
	Create(ctx context.Context, entry *models.ActivityEntry) error
	// ListByFamily returns a family's activity entries, newest first.
	ListByFamily(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error)
}

// postgresActivityRepository is the Postgres-backed implementation of ActivityRepository.
type postgresActivityRepository struct {
	db *pgxpool.Pool
}

// NewActivityRepository constructs a Postgres-backed ActivityRepository.
func NewActivityRepository(db *pgxpool.Pool) ActivityRepository {
	return &postgresActivityRepository{db: db}
}

func (r *postgresActivityRepository) Create(ctx context.Context, entry *models.ActivityEntry) error {
	query, args, err := dialect.Insert(activityLogTable).
		Rows(goqu.Record{
			"family_id":     entry.FamilyID,
			"actor_id":      entry.ActorID,
			"actor_name":    entry.ActorName,
			"verb":          entry.Verb,
			"subject_type":  entry.SubjectType,
			"subject_id":    entry.SubjectID,
			"subject_label": entry.SubjectLabel,
		}).
		Returning("id", "created_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	return r.db.QueryRow(ctx, query, args...).Scan(&entry.ID, &entry.CreatedAt)
}

func (r *postgresActivityRepository) ListByFamily(ctx context.Context, familyID string, limit, offset int) ([]*models.ActivityEntry, error) {
	query, args, err := dialect.From(activityLogTable).
		Select("id", "family_id", "actor_id", "actor_name", "verb", "subject_type", "subject_id", "subject_label", "created_at").
		Where(goqu.Ex{"family_id": familyID}).
		Order(goqu.I("created_at").Desc()).
		Limit(uint(limit)).
		Offset(uint(offset)).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]*models.ActivityEntry, 0)
	for rows.Next() {
		var e models.ActivityEntry
		if err := rows.Scan(
			&e.ID, &e.FamilyID, &e.ActorID, &e.ActorName, &e.Verb,
			&e.SubjectType, &e.SubjectID, &e.SubjectLabel, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, &e)
	}
	return entries, rows.Err()
}
