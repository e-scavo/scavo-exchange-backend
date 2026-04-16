package app

import (
	"time"

	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"

	authdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/domain"
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	UserID      string `json:"user_id"`
}

type SessionView struct {
	Authenticated bool          `json:"authenticated"`
	TokenType     string        `json:"token_type"`
	UserID        string        `json:"user_id"`
	Email         string        `json:"email,omitempty"`
	WalletID      string        `json:"wallet_id,omitempty"`
	WalletAddress string        `json:"wallet_address,omitempty"`
	AuthMethod    string        `json:"auth_method,omitempty"`
	Chain         string        `json:"chain,omitempty"`
	Subject       string        `json:"subject,omitempty"`
	Issuer        string        `json:"issuer,omitempty"`
	ExpiresAt     *time.Time    `json:"expires_at,omitempty"`
	User          *usermod.User `json:"user,omitempty"`
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
	User                *usermod.User        `json:"user,omitempty"`
	UserID              string               `json:"user_id"`
	AuthMethod          string               `json:"auth_method,omitempty"`
	WalletID            string               `json:"wallet_id,omitempty"`
	WalletAddress       string               `json:"wallet_address,omitempty"`
	Chain               string               `json:"chain,omitempty"`
	PrimaryWallet       *ProfileWalletView   `json:"primary_wallet,omitempty"`
	Wallets             []*ProfileWalletView `json:"wallets"`
	WalletCount         int                  `json:"wallet_count"`
	ActiveWalletCount   int                  `json:"active_wallet_count"`
	DetachedWalletCount int                  `json:"detached_wallet_count"`
	HasWalletSession    bool                 `json:"has_wallet_session"`
}

type MeResponse struct {
	User    *usermod.User `json:"user"`
	Profile *ProfileView  `json:"profile,omitempty"`
}

type SessionResponse struct {
	Session *SessionView `json:"session"`
}

type BootstrapWalletsView struct {
	Items []*WalletReadModel `json:"items"`
	Total int                `json:"total"`
}

type BootstrapResponse struct {
	Session  *SessionView         `json:"session"`
	User     *usermod.User        `json:"user,omitempty"`
	Profile  *ProfileView         `json:"profile,omitempty"`
	Settings usersettingsmod.View `json:"settings"`
	Wallets  BootstrapWalletsView `json:"wallets"`
}

type WalletReadModel struct {
	ID                 string     `json:"id"`
	Address            string     `json:"address"`
	UserID             string     `json:"user_id,omitempty"`
	LinkedAt           *time.Time `json:"linked_at,omitempty"`
	DetachedAt         *time.Time `json:"detached_at,omitempty"`
	IsPrimary          bool       `json:"is_primary"`
	Status             string     `json:"status"`
	CanSetPrimary      bool       `json:"can_set_primary"`
	CanDetach          bool       `json:"can_detach"`
	DetachBlockReasons []string   `json:"detach_block_reasons"`
}

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
	Challenge *authdomain.WalletChallenge `json:"challenge"`
}

type WalletLinkVerifyResponse struct {
	LinkedWallet *authdomain.WalletIdentity   `json:"linked_wallet,omitempty"`
	Wallets      []*authdomain.WalletIdentity `json:"wallets"`
	Challenge    *authdomain.WalletChallenge  `json:"challenge,omitempty"`
}

type WalletAccountMergeChallengeResponse struct {
	Challenge *authdomain.WalletChallenge `json:"challenge"`
}

type WalletAccountMergeVerifyResponse struct {
	MergedWallet *authdomain.WalletIdentity   `json:"merged_wallet,omitempty"`
	Wallets      []*authdomain.WalletIdentity `json:"wallets"`
	Challenge    *authdomain.WalletChallenge  `json:"challenge,omitempty"`
	SourceUserID string                       `json:"source_user_id"`
	TargetUserID string                       `json:"target_user_id"`
}

type WalletDetachCheckResponse struct {
	WalletAddress    string   `json:"wallet_address"`
	Eligible         bool     `json:"eligible"`
	IsPrimary        bool     `json:"is_primary"`
	OwnedWalletCount int      `json:"owned_wallet_count"`
	Reasons          []string `json:"reasons"`
}

type WalletDetachExecuteResponse struct {
	DetachedWallet *authdomain.WalletIdentity   `json:"detached_wallet,omitempty"`
	Wallets        []*authdomain.WalletIdentity `json:"wallets"`
	Check          *WalletDetachCheckResponse   `json:"check,omitempty"`
}

type WalletPrimarySetResponse struct {
	PrimaryWallet *authdomain.WalletIdentity   `json:"primary_wallet,omitempty"`
	Wallets       []*authdomain.WalletIdentity `json:"wallets"`
}
