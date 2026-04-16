package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"

	rootauth "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"
)

type WalletChallengeStorePG = rootauth.WalletChallengeStorePG

func NewWalletChallengeStorePG(db *pgxpool.Pool) *WalletChallengeStorePG {
	return rootauth.NewWalletChallengeStorePG(db)
}
