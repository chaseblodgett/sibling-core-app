package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

// TaskRepository defines data access for the shared family task board.
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id string) (*models.Task, error)
	ListByFamily(ctx context.Context, familyID string) ([]*models.Task, error)
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
	return errors.New("not implemented")
}

func (r *postgresTaskRepository) GetByID(ctx context.Context, id string) (*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresTaskRepository) ListByFamily(ctx context.Context, familyID string) ([]*models.Task, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresTaskRepository) Update(ctx context.Context, task *models.Task) error {
	return errors.New("not implemented")
}

func (r *postgresTaskRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
