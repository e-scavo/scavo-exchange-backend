package repository

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/domain"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string) (*domain.UserSettings, error) {
	if userID == "" {
		return nil, domain.ErrUserIDRequired
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

	var settings domain.UserSettings
	var preferencesBytes []byte

	err := row.Scan(
		&settings.UserID,
		&preferencesBytes,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
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

func (r *PostgresRepository) UpsertPreferences(ctx context.Context, userID string, preferences map[string]any) (*domain.UserSettings, error) {
	if userID == "" {
		return nil, domain.ErrUserIDRequired
	}

	if preferences == nil {
		preferences = map[string]any{}
	}

	if r.pool == nil {
		return &domain.UserSettings{
			UserID:      userID,
			Preferences: preferences,
		}, nil
	}

	preferencesBytes, err := json.Marshal(preferences)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO user_settings (user_id, preferences)
		VALUES ($1, $2::jsonb)
		ON CONFLICT (user_id)
		DO UPDATE SET
			preferences = EXCLUDED.preferences,
			updated_at = NOW()
		RETURNING user_id, preferences, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query, userID, preferencesBytes)

	var settings domain.UserSettings
	var storedPreferencesBytes []byte

	if err := row.Scan(
		&settings.UserID,
		&storedPreferencesBytes,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if len(storedPreferencesBytes) > 0 {
		if err := json.Unmarshal(storedPreferencesBytes, &settings.Preferences); err != nil {
			return nil, err
		}
	} else {
		settings.Preferences = map[string]any{}
	}

	return &settings, nil
}
