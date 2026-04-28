package readmodels

import authdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/domain"

// NewAuthLoginReadModel builds the explicit login output projection from
// resolved authentication values. It keeps login response construction outside
// HTTP handlers while avoiding any dependency on transport-layer structs.
func NewAuthLoginReadModel(accessToken, tokenType string, expiresIn int64, userID string) AuthLoginReadModel {
	return AuthLoginReadModel{
		AccessToken: accessToken,
		TokenType:   tokenType,
		ExpiresIn:   expiresIn,
		UserID:      userID,
	}
}

// FromWalletIdentity maps a canonical wallet identity domain model into its
// output-only auth wallet read model projection.
func FromWalletIdentity(wallet *authdomain.WalletIdentity) *AuthWalletReadModel {
	if wallet == nil {
		return nil
	}

	return &AuthWalletReadModel{
		ID:                 wallet.ID,
		Address:            wallet.Address,
		UserID:             wallet.UserID,
		LinkedAt:           wallet.LinkedAt,
		DetachedAt:         wallet.DetachedAt,
		IsPrimary:          wallet.IsPrimary,
		Status:             walletStatus(wallet),
		CanSetPrimary:      false,
		CanDetach:          false,
		DetachBlockReasons: []string{},
	}
}

// FromWalletIdentities maps a wallet identity collection into output-only read
// models. Nil inputs return an empty slice to preserve response stability.
func FromWalletIdentities(wallets []*authdomain.WalletIdentity) []*AuthWalletReadModel {
	if wallets == nil {
		return []*AuthWalletReadModel{}
	}

	out := make([]*AuthWalletReadModel, 0, len(wallets))
	for _, wallet := range wallets {
		mapped := FromWalletIdentity(wallet)
		if mapped != nil {
			out = append(out, mapped)
		}
	}
	return out
}

// FromWalletChallenge maps a canonical wallet challenge domain model into its
// output-only auth wallet challenge projection.
func FromWalletChallenge(challenge *authdomain.WalletChallenge) *AuthWalletChallengeReadModel {
	if challenge == nil {
		return nil
	}

	return &AuthWalletChallengeReadModel{
		ID:                challenge.ID,
		Address:           challenge.Address,
		Chain:             challenge.Chain,
		Nonce:             challenge.Nonce,
		Message:           challenge.Message,
		Purpose:           challenge.Purpose,
		RequestedByUserID: challenge.RequestedByUserID,
		IssuedAt:          challenge.IssuedAt,
		ExpiresAt:         challenge.ExpiresAt,
		UsedAt:            challenge.UsedAt,
	}
}

func walletStatus(wallet *authdomain.WalletIdentity) string {
	if wallet == nil {
		return "unknown"
	}
	if wallet.UserID != "" {
		return "active"
	}
	if wallet.DetachedAt != nil {
		return "detached"
	}
	return "unlinked"
}
