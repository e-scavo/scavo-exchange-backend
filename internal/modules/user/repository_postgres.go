package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
)

var ErrUserNotFound = errors.New("user not found")

type PostgresRepository struct {
	pool *pgxpool.Pool
	log  *logger.Logger
}

func NewPostgresRepository(pool *pgxpool.Pool, log *logger.Logger) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
		log:  log,
	}
}

func (r *PostgresRepository) UpsertDevUser(ctx context.Context, email string) (*User, error) {
	const q = `
INSERT INTO users (
    id,
    email,
    display_name,
    created_at,
    updated_at,
    last_login_at
)
VALUES (
    $1,
    $2,
    '',
    NOW(),
    NOW(),
    NOW()
)
ON CONFLICT (email)
DO UPDATE SET
    updated_at = NOW(),
    last_login_at = NOW()
RETURNING
    id,
    email,
    display_name,
    created_at,
    updated_at,
    last_login_at
`

	u := &User{}
	err := r.pool.QueryRow(ctx, q, devUserID(email), email).Scan(
		&u.ID,
		&u.Email,
		&u.DisplayName,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *PostgresRepository) UpsertWalletUser(ctx context.Context, id, email, displayName string) (*User, error) {
	const q = `
INSERT INTO users (
    id,
    email,
    display_name,
    created_at,
    updated_at,
    last_login_at
)
VALUES (
    $1,
    $2,
    $3,
    NOW(),
    NOW(),
    NOW()
)
ON CONFLICT (id)
DO UPDATE SET
    email = EXCLUDED.email,
    display_name = EXCLUDED.display_name,
    updated_at = NOW(),
    last_login_at = NOW()
RETURNING
    id,
    email,
    display_name,
    created_at,
    updated_at,
    last_login_at
`

	u := &User{}
	err := r.pool.QueryRow(ctx, q, id, email, displayName).Scan(
		&u.ID,
		&u.Email,
		&u.DisplayName,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLoginAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*User, error) {
	const q = `
SELECT
    id,
    email,
    display_name,
    created_at,
    updated_at,
    last_login_at
FROM users
WHERE id = $1
`

	u := &User{}
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&u.ID,
		&u.Email,
		&u.DisplayName,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLoginAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *PostgresRepository) UpdateDisplayName(ctx context.Context, id, displayName string) (*User, error) {
	const q = `
UPDATE users
SET
    display_name = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    email,
    display_name,
    created_at,
    updated_at,
    last_login_at
`

	u := &User{}
	err := r.pool.QueryRow(ctx, q, id, displayName).Scan(
		&u.ID,
		&u.Email,
		&u.DisplayName,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLoginAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}
