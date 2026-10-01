package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

// DocumentRepository defines data access for document metadata.
// The underlying file bytes live in Supabase Storage; this stores the
// pointer/metadata row tied to a family.
type DocumentRepository interface {
	Create(ctx context.Context, document *models.Document) error
	GetByID(ctx context.Context, id string) (*models.Document, error)
	ListByFamily(ctx context.Context, familyID string) ([]*models.Document, error)
	Delete(ctx context.Context, id string) error
}

// postgresDocumentRepository is the Postgres-backed implementation of DocumentRepository.
type postgresDocumentRepository struct {
	db *pgxpool.Pool
}

// NewDocumentRepository constructs a Postgres-backed DocumentRepository.
func NewDocumentRepository(db *pgxpool.Pool) DocumentRepository {
	return &postgresDocumentRepository{db: db}
}

func (r *postgresDocumentRepository) Create(ctx context.Context, document *models.Document) error {
	return errors.New("not implemented")
}

func (r *postgresDocumentRepository) GetByID(ctx context.Context, id string) (*models.Document, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresDocumentRepository) ListByFamily(ctx context.Context, familyID string) ([]*models.Document, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresDocumentRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
