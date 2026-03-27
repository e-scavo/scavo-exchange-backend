package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletIdentityStorePG struct {
	db *pgxpool.Pool
}

func NewWalletIdentityStorePG(db *pgxpool.Pool) *WalletIdentityStorePG {
	return &WalletIdentityStorePG{db: db}
}

func (s *WalletIdentityStorePG) GetOrCreate(ctx context.Context, address string) (*WalletIdentity, error) {
	address = normalizeWalletAddress(address)
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}
	if s == nil || s.db == nil {
		return nil, ErrChallengeStore
	}

	identity, err := s.getByAddress(ctx, address)
	if err == nil {
		return identity, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	newID := uuid.NewString()

	_, err = s.db.Exec(ctx, `
		INSERT INTO auth_wallet_identities (id, address)
		VALUES ($1, $2)
		ON CONFLICT (address) DO NOTHING
	`, newID, address)
	if err != nil {
		return nil, err
	}

	identity, err = s.getByAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (s *WalletIdentityStorePG) GetByAddress(ctx context.Context, address string) (*WalletIdentity, error) {
	address = normalizeWalletAddress(address)
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}
	if s == nil || s.db == nil {
		return nil, ErrChallengeStore
	}

	identity, err := s.getByAddress(ctx, address)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletIdentityNotFound
		}
		return nil, err
	}

	return identity, nil
}

func (s *WalletIdentityStorePG) AttachUser(ctx context.Context, walletID, userID string, primary bool) (*WalletIdentity, error) {
	if s == nil || s.db == nil {
		return nil, ErrChallengeStore
	}

	walletID = strings.TrimSpace(walletID)
	userID = strings.TrimSpace(userID)
	if walletID == "" || userID == "" {
		return nil, ErrUnauthorized
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	current, err := scanWalletIdentityRow(tx.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE id = $1::uuid
		FOR UPDATE
	`, walletID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletIdentityNotFound
		}
		return nil, err
	}

	if strings.TrimSpace(current.UserID) != "" && strings.TrimSpace(current.UserID) != userID {
		return nil, ErrWalletIdentityAlreadyLinked
	}

	if primary {
		_, err = tx.Exec(ctx, `
			UPDATE auth_wallet_identities
			SET is_primary = FALSE
			WHERE user_id = $1 AND id <> $2::uuid
		`, userID, walletID)
		if err != nil {
			return nil, err
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE auth_wallet_identities
		SET
			user_id = $2,
			linked_at = COALESCE(linked_at, NOW()),
			is_primary = $3
		WHERE id = $1::uuid
	`, walletID, userID, primary)
	if err != nil {
		return nil, err
	}

	identity, err := scanWalletIdentityRow(tx.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE id = $1::uuid
	`, walletID))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return identity, nil
}

func (s *WalletIdentityStorePG) ReassignUser(ctx context.Context, walletID, fromUserID, toUserID string, primary bool) (*WalletIdentity, error) {
	if s == nil || s.db == nil {
		return nil, ErrChallengeStore
	}

	walletID = strings.TrimSpace(walletID)
	fromUserID = strings.TrimSpace(fromUserID)
	toUserID = strings.TrimSpace(toUserID)
	if walletID == "" || fromUserID == "" || toUserID == "" {
		return nil, ErrUnauthorized
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	current, err := scanWalletIdentityRow(tx.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE id = $1::uuid
		FOR UPDATE
	`, walletID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletIdentityNotFound
		}
		return nil, err
	}

	currentUserID := strings.TrimSpace(current.UserID)
	if currentUserID != fromUserID {
		if currentUserID == toUserID {
			return nil, ErrWalletMergeSameUser
		}
		return nil, ErrWalletIdentityAlreadyLinked
	}

	if primary {
		_, err = tx.Exec(ctx, `
			UPDATE auth_wallet_identities
			SET is_primary = FALSE
			WHERE user_id = $1 AND id <> $2::uuid
		`, toUserID, walletID)
		if err != nil {
			return nil, err
		}
	}

	_, err = tx.Exec(ctx, `
		UPDATE auth_wallet_identities
		SET
			user_id = $2,
			linked_at = COALESCE(linked_at, NOW()),
			is_primary = $3
		WHERE id = $1::uuid
	`, walletID, toUserID, primary)
	if err != nil {
		return nil, err
	}

	identity, err := scanWalletIdentityRow(tx.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE id = $1::uuid
	`, walletID))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return identity, nil
}

func (s *WalletIdentityStorePG) MergeUsers(ctx context.Context, sourceUserID, targetUserID string) ([]*WalletIdentity, error) {
	if s == nil || s.db == nil {
		return nil, ErrChallengeStore
	}

	sourceUserID = strings.TrimSpace(sourceUserID)
	targetUserID = strings.TrimSpace(targetUserID)
	if sourceUserID == "" || targetUserID == "" {
		return nil, ErrUnauthorized
	}
	if sourceUserID == targetUserID {
		return s.ListByUser(ctx, targetUserID)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var sourceCount int
	if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM auth_wallet_identities WHERE user_id = $1`, sourceUserID).Scan(&sourceCount); err != nil {
		return nil, err
	}
	if sourceCount == 0 {
		return []*WalletIdentity{}, nil
	}

	var targetHasPrimary bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM auth_wallet_identities WHERE user_id = $1 AND is_primary = TRUE)`, targetUserID).Scan(&targetHasPrimary); err != nil {
		return nil, err
	}

	_, err = tx.Exec(ctx, `
		UPDATE auth_wallet_identities
		SET
			user_id = $2,
			linked_at = COALESCE(linked_at, NOW()),
			is_primary = CASE WHEN $3 THEN FALSE ELSE is_primary END
		WHERE user_id = $1
	`, sourceUserID, targetUserID, targetHasPrimary)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE user_id = $1
		ORDER BY is_primary DESC, linked_at ASC NULLS LAST, address ASC
	`, targetUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*WalletIdentity, 0)
	for rows.Next() {
		identity, err := scanWalletIdentityRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, identity)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *WalletIdentityStorePG) SetPrimary(ctx context.Context, userID, address string) (*WalletIdentity, error) {
	if s == nil || s.db == nil {
		return nil, ErrChallengeStore
	}

	userID = strings.TrimSpace(userID)
	address = normalizeWalletAddress(address)
	if userID == "" {
		return nil, ErrUnauthorized
	}
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	current, err := scanWalletIdentityRow(tx.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE address = $1
		FOR UPDATE
	`, address))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWalletIdentityNotFound
		}
		return nil, err
	}
	if strings.TrimSpace(current.UserID) != userID {
		return nil, ErrWalletNotOwnedByUser
	}

	_, err = tx.Exec(ctx, `
		UPDATE auth_wallet_identities
		SET is_primary = CASE WHEN address = $2 THEN TRUE ELSE FALSE END
		WHERE user_id = $1
	`, userID, address)
	if err != nil {
		return nil, err
	}

	identity, err := scanWalletIdentityRow(tx.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE address = $1
	`, address))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return identity, nil
}

