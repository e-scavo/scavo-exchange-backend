// Package writemodels contains explicit input-only models for the auth module.
//
// Write models represent caller intent and must not be reused as responses or
// domain state. They are introduced additively in Phase 0.12.3 so existing HTTP
// contracts and handlers can migrate progressively without runtime breakage.
package writemodels

import "encoding/json"

// AuthLoginWriteModel represents the login input intent.
type AuthLoginWriteModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthUpdateProfileWriteModel represents the authenticated profile update input
// intent currently exposed through the /me update surface.
type AuthUpdateProfileWriteModel struct {
	DisplayName string `json:"display_name"`
}

// AuthUpdateSettingsWriteModel represents authenticated user settings update
// input intent. Raw JSON is preserved to keep existing request compatibility
// until handler alignment introduces explicit conversion.
type AuthUpdateSettingsWriteModel struct {
	Preferences json.RawMessage `json:"preferences"`
}

// AuthWalletChallengeWriteModel represents public wallet challenge creation
// input intent.
type AuthWalletChallengeWriteModel struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
}

// AuthWalletVerifyWriteModel represents public wallet verification input intent.
type AuthWalletVerifyWriteModel struct {
	ChallengeID string `json:"challenge_id"`
	Address     string `json:"address"`
	Signature   string `json:"signature"`
}

// AuthWalletLinkChallengeWriteModel represents authenticated wallet-link
// challenge creation input intent.
type AuthWalletLinkChallengeWriteModel struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
}

// AuthWalletLinkVerifyWriteModel represents authenticated wallet-link
// verification input intent.
type AuthWalletLinkVerifyWriteModel struct {
	ChallengeID string `json:"challenge_id"`
	Address     string `json:"address"`
	Signature   string `json:"signature"`
}

// AuthWalletAccountMergeChallengeWriteModel represents authenticated account
// merge challenge creation input intent.
type AuthWalletAccountMergeChallengeWriteModel struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
}

// AuthWalletAccountMergeVerifyWriteModel represents authenticated account merge
// verification input intent.
type AuthWalletAccountMergeVerifyWriteModel struct {
	ChallengeID string `json:"challenge_id"`
	Address     string `json:"address"`
	Signature   string `json:"signature"`
}

// AuthWalletDetachCheckWriteModel represents authenticated wallet detach safety
// check input intent.
type AuthWalletDetachCheckWriteModel struct {
	Address string `json:"wallet_address"`
}

// AuthWalletDetachExecuteWriteModel represents authenticated wallet detach
// execution input intent.
type AuthWalletDetachExecuteWriteModel struct {
	Address string `json:"wallet_address"`
}

// AuthWalletPrimarySetWriteModel represents authenticated primary wallet update
// input intent.
type AuthWalletPrimarySetWriteModel struct {
	Address string `json:"wallet_address"`
}
