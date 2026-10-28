package repository

import (
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

var ErrParentNotFound = errors.New("parent not found")

const parentsTable = "parents"

// ParentRepository defines data access for a family's tracked parents
// (care-recipients).
type ParentRepository interface {
	Create(ctx context.Context, parent *models.Parent) error
	GetByID(ctx context.Context, id string) (*models.Parent, error)
	ListByFamily(ctx context.Context, familyID string) ([]*models.Parent, error)
	Update(ctx context.Context, parent *models.Parent) error
	Delete(ctx context.Context, id string) error
}

// postgresParentRepository is the Postgres-backed implementation of ParentRepository.
type postgresParentRepository struct {
	db *pgxpool.Pool
}

// NewParentRepository constructs a Postgres-backed ParentRepository.
func NewParentRepository(db *pgxpool.Pool) ParentRepository {
	return &postgresParentRepository{db: db}
}

func (r *postgresParentRepository) Create(ctx context.Context, parent *models.Parent) error {
	query, args, err := dialect.Insert(parentsTable).
		Rows(goqu.Record{
			"family_id":  parent.FamilyID,
			"name":       parent.Name,
			"relation":   parent.Relation,
			"created_by": parent.CreatedBy,
		}).
		Returning("id", "created_at", "updated_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	return r.db.QueryRow(ctx, query, args...).Scan(&parent.ID, &parent.CreatedAt, &parent.UpdatedAt)
}

func (r *postgresParentRepository) GetByID(ctx context.Context, id string) (*models.Parent, error) {
	query, args, err := parentSelectDataset().Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanParent(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresParentRepository) ListByFamily(ctx context.Context, familyID string) ([]*models.Parent, error) {
	query, args, err := parentSelectDataset().
		Where(goqu.Ex{"family_id": familyID}).
		Order(goqu.I("created_at").Asc()).
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

	parents := make([]*models.Parent, 0)
	for rows.Next() {
		var p models.Parent
		if err := scanParentRow(rows, &p); err != nil {
			return nil, err
		}
		parents = append(parents, &p)
	}
	return parents, rows.Err()
}

func (r *postgresParentRepository) Update(ctx context.Context, parent *models.Parent) error {
	query, args, err := dialect.Update(parentsTable).
		Set(goqu.Record{
			"name":                    parent.Name,
			"relation":                parent.Relation,
			"doctor_name":             parent.DoctorName,
			"doctor_phone":            parent.DoctorPhone,
			"allergies":               parent.Allergies,
			"medications":             parent.Medications,
			"insurance_provider":      parent.InsuranceProvider,
			"insurance_member_id":     parent.InsuranceMemberID,
			"emergency_contact_name":  parent.EmergencyContactName,
			"emergency_contact_phone": parent.EmergencyContactPhone,
			"updated_at":              goqu.L("now()"),
		}).
		Where(goqu.Ex{"id": parent.ID}).
		Returning("updated_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&parent.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrParentNotFound
	}
	return err
}

func (r *postgresParentRepository) Delete(ctx context.Context, id string) error {
	query, args, err := dialect.Delete(parentsTable).Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrParentNotFound
	}
	return nil
}

func parentSelectDataset() *goqu.SelectDataset {
	return dialect.From(parentsTable).
		Select(
			"id", "family_id", "name", "relation",
			"doctor_name", "doctor_phone", "allergies", "medications",
			"insurance_provider", "insurance_member_id",
			"emergency_contact_name", "emergency_contact_phone",
			"created_by", "created_at", "updated_at",
		)
}

type parentScanner interface {
	Scan(dest ...any) error
}

func scanParentRow(row parentScanner, p *models.Parent) error {
	return row.Scan(
		&p.ID, &p.FamilyID, &p.Name, &p.Relation,
		&p.DoctorName, &p.DoctorPhone, &p.Allergies, &p.Medications,
		&p.InsuranceProvider, &p.InsuranceMemberID,
		&p.EmergencyContactName, &p.EmergencyContactPhone,
		&p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
	)
}

func scanParent(row pgx.Row) (*models.Parent, error) {
	var p models.Parent
	if err := scanParentRow(row, &p); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrParentNotFound
		}
		return nil, err
	}
	return &p, nil
}
