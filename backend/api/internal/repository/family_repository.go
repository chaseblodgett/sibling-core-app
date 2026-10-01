package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
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
	return errors.New("not implemented")
}

func (r *postgresFamilyRepository) GetByID(ctx context.Context, id string) (*models.Family, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresFamilyRepository) GetByInviteCode(ctx context.Context, code string) (*models.Family, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresFamilyRepository) Update(ctx context.Context, family *models.Family) error {
	return errors.New("not implemented")
}

func (r *postgresFamilyRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func (r *postgresFamilyRepository) CreateInvite(ctx context.Context, invite *models.FamilyInvite) error {
	return errors.New("not implemented")
}

func (r *postgresFamilyRepository) GetInviteByID(ctx context.Context, id string) (*models.FamilyInvite, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresFamilyRepository) ListInvitesByFamily(ctx context.Context, familyID string) ([]*models.FamilyInvite, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresFamilyRepository) UpdateInviteStatus(ctx context.Context, id string, status string) error {
	return errors.New("not implemented")
}
