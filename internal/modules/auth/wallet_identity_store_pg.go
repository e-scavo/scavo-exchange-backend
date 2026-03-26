package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (s *WalletIdentityStorePG) AttachUser(ctx context.Context, walletID, userID string) (*WalletIdentity, error) {
	if s == nil || s.db == nil {
		return nil, ErrChallengeStore
	}

	walletID = strings.TrimSpace(walletID)
	userID = strings.TrimSpace(userID)
	if walletID == "" || userID == "" {
		return nil, ErrUnauthorized
	}

	cmd, err := s.db.Exec(ctx, `
		UPDATE auth_wallet_identities
		SET user_id = $2
		WHERE id = $1::uuid
	`, walletID, userID)
	if err != nil {
		return nil, err
	}
	if cmd.RowsAffected() == 0 {
		return nil, ErrWalletIdentityNotFound
	}

	var identity WalletIdentity
	err = s.db.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, '')
		FROM auth_wallet_identities
		WHERE id = $1::uuid
	`, walletID).Scan(&identity.ID, &identity.Address, &identity.UserID)
	if err != nil {
		return nil, err
	}

	return &identity, nil
}

func (s *WalletIdentityStorePG) getByAddress(ctx context.Context, address string) (*WalletIdentity, error) {
	var identity WalletIdentity
	err := s.db.QueryRow(ctx, `
		SELECT id::text, address, COALESCE(user_id, '')
		FROM auth_wallet_identities
		WHERE address = $1
	`, address).Scan(&identity.ID, &identity.Address, &identity.UserID)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}
