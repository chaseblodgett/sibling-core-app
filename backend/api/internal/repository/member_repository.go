package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"sibling-core/api/internal/models"
)

// MemberRepository defines data access for family membership records.
type MemberRepository interface {
	Create(ctx context.Context, member *models.Member) error
	GetByID(ctx context.Context, id string) (*models.Member, error)
	GetByUserAndFamily(ctx context.Context, userID, familyID string) (*models.Member, error)
	ListByFamily(ctx context.Context, familyID string) ([]*models.Member, error)
	Delete(ctx context.Context, id string) error
}

// postgresMemberRepository is the Postgres-backed implementation of MemberRepository.
type postgresMemberRepository struct {
	db *pgxpool.Pool
}

// NewMemberRepository constructs a Postgres-backed MemberRepository.
func NewMemberRepository(db *pgxpool.Pool) MemberRepository {
	return &postgresMemberRepository{db: db}
}

func (r *postgresMemberRepository) Create(ctx context.Context, member *models.Member) error {
	return errors.New("not implemented")
}

func (r *postgresMemberRepository) GetByID(ctx context.Context, id string) (*models.Member, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresMemberRepository) GetByUserAndFamily(ctx context.Context, userID, familyID string) (*models.Member, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresMemberRepository) ListByFamily(ctx context.Context, familyID string) ([]*models.Member, error) {
	return nil, errors.New("not implemented")
}

func (r *postgresMemberRepository) Delete(ctx context.Context, id string) error {
	return errors.New("not implemented")
}
