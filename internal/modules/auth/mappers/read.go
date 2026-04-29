// Package mappers centralizes explicit model transformations for the auth module.
//
// It is introduced in Phase 0.12.4 as the module-level mapping boundary between
// write models, domain inputs, domain state and read models. Existing mapping
// functions under readmodels/ and writemodels/ remain in place temporarily for
// backward compatibility until the consolidation subphases migrate call sites.
package mappers

import (
	authdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/domain"
	authreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/readmodels"
)

// NewAuthLoginReadModel builds the explicit login output projection from
// resolved authentication values. It keeps response construction outside HTTP
// handlers and application orchestration.
func NewAuthLoginReadModel(accessToken, tokenType string, expiresIn int64, userID string) authreadmodels.AuthLoginReadModel {
	return authreadmodels.AuthLoginReadModel{
		AccessToken: accessToken,
		TokenType:   tokenType,
		ExpiresIn:   expiresIn,
		UserID:      userID,
	}
}

// WalletIdentityToReadModel maps a canonical wallet identity domain model into
// its output-only auth wallet read projection.
func WalletIdentityToReadModel(wallet *authdomain.WalletIdentity) *authreadmodels.AuthWalletReadModel {
	if wallet == nil {
		return nil
	}

	return &authreadmodels.AuthWalletReadModel{
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

// WalletIdentitiesToReadModels maps wallet identity domain models into
// output-only read projections. Nil inputs return an empty slice to preserve
// response stability.
func WalletIdentitiesToReadModels(wallets []*authdomain.WalletIdentity) []*authreadmodels.AuthWalletReadModel {
	if wallets == nil {
		return []*authreadmodels.AuthWalletReadModel{}
	}

	out := make([]*authreadmodels.AuthWalletReadModel, 0, len(wallets))
	for _, wallet := range wallets {
		mapped := WalletIdentityToReadModel(wallet)
		if mapped != nil {
			out = append(out, mapped)
		}
	}
	return out
}

// WalletIdentitiesToActionableReadModels maps wallet identity domain models into
// output-only wallet read projections and enriches them with response
// actionability metadata. This keeps response-oriented mapping concerns inside
// the module-level mapping layer instead of application orchestration.
func WalletIdentitiesToActionableReadModels(wallets []*authdomain.WalletIdentity) []*authreadmodels.AuthWalletReadModel {
	return EnrichWalletReadModelsActionability(WalletIdentitiesToReadModels(wallets))
}

// EnrichWalletReadModelsActionability applies response actionability metadata to
// wallet read models. It preserves the existing contract while centralizing the
// calculation outside handlers and application flows.
func EnrichWalletReadModelsActionability(wallets []*authreadmodels.AuthWalletReadModel) []*authreadmodels.AuthWalletReadModel {
	if len(wallets) == 0 {
		return []*authreadmodels.AuthWalletReadModel{}
	}

	activeOwnedCount := 0
	for _, wallet := range wallets {
		if wallet != nil && wallet.Status == "active" {
			activeOwnedCount++
		}
	}

	for _, wallet := range wallets {
		if wallet == nil {
			continue
		}

		wallet.CanSetPrimary = wallet.Status == "active" && !wallet.IsPrimary
		wallet.CanDetach = false
		wallet.DetachBlockReasons = []string{}

		if wallet.Status != "active" {
			continue
		}

		if wallet.IsPrimary {
			wallet.DetachBlockReasons = append(wallet.DetachBlockReasons, authdomain.WalletDetachReasonWalletIsPrimary)
		}
		if activeOwnedCount <= 1 {
			wallet.DetachBlockReasons = append(wallet.DetachBlockReasons, authdomain.WalletDetachReasonUserWouldBeEmpty)
		}
		if len(wallet.DetachBlockReasons) == 0 {
			wallet.CanDetach = true
		}
	}

	return wallets
}

// WalletChallengeToReadModel maps a canonical wallet challenge domain model into
// its output-only auth wallet challenge projection.
func WalletChallengeToReadModel(challenge *authdomain.WalletChallenge) *authreadmodels.AuthWalletChallengeReadModel {
	if challenge == nil {
		return nil
	}

	return &authreadmodels.AuthWalletChallengeReadModel{
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
