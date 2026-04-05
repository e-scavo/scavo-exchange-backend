package usersettings

import (
	"context"
	"database/sql"
	"encoding/json"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*UserSettings, error) {
	if userID == "" {
		return nil, ErrUserIDRequired
	}

	query := `
		SELECT user_id, preferences, created_at, updated_at
		FROM user_settings
		WHERE user_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, userID)

	var settings UserSettings
	var preferencesBytes []byte

	err := row.Scan(
		&settings.UserID,
		&preferencesBytes,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
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
