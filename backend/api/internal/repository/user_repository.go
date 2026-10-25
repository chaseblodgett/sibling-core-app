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
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("email already exists")
)

const usersTable = "users"

// UserRepository defines data access for user accounts.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
}

// postgresUserRepository is the Postgres-backed implementation of UserRepository.
type postgresUserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository constructs a Postgres-backed UserRepository.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *models.User) error {
	query, args, err := dialect.Insert(usersTable).
		Rows(goqu.Record{
			"email":         user.Email,
			"password_hash": user.PasswordHash,
			"name":          user.Name,
		}).
		Returning("id", "created_at").
		Prepared(true).
		ToSQL()
	if err != nil {
		return err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	query, args, err := userSelectDataset().Where(goqu.Ex{"email": email}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanUser(r.db.QueryRow(ctx, query, args...))
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	query, args, err := userSelectDataset().Where(goqu.Ex{"id": id}).Prepared(true).ToSQL()
	if err != nil {
		return nil, err
	}
	return scanUser(r.db.QueryRow(ctx, query, args...))
}

func userSelectDataset() *goqu.SelectDataset {
	return dialect.From(usersTable).Select("id", "email", "password_hash", "name", "created_at")
}

func scanUser(row pgx.Row) (*models.User, error) {
	var u models.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}
