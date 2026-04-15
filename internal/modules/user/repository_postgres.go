package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
	userrepo "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/repository"
)

var ErrUserNotFound = userrepo.ErrUserNotFound

type PostgresRepository = userrepo.PostgresRepository

func NewPostgresRepository(pool *pgxpool.Pool, log *logger.Logger) *PostgresRepository {
	return userrepo.NewPostgresRepository(pool, log)
}
