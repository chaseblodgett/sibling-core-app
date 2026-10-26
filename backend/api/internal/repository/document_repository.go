package repository

import (
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

var ErrDocumentNotFound = errors.New("document not found")

const documentsTable = "documents"

// DocumentRepository defines data access for document metadata.
// The underlying file bytes live in Supabase Storage; this stores the
// pointer/metadata row tied to a family.
type DocumentRepository interface {
	Create(ctx context.Context, document *models.Document) error
	GetByID(ctx context.Context, id string) (*models.Document, error)
	// ListByFamily returns a family's documents, optionally filtered to
	// those tied to a specific parent (nil parentID returns all of them).
	ListByFamily(ctx context.Context, familyID string, parentID *string) ([]*models.Document, error)
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
	query, args, err := dialect.Insert(documentsTable).
		Rows(goqu.Record{
			"family_id":    document.FamilyID,
			"uploaded_by":  document.UploadedBy,
			"file_name":    document.FileName,
			"storage_path": document.StoragePath,
			"content_type": document.ContentType,
			"size_bytes":   document.SizeBytes,
			"parent_id":    document.ParentID,
		}).
		Returning("id", "created_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}
	return r.db.QueryRow(ctx, query, args...).Scan(&document.ID, &document.CreatedAt)
}

func (r *postgresDocumentRepository) GetByID(ctx context.Context, id string) (*models.Document, error) {
	query, args, err := documentSelectDataset().Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanDocument(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresDocumentRepository) ListByFamily(ctx context.Context, familyID string, parentID *string) ([]*models.Document, error) {
	ds := documentSelectDataset().Where(goqu.Ex{"family_id": familyID})
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

	documents := make([]*models.Document, 0)
	for rows.Next() {
		var d models.Document
		if err := scanDocumentRow(rows, &d); err != nil {
			return nil, err
		}
		documents = append(documents, &d)
	}
	return documents, rows.Err()
}

func (r *postgresDocumentRepository) Delete(ctx context.Context, id string) error {
	query, args, err := dialect.Delete(documentsTable).Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrDocumentNotFound
	}
	return nil
}

func documentSelectDataset() *goqu.SelectDataset {
	return dialect.From(documentsTable).
		Select("id", "family_id", "uploaded_by", "file_name", "storage_path", "content_type", "size_bytes", "parent_id", "created_at")
}

func scanDocumentRow(row taskScanner, d *models.Document) error {
	return row.Scan(&d.ID, &d.FamilyID, &d.UploadedBy, &d.FileName, &d.StoragePath, &d.ContentType, &d.SizeBytes, &d.ParentID, &d.CreatedAt)
}

func scanDocument(row pgx.Row) (*models.Document, error) {
	var d models.Document
	if err := scanDocumentRow(row, &d); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}
	return &d, nil
}
