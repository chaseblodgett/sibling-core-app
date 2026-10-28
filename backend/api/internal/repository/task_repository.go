package repository

import (
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

var ErrTaskNotFound = errors.New("task not found")

const tasksTable = "tasks"

// TaskRepository defines data access for the shared family task board.
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	// ListByFamily returns a family's tasks, optionally filtered to those
	// tied to a specific parent (nil parentID returns all of them).
	ListByFamily(ctx context.Context, familyID string, parentID *string) ([]*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id string) error
}

// postgresTaskRepository is the Postgres-backed implementation of TaskRepository.
type postgresTaskRepository struct {
	db *pgxpool.Pool
}

// NewTaskRepository constructs a Postgres-backed TaskRepository.
func NewTaskRepository(db *pgxpool.Pool) TaskRepository {
	return &postgresTaskRepository{db: db}
}

func (r *postgresTaskRepository) Create(ctx context.Context, task *models.Task) error {
	query, args, err := dialect.Insert(tasksTable).
		Rows(goqu.Record{
			"family_id":       task.FamilyID,
			"title":           task.Title,
			"description":     task.Description,
			"assignee_id":     task.AssigneeID,
			"parent_id":       task.ParentID,
			"created_by":      task.CreatedBy,
			"due_date":        task.DueDate,
			"recurrence_rule": task.RecurrenceRule,
		}).
		Returning("id", "status", "created_at", "updated_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	return r.db.QueryRow(ctx, query, args...).Scan(&task.ID, &task.Status, &task.CreatedAt, &task.UpdatedAt)
}

func (r *postgresTaskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	query, args, err := taskSelectDataset().Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanTask(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresTaskRepository) ListByFamily(ctx context.Context, familyID string, parentID *string) ([]*models.Task, error) {
	ds := taskSelectDataset().Where(goqu.Ex{"family_id": familyID})
	if parentID != nil {
		ds = ds.Where(goqu.Ex{"parent_id": *parentID})
	}
	query, args, err := ds.Order(goqu.I("created_at").Desc()).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*models.Task, 0)
	for rows.Next() {
		var t models.Task
		if err := scanTaskRow(rows, &t); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	return tasks, rows.Err()
}

func (r *postgresTaskRepository) Update(ctx context.Context, task *models.Task) error {
	query, args, err := dialect.Update(tasksTable).
		Set(goqu.Record{
			"title":           task.Title,
			"description":     task.Description,
			"status":          task.Status,
			"assignee_id":     task.AssigneeID,
			"parent_id":       task.ParentID,
			"due_date":        task.DueDate,
			"recurrence_rule": task.RecurrenceRule,
			"updated_at":      goqu.L("now()"),
		}).
		Where(goqu.Ex{"id": task.ID}).
		Returning("updated_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&task.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrTaskNotFound
	}
	return err
}

func (r *postgresTaskRepository) Delete(ctx context.Context, id string) error {
	query, args, err := dialect.Delete(tasksTable).Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func taskSelectDataset() *goqu.SelectDataset {
	return dialect.From(tasksTable).
		Select("id", "family_id", "title", "description", "status", "assignee_id", "parent_id", "created_by", "due_date", "recurrence_rule", "created_at", "updated_at")
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTaskRow(row taskScanner, t *models.Task) error {
	return row.Scan(&t.ID, &t.FamilyID, &t.Title, &t.Description, &t.Status, &t.AssigneeID, &t.ParentID, &t.CreatedBy, &t.DueDate, &t.RecurrenceRule, &t.CreatedAt, &t.UpdatedAt)
}

func scanTask(row pgx.Row) (*models.Task, error) {
	var t models.Task
	if err := scanTaskRow(row, &t); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &t, nil
}
