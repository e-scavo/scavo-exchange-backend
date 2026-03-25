package auth

import (
	"context"
	"errors"

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

	var id string
	err := s.db.QueryRow(ctx, `
		SELECT id::text
		FROM auth_wallet_identities
		WHERE address = $1
	`, address).Scan(&id)
	if err == nil {
		return &WalletIdentity{
			ID:      id,
			Address: address,
		}, nil
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

	err = s.db.QueryRow(ctx, `
		SELECT id::text
		FROM auth_wallet_identities
		WHERE address = $1
	`, address).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &WalletIdentity{
		ID:      id,
		Address: address,
	}, nil
}
