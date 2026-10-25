package repository

import (
	"context"
	"errors"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

var (
	ErrFamilyNotFound      = errors.New("family not found")
	ErrDuplicateInviteCode = errors.New("invite code already exists")
	ErrInviteNotFound      = errors.New("family invite not found")
)

const (
	familiesTable      = "families"
	familyInvitesTable = "family_invites"
)

// FamilyRepository defines data access for families and their invites.
type FamilyRepository interface {
	Create(ctx context.Context, family *models.Family) error
	GetByID(ctx context.Context, id string) (*models.Family, error)
	GetByInviteCode(ctx context.Context, code string) (*models.Family, error)
	Update(ctx context.Context, family *models.Family) error
	Delete(ctx context.Context, id string) error

	CreateInvite(ctx context.Context, invite *models.FamilyInvite) error
	GetInviteByID(ctx context.Context, id string) (*models.FamilyInvite, error)
	ListInvitesByFamily(ctx context.Context, familyID string) ([]*models.FamilyInvite, error)
	UpdateInviteStatus(ctx context.Context, id string, status string) error
}

// postgresFamilyRepository is the Postgres-backed implementation of FamilyRepository.
type postgresFamilyRepository struct {
	db *pgxpool.Pool
}

// NewFamilyRepository constructs a Postgres-backed FamilyRepository.
func NewFamilyRepository(db *pgxpool.Pool) FamilyRepository {
	return &postgresFamilyRepository{db: db}
}

func (r *postgresFamilyRepository) Create(ctx context.Context, family *models.Family) error {
	query, args, err := dialect.Insert(familiesTable).
		Rows(goqu.Record{
			"name":        family.Name,
			"invite_code": family.InviteCode,
			"created_by":  family.CreatedBy,
		}).
		Returning("id", "created_at", "updated_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&family.ID, &family.CreatedAt, &family.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateInviteCode
		}
		return err
	}
	return nil
}

func (r *postgresFamilyRepository) GetByID(ctx context.Context, id string) (*models.Family, error) {
	query, args, err := familySelectDataset().Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanFamily(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresFamilyRepository) GetByInviteCode(ctx context.Context, code string) (*models.Family, error) {
	query, args, err := familySelectDataset().Where(goqu.Ex{"invite_code": code}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanFamily(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresFamilyRepository) Update(ctx context.Context, family *models.Family) error {
	query, args, err := dialect.Update(familiesTable).
		Set(goqu.Record{"name": family.Name, "updated_at": goqu.L("now()")}).
		Where(goqu.Ex{"id": family.ID}).
		Returning("updated_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&family.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrFamilyNotFound
	}
	return err
}

func (r *postgresFamilyRepository) Delete(ctx context.Context, id string) error {
	query, args, err := dialect.Delete(familiesTable).Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrFamilyNotFound
	}
	return nil
}

func (r *postgresFamilyRepository) CreateInvite(ctx context.Context, invite *models.FamilyInvite) error {
	query, args, err := dialect.Insert(familyInvitesTable).
		Rows(goqu.Record{
			"family_id":  invite.FamilyID,
			"email":      invite.Email,
			"invited_by": invite.InvitedBy,
			"status":     invite.Status,
		}).
		Returning("id", "created_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	return r.db.QueryRow(ctx, query, args...).Scan(&invite.ID, &invite.CreatedAt)
}

func (r *postgresFamilyRepository) GetInviteByID(ctx context.Context, id string) (*models.FamilyInvite, error) {
	query, args, err := familyInviteSelectDataset().Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanFamilyInvite(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresFamilyRepository) ListInvitesByFamily(ctx context.Context, familyID string) ([]*models.FamilyInvite, error) {
	query, args, err := familyInviteSelectDataset().
		Where(goqu.Ex{"family_id": familyID}).
		Order(goqu.I("created_at").Desc()).
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

	invites := make([]*models.FamilyInvite, 0)
	for rows.Next() {
		invite, err := scanFamilyInviteRow(rows)
		if err != nil {
			return nil, err
		}
		invites = append(invites, invite)
	}
	return invites, rows.Err()
}

func (r *postgresFamilyRepository) UpdateInviteStatus(ctx context.Context, id string, status string) error {
	values := goqu.Record{"status": status}
	if status == "accepted" {
		values["accepted_at"] = goqu.L("now()")
	}

	query, args, err := dialect.Update(familyInvitesTable).
		Set(values).
		Where(goqu.Ex{"id": id}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrInviteNotFound
	}
	return nil
}

func familySelectDataset() *goqu.SelectDataset {
	return dialect.From(familiesTable).
		Select("id", "name", "invite_code", "created_by", "created_at", "updated_at")
}

func scanFamily(row pgx.Row) (*models.Family, error) {
	var f models.Family
	if err := row.Scan(&f.ID, &f.Name, &f.InviteCode, &f.CreatedBy, &f.CreatedAt, &f.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFamilyNotFound
		}
		return nil, err
	}
	return &f, nil
}

func familyInviteSelectDataset() *goqu.SelectDataset {
	return dialect.From(familyInvitesTable).
		Select("id", "family_id", "email", "invited_by", "status", "created_at", "accepted_at")
}

func scanFamilyInvite(row pgx.Row) (*models.FamilyInvite, error) {
	var inv models.FamilyInvite
	if err := row.Scan(&inv.ID, &inv.FamilyID, &inv.Email, &inv.InvitedBy, &inv.Status, &inv.CreatedAt, &inv.AcceptedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInviteNotFound
		}
		return nil, err
	}
	return &inv, nil
}

func scanFamilyInviteRow(rows pgx.Rows) (*models.FamilyInvite, error) {
	var inv models.FamilyInvite
	if err := rows.Scan(&inv.ID, &inv.FamilyID, &inv.Email, &inv.InvitedBy, &inv.Status, &inv.CreatedAt, &inv.AcceptedAt); err != nil {
		return nil, err
	}
	return &inv, nil
}
