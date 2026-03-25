package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletChallengeStorePG struct {
	db *pgxpool.Pool
}

func NewWalletChallengeStorePG(db *pgxpool.Pool) *WalletChallengeStorePG {
	return &WalletChallengeStorePG{db: db}
}

func (s *WalletChallengeStorePG) Save(ctx context.Context, challenge *WalletChallenge) error {
	if s == nil || s.db == nil {
		return ErrChallengeStore
	}
	if challenge == nil || challenge.ID == "" {
		return errors.New("invalid wallet challenge")
	}

	_, err := s.db.Exec(ctx, `
		INSERT INTO auth_wallet_challenges (
			id,
			address,
			chain,
			nonce,
			message,
			issued_at,
			expires_at,
			used_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`,
		challenge.ID,
		challenge.Address,
		challenge.Chain,
		challenge.Nonce,
		challenge.Message,
		challenge.IssuedAt.UTC(),
		challenge.ExpiresAt.UTC(),
		challenge.UsedAt,
	)
	return err
}

func (s *WalletChallengeStorePG) GetByID(ctx context.Context, id string) (*WalletChallenge, error) {
	if s == nil || s.db == nil {
		return nil, ErrChallengeStore
	}

	row := s.db.QueryRow(ctx, `
		SELECT
			id::text,
			address,
			chain,
			nonce,
			message,
			issued_at,
			expires_at,
			used_at
		FROM auth_wallet_challenges
		WHERE id = $1
	`, id)

	var challenge WalletChallenge
	var usedAt *time.Time

	err := row.Scan(
		&challenge.ID,
		&challenge.Address,
		&challenge.Chain,
		&challenge.Nonce,
		&challenge.Message,
		&challenge.IssuedAt,
		&challenge.ExpiresAt,
		&usedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrWalletChallengeNotFound
	}
	if err != nil {
		return nil, err
	}

	if usedAt != nil {
		ts := usedAt.UTC()
		challenge.UsedAt = &ts
	}

	challenge.IssuedAt = challenge.IssuedAt.UTC()
	challenge.ExpiresAt = challenge.ExpiresAt.UTC()

	return &challenge, nil
}

func (s *WalletChallengeStorePG) Use(ctx context.Context, id string, usedAt time.Time) (*WalletChallenge, error) {
	if s == nil || s.db == nil {
		return nil, ErrChallengeStore
	}

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	row := tx.QueryRow(ctx, `
		SELECT
			id::text,
			address,
			chain,
			nonce,
			message,
			issued_at,
			expires_at,
			used_at
		FROM auth_wallet_challenges
		WHERE id = $1
		FOR UPDATE
	`, id)

	var challenge WalletChallenge
	var currentUsedAt *time.Time

	err = row.Scan(
		&challenge.ID,
		&challenge.Address,
		&challenge.Chain,
		&challenge.Nonce,
		&challenge.Message,
		&challenge.IssuedAt,
		&challenge.ExpiresAt,
		&currentUsedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrWalletChallengeNotFound
	}
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if currentUsedAt != nil {
		return nil, ErrChallengeUsed
	}
	if now.After(challenge.ExpiresAt.UTC()) {
		return nil, ErrChallengeExpired
	}

	ts := usedAt.UTC()
	_, err = tx.Exec(ctx, `
		UPDATE auth_wallet_challenges
		SET used_at = $2
		WHERE id = $1
	`, id, ts)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	challenge.IssuedAt = challenge.IssuedAt.UTC()
	challenge.ExpiresAt = challenge.ExpiresAt.UTC()
	challenge.UsedAt = &ts

	return &challenge, nil
}
