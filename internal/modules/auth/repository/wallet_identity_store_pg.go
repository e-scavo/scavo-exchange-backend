package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"

	rootauth "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"
)

type WalletIdentityStorePG = rootauth.WalletIdentityStorePG

func NewWalletIdentityStorePG(db *pgxpool.Pool) *WalletIdentityStorePG {
	return rootauth.NewWalletIdentityStorePG(db)
}
