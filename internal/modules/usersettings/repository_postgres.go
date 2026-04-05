package usersettings

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*UserSettings, error) {
	if userID == "" {
		return nil, ErrUserIDRequired
	}

	if r.pool == nil {
		// DB no configurada → comportamiento consistente con el resto del proyecto
		return nil, nil
	}

	query := `
		SELECT user_id, preferences, created_at, updated_at
		FROM user_settings
		WHERE user_id = $1
	`

	row := r.pool.QueryRow(ctx, query, userID)

	var settings UserSettings
	var preferencesBytes []byte

	err := row.Scan(
		&settings.UserID,
		&preferencesBytes,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		// no row → defaults se resuelven en service
		return nil, nil
	}

	if len(preferencesBytes) > 0 {
		if err := json.Unmarshal(preferencesBytes, &settings.Preferences); err != nil {
			return nil, err
		}
	} else {
		settings.Preferences = map[string]any{}
	}

	return &settings, nil
}
