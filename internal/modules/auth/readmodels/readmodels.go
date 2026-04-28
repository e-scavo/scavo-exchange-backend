package readmodels

import "time"

// AuthLoginReadModel is the explicit output projection for successful
// authentication responses. It is intentionally output-only and must not be
// reused as an input or command payload.
type AuthLoginReadModel struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	UserID      string `json:"user_id"`
}

// AuthWalletReadModel is the explicit output projection for wallet identity
// data exposed by auth-facing read surfaces.
type AuthWalletReadModel struct {
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

// AuthWalletChallengeReadModel is the explicit output projection for wallet
// challenge responses. It does not replace the internal challenge domain state.
type AuthWalletChallengeReadModel struct {
	ID                string     `json:"id"`
	Address           string     `json:"address"`
	Chain             string     `json:"chain"`
	Nonce             string     `json:"nonce"`
	Message           string     `json:"message"`
	Purpose           string     `json:"purpose"`
	RequestedByUserID string     `json:"requested_by_user_id,omitempty"`
	IssuedAt          time.Time  `json:"issued_at"`
	ExpiresAt         time.Time  `json:"expires_at"`
	UsedAt            *time.Time `json:"used_at,omitempty"`
}
