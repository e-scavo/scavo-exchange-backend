package auth

import (
	authapp "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/app"
	authreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/readmodels"
	authwritemodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/writemodels"
	userreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/readmodels"
	usersettingsreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/readmodels"
)

// LoginRequest preserves the public login HTTP contract while pointing the
// handler-facing input to the explicit auth write model.
type LoginRequest = authwritemodels.AuthLoginWriteModel

// LoginResponse preserves the public login HTTP response contract.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	UserID      string `json:"user_id"`
}

// UpdateMeRequest preserves the public profile-update HTTP contract while
// pointing the handler-facing input to the explicit auth write model.
type UpdateMeRequest = authwritemodels.AuthUpdateProfileWriteModel

// MeResponse preserves the public /me HTTP response contract while exposing
// read-model output for the user surface.
type MeResponse struct {
	User    *userreadmodels.UserReadModel `json:"user"`
	Profile *ProfileView                  `json:"profile,omitempty"`
}

// SessionResponse preserves the public session HTTP response contract.
type SessionResponse struct {
	Session *SessionView `json:"session"`
}

// UpdateMeSettingsRequest preserves the public settings-update HTTP contract
// while pointing the handler-facing input to the explicit auth write model.
type UpdateMeSettingsRequest = authwritemodels.AuthUpdateSettingsWriteModel

// MeSettingsResponse preserves the public settings HTTP response contract while
// exposing read-model output for the settings surface.
type MeSettingsResponse struct {
	Settings usersettingsreadmodels.UserSettingsReadModel `json:"settings"`
}

// WalletChallengeRequest preserves the public wallet-challenge HTTP contract
// while pointing the handler-facing input to the explicit auth write model.
type WalletChallengeRequest = authwritemodels.AuthWalletChallengeWriteModel

// WalletChallengeResponse preserves the public wallet-challenge HTTP response
// contract while exposing the explicit auth challenge read model.
type WalletChallengeResponse struct {
	Challenge *authreadmodels.AuthWalletChallengeReadModel `json:"challenge"`
}

// WalletVerifyRequest preserves the public wallet-verification HTTP contract
// while pointing the handler-facing input to the explicit auth write model.
type WalletVerifyRequest = authwritemodels.AuthWalletVerifyWriteModel

// WalletVerifyResponse preserves the public wallet-verification HTTP response
// contract while exposing read-model output where applicable.
type WalletVerifyResponse struct {
	AccessToken   string                                       `json:"access_token"`
	TokenType     string                                       `json:"token_type"`
	ExpiresIn     int64                                        `json:"expires_in"`
	UserID        string                                       `json:"user_id"`
	WalletID      string                                       `json:"wallet_id,omitempty"`
	WalletAddress string                                       `json:"wallet_address"`
	Chain         string                                       `json:"chain"`
	AuthMethod    string                                       `json:"auth_method"`
	User          *userreadmodels.UserReadModel                `json:"user,omitempty"`
	Challenge     *authreadmodels.AuthWalletChallengeReadModel `json:"challenge,omitempty"`
}

// Authenticated wallet-management request contracts preserve public HTTP names
// while pointing handler-facing inputs to explicit auth write models.
type WalletLinkChallengeRequest = authwritemodels.AuthWalletLinkChallengeWriteModel
type WalletLinkVerifyRequest = authwritemodels.AuthWalletLinkVerifyWriteModel
type WalletAccountMergeChallengeRequest = authwritemodels.AuthWalletAccountMergeChallengeWriteModel
type WalletAccountMergeVerifyRequest = authwritemodels.AuthWalletAccountMergeVerifyWriteModel
type WalletDetachCheckRequest = authwritemodels.AuthWalletDetachCheckWriteModel
type WalletDetachExecuteRequest = authwritemodels.AuthWalletDetachExecuteWriteModel
type WalletPrimarySetRequest = authwritemodels.AuthWalletPrimarySetWriteModel

// Authenticated wallet-management response contracts preserve the stable
// application response shapes already aligned during Phase 0.12.2 and Phase
// 0.12.4.
type WalletLinkChallengeResponse = authapp.WalletLinkChallengeResponse
type WalletLinkVerifyResponse = authapp.WalletLinkVerifyResponse
type WalletAccountMergeChallengeResponse = authapp.WalletAccountMergeChallengeResponse
type WalletAccountMergeVerifyResponse = authapp.WalletAccountMergeVerifyResponse
type WalletDetachCheckResponse = authapp.WalletDetachCheckResponse
type WalletDetachExecuteResponse = authapp.WalletDetachExecuteResponse
type WalletPrimarySetResponse = authapp.WalletPrimarySetResponse
