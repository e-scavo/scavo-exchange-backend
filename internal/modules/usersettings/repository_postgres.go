package usersettings

import (
	"github.com/jackc/pgx/v5/pgxpool"

	usersettingsrepo "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/repository"
)

type PostgresRepository = usersettingsrepo.PostgresRepository

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return usersettingsrepo.NewPostgresRepository(pool)
}
