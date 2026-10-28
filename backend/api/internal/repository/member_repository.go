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
	ErrMemberNotFound = errors.New("member not found")
	ErrAlreadyMember  = errors.New("user is already a member of this family")
)

const membersTable = "members"

// MemberRepository defines data access for family membership records.
type MemberRepository interface {
	Create(ctx context.Context, member *models.Member) error
	GetByID(ctx context.Context, id string) (*models.Member, error)
	GetByUserAndFamily(ctx context.Context, userID, familyID string) (*models.Member, error)
	ListByFamily(ctx context.Context, familyID string) ([]*models.Member, error)
	// ListByUser returns every family membership for a user, used to let a
	// logged-in user discover which family/families they belong to.
	ListByUser(ctx context.Context, userID string) ([]*models.Member, error)
	// SetRole changes a member's role (owner, member, or viewer).
	SetRole(ctx context.Context, id, role string) error
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
	query, args, err := dialect.Insert(membersTable).
		Rows(goqu.Record{
			"family_id": member.FamilyID,
			"user_id":   member.UserID,
			"email":     member.Email,
			"full_name": member.FullName,
			"role":      member.Role,
		}).
		Returning("id", "joined_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&member.ID, &member.JoinedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyMember
		}
		return err
	}
	return nil
}

func (r *postgresMemberRepository) GetByID(ctx context.Context, id string) (*models.Member, error) {
	query, args, err := memberSelectDataset().Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanMember(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresMemberRepository) GetByUserAndFamily(ctx context.Context, userID, familyID string) (*models.Member, error) {
	query, args, err := memberSelectDataset().
		Where(goqu.Ex{"user_id": userID, "family_id": familyID}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}
	return scanMember(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresMemberRepository) ListByFamily(ctx context.Context, familyID string) ([]*models.Member, error) {
	query, args, err := memberSelectDataset().
		Where(goqu.Ex{"family_id": familyID}).
		Order(goqu.I("joined_at").Asc()).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}
	return queryMembers(ctx, r.db, query, args...)
}

func (r *postgresMemberRepository) ListByUser(ctx context.Context, userID string) ([]*models.Member, error) {
	query, args, err := memberSelectDataset().
		Where(goqu.Ex{"user_id": userID}).
		Order(goqu.I("joined_at").Asc()).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, err
	}
	return queryMembers(ctx, r.db, query, args...)
}

func (r *postgresMemberRepository) SetRole(ctx context.Context, id, role string) error {
	query, args, err := dialect.Update(membersTable).
		Set(goqu.Record{"role": role}).
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
		return ErrMemberNotFound
	}
	return nil
}

func (r *postgresMemberRepository) Delete(ctx context.Context, id string) error {
	query, args, err := dialect.Delete(membersTable).Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return err
	}
	tag, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrMemberNotFound
	}
	return nil
}

func memberSelectDataset() *goqu.SelectDataset {
	return dialect.From(membersTable).
		Select("id", "family_id", "user_id", "email", "full_name", "role", "joined_at")
}

func scanMember(row pgx.Row) (*models.Member, error) {
	var m models.Member
	if err := row.Scan(&m.ID, &m.FamilyID, &m.UserID, &m.Email, &m.FullName, &m.Role, &m.JoinedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}
	return &m, nil
}

func queryMembers(ctx context.Context, db *pgxpool.Pool, query string, args ...any) ([]*models.Member, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]*models.Member, 0)
	for rows.Next() {
		var m models.Member
		if err := rows.Scan(&m.ID, &m.FamilyID, &m.UserID, &m.Email, &m.FullName, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	return members, rows.Err()
}
