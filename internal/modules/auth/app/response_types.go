package app

import (
	"time"

	authreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/readmodels"
	userreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/readmodels"
	usersettingsreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/readmodels"
)

type LoginResponse = authreadmodels.AuthLoginReadModel

type SessionView struct {
	Authenticated bool                          `json:"authenticated"`
	TokenType     string                        `json:"token_type"`
	UserID        string                        `json:"user_id"`
	Email         string                        `json:"email,omitempty"`
	WalletID      string                        `json:"wallet_id,omitempty"`
	WalletAddress string                        `json:"wallet_address,omitempty"`
	AuthMethod    string                        `json:"auth_method,omitempty"`
	Chain         string                        `json:"chain,omitempty"`
	Subject       string                        `json:"subject,omitempty"`
	Issuer        string                        `json:"issuer,omitempty"`
	ExpiresAt     *time.Time                    `json:"expires_at,omitempty"`
	User          *userreadmodels.UserReadModel `json:"user,omitempty"`
}

type ProfileWalletView struct {
	ID         string     `json:"id"`
	Address    string     `json:"address"`
	IsPrimary  bool       `json:"is_primary"`
	Status     string     `json:"status"`
	LinkedAt   *time.Time `json:"linked_at,omitempty"`
	DetachedAt *time.Time `json:"detached_at,omitempty"`
}

type ProfileView struct {
	User                *userreadmodels.UserReadModel `json:"user,omitempty"`
	UserID              string                        `json:"user_id"`
	AuthMethod          string                        `json:"auth_method,omitempty"`
	WalletID            string                        `json:"wallet_id,omitempty"`
	WalletAddress       string                        `json:"wallet_address,omitempty"`
	Chain               string                        `json:"chain,omitempty"`
	PrimaryWallet       *ProfileWalletView            `json:"primary_wallet,omitempty"`
	Wallets             []*ProfileWalletView          `json:"wallets"`
	WalletCount         int                           `json:"wallet_count"`
	ActiveWalletCount   int                           `json:"active_wallet_count"`
	DetachedWalletCount int                           `json:"detached_wallet_count"`
	HasWalletSession    bool                          `json:"has_wallet_session"`
}

type MeResponse struct {
	User    *userreadmodels.UserReadModel `json:"user"`
	Profile *ProfileView                  `json:"profile,omitempty"`
}

type SessionResponse struct {
	Session *SessionView `json:"session"`
}

type MeSettingsResponse struct {
	Settings usersettingsreadmodels.UserSettingsReadModel `json:"settings"`
}

type WalletChallengeResponse struct {
	Challenge *authreadmodels.AuthWalletChallengeReadModel `json:"challenge"`
}

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

type BootstrapWalletsView struct {
	Items []*WalletReadModel `json:"items"`
	Total int                `json:"total"`
}

type BootstrapResponse struct {
	Session  *SessionView                                 `json:"session"`
	User     *userreadmodels.UserReadModel                `json:"user,omitempty"`
	Profile  *ProfileView                                 `json:"profile,omitempty"`
	Settings usersettingsreadmodels.UserSettingsReadModel `json:"settings"`
	Wallets  BootstrapWalletsView                         `json:"wallets"`
}

type WalletReadModel = authreadmodels.AuthWalletReadModel

type WalletsQuery struct {
	Status         string
	Primary        *bool
	Sort           string
	Order          string
	Limit          int
	Offset         int
	SortProvided   bool
	OrderProvided  bool
	LimitProvided  bool
	OffsetProvided bool
}

type WalletsResponse struct {
	Items          []*WalletReadModel `json:"items"`
	Wallets        []*WalletReadModel `json:"wallets"`
	Total          int                `json:"total"`
	Limit          int                `json:"limit"`
	Offset         int                `json:"offset"`
	Returned       int                `json:"returned"`
	HasMore        bool               `json:"has_more"`
	NextOffset     *int               `json:"next_offset,omitempty"`
	PreviousOffset *int               `json:"previous_offset,omitempty"`
}

type WalletLinkChallengeResponse struct {
	Challenge *authreadmodels.AuthWalletChallengeReadModel `json:"challenge"`
}

type WalletLinkVerifyResponse struct {
	LinkedWallet *WalletReadModel                             `json:"linked_wallet,omitempty"`
	Wallets      []*WalletReadModel                           `json:"wallets"`
	Challenge    *authreadmodels.AuthWalletChallengeReadModel `json:"challenge,omitempty"`
}

type WalletAccountMergeChallengeResponse struct {
	Challenge *authreadmodels.AuthWalletChallengeReadModel `json:"challenge"`
}

type WalletAccountMergeVerifyResponse struct {
	MergedWallet *WalletReadModel                             `json:"merged_wallet,omitempty"`
	Wallets      []*WalletReadModel                           `json:"wallets"`
	Challenge    *authreadmodels.AuthWalletChallengeReadModel `json:"challenge,omitempty"`
	SourceUserID string                                       `json:"source_user_id"`
	TargetUserID string                                       `json:"target_user_id"`
}

type WalletDetachCheckResponse struct {
	WalletAddress    string   `json:"wallet_address"`
	Eligible         bool     `json:"eligible"`
	IsPrimary        bool     `json:"is_primary"`
	OwnedWalletCount int      `json:"owned_wallet_count"`
	Reasons          []string `json:"reasons"`
}

type WalletDetachExecuteResponse struct {
	DetachedWallet *WalletReadModel           `json:"detached_wallet,omitempty"`
	Wallets        []*WalletReadModel         `json:"wallets"`
	Check          *WalletDetachCheckResponse `json:"check,omitempty"`
}

type WalletPrimarySetResponse struct {
	PrimaryWallet *WalletReadModel   `json:"primary_wallet,omitempty"`
	Wallets       []*WalletReadModel `json:"wallets"`
}
