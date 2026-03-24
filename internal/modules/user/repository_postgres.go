package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
)

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
