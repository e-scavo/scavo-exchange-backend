package domain

// LoginInput represents canonical authentication input after transport-level
// write models have been decoded and mapped.
type LoginInput struct {
	Email    string
	Password string
}

// ProfileUpdateInput represents canonical authenticated profile update intent.
type ProfileUpdateInput struct {
	DisplayName string
}

// SettingsUpdateInput represents canonical authenticated settings update
// intent. Preferences are already decoded from the transport payload.
type SettingsUpdateInput struct {
	Preferences map[string]any
}

// WalletChallengeInput represents canonical wallet challenge creation intent.
type WalletChallengeInput struct {
	Address string
	Chain   string
}

// WalletVerifyInput represents canonical wallet verification intent.
type WalletVerifyInput struct {
	ChallengeID string
	Address     string
	Signature   string
}

// WalletDetachInput represents canonical wallet detach/check intent.
type WalletDetachInput struct {
	Address string
}

// WalletPrimarySetInput represents canonical primary wallet update intent.
type WalletPrimarySetInput struct {
	Address string
}