func (s *WalletIdentityStorePG) DetachUser(ctx context.Context, userID, address string) (*WalletIdentity, []*WalletIdentity, error) {
	if s == nil || s.db == nil {
		return nil, nil, ErrChallengeStore
	}

	userID = strings.TrimSpace(userID)
	address = normalizeWalletAddress(address)
	if userID == "" {
		return nil, nil, ErrUnauthorized
	}
	if !evmAddressRE.MatchString(address) {
		return nil, nil, ErrInvalidWalletAddress
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	current, err := scanWalletIdentityRow(tx.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE address = $1
		FOR UPDATE
	`, address))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrWalletIdentityNotFound
		}
		return nil, nil, err
	}
	if strings.TrimSpace(current.UserID) != userID {
		return nil, nil, ErrWalletNotOwnedByUser
	}

	_, err = tx.Exec(ctx, `
		UPDATE auth_wallet_identities
		SET
			user_id = NULL,
			linked_at = NULL,
			detached_at = NOW(),
			is_primary = FALSE
		WHERE address = $1
	`, address)
	if err != nil {
		return nil, nil, err
	}

	detached, err := scanWalletIdentityRow(tx.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE address = $1
	`, address))
	if err != nil {
		return nil, nil, err
	}

	rows, err := tx.Query(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE user_id = $1
		ORDER BY is_primary DESC, linked_at ASC NULLS LAST, address ASC
	`, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	remaining := make([]*WalletIdentity, 0)
	for rows.Next() {
		identity, err := scanWalletIdentityRows(rows)
		if err != nil {
			return nil, nil, err
		}
		remaining = append(remaining, identity)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return detached, remaining, nil
}

func (s *WalletIdentityStorePG) ListByUser(ctx context.Context, userID string) ([]*WalletIdentity, error) {
	if s == nil || s.db == nil {
		return []*WalletIdentity{}, nil
	}

	userID = strings.TrimSpace(userID)
	if userID == "" {
		return []*WalletIdentity{}, nil
	}

	rows, err := s.db.Query(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE user_id = $1
		ORDER BY is_primary DESC, linked_at ASC NULLS LAST, address ASC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*WalletIdentity, 0)
	for rows.Next() {
		identity, err := scanWalletIdentityRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, identity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (s *WalletIdentityStorePG) getByAddress(ctx context.Context, address string) (*WalletIdentity, error) {
	return scanWalletIdentityRow(s.db.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, ''), linked_at, detached_at, is_primary
		FROM auth_wallet_identities
		WHERE address = $1
	`, address))
}

func scanWalletIdentityRow(row pgx.Row) (*WalletIdentity, error) {
	var identity WalletIdentity
	var linkedAt pgtype.Timestamptz
	var detachedAt pgtype.Timestamptz

	err := row.Scan(
		&identity.ID,
		&identity.Address,
		&identity.UserID,
		&linkedAt,
		&detachedAt,
		&identity.IsPrimary,
	)
	if err != nil {
		return nil, err
	}

	if linkedAt.Valid {
		ts := linkedAt.Time.UTC()
		identity.LinkedAt = &ts
	}
	if detachedAt.Valid {
		ts := detachedAt.Time.UTC()
		identity.DetachedAt = &ts
	}

	return &identity, nil
}

func scanWalletIdentityRows(rows pgx.Rows) (*WalletIdentity, error) {
	var identity WalletIdentity
	var linkedAt pgtype.Timestamptz
	var detachedAt pgtype.Timestamptz

	err := rows.Scan(
		&identity.ID,
		&identity.Address,
		&identity.UserID,
		&linkedAt,
		&detachedAt,
		&identity.IsPrimary,
	)
	if err != nil {
		return nil, err
	}

	if linkedAt.Valid {
		ts := linkedAt.Time.UTC()
		identity.LinkedAt = &ts
	}
	if detachedAt.Valid {
		ts := detachedAt.Time.UTC()
		identity.DetachedAt = &ts
	}

	return &identity, nil
}
