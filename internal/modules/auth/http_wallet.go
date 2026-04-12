package auth

import (
	"errors"
	"net/http"
	"time"

	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

type WalletChallengeRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
}

type WalletChallengeResponse struct {
	Challenge *WalletChallenge `json:"challenge"`
}

type WalletVerifyRequest struct {
	ChallengeID string `json:"challenge_id"`
	Address     string `json:"address"`
	Signature   string `json:"signature"`
}

type WalletVerifyResponse struct {
	AccessToken   string           `json:"access_token"`
	TokenType     string           `json:"token_type"`
	ExpiresIn     int64            `json:"expires_in"`
	UserID        string           `json:"user_id"`
	WalletID      string           `json:"wallet_id,omitempty"`
	WalletAddress string           `json:"wallet_address"`
	Chain         string           `json:"chain"`
	AuthMethod    string           `json:"auth_method"`
	User          *usermod.User    `json:"user,omitempty"`
	Challenge     *WalletChallenge `json:"challenge,omitempty"`
}

type WalletLinkChallengeRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
}

type WalletLinkChallengeResponse struct {
	Challenge *WalletChallenge `json:"challenge"`
}

type WalletLinkVerifyRequest struct {
	ChallengeID string `json:"challenge_id"`
	Address     string `json:"address"`
	Signature   string `json:"signature"`
}

type WalletLinkVerifyResponse struct {
	LinkedWallet *WalletIdentity   `json:"linked_wallet,omitempty"`
	Wallets      []*WalletIdentity `json:"wallets"`
	Challenge    *WalletChallenge  `json:"challenge,omitempty"`
}

type WalletAccountMergeChallengeRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
}

type WalletAccountMergeChallengeResponse struct {
	Challenge *WalletChallenge `json:"challenge"`
}

type WalletAccountMergeVerifyRequest struct {
	ChallengeID string `json:"challenge_id"`
	Address     string `json:"address"`
	Signature   string `json:"signature"`
}

type WalletAccountMergeVerifyResponse struct {
	MergedWallet *WalletIdentity   `json:"merged_wallet,omitempty"`
	Wallets      []*WalletIdentity `json:"wallets"`
	Challenge    *WalletChallenge  `json:"challenge,omitempty"`
	SourceUserID string            `json:"source_user_id"`
	TargetUserID string            `json:"target_user_id"`
}

type WalletDetachCheckRequest struct {
	Address string `json:"wallet_address"`
}

type WalletDetachCheckResponse struct {
	WalletAddress    string   `json:"wallet_address"`
	Eligible         bool     `json:"eligible"`
	IsPrimary        bool     `json:"is_primary"`
	OwnedWalletCount int      `json:"owned_wallet_count"`
	Reasons          []string `json:"reasons"`
}

type WalletDetachExecuteRequest struct {
	Address string `json:"wallet_address"`
}

type WalletDetachExecuteResponse struct {
	DetachedWallet *WalletIdentity            `json:"detached_wallet,omitempty"`
	Wallets        []*WalletIdentity          `json:"wallets"`
	Check          *WalletDetachCheckResponse `json:"check,omitempty"`
}

type WalletPrimarySetRequest struct {
	Address string `json:"wallet_address"`
}

type WalletPrimarySetResponse struct {
	PrimaryWallet *WalletIdentity   `json:"primary_wallet,omitempty"`
	Wallets       []*WalletIdentity `json:"wallets"`
}

func (h HTTPHandlers) WalletChallenge(w http.ResponseWriter, r *http.Request) {
	var req WalletChallengeRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	challengeTTL := h.ChallengeTTL
	if challengeTTL <= 0 {
		challengeTTL = 5 * time.Minute
	}

	svc := NewWalletChallengeService(h.Challenges, h.PublicBaseURL, challengeTTL)
	challenge, err := svc.Create(r.Context(), req.Address, req.Chain)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidWalletAddress):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_wallet_address")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "wallet_challenge_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, WalletChallengeResponse{Challenge: challenge})
}

func (h HTTPHandlers) WalletVerify(w http.ResponseWriter, r *http.Request) {
	var req WalletVerifyRequest
	if !decodeRequest(w, r, &req, 8<<10) {
		return
	}

	challengeTTL := h.ChallengeTTL
	if challengeTTL <= 0 {
		challengeTTL = 5 * time.Minute
	}

	challengeSvc := NewWalletChallengeService(h.Challenges, h.PublicBaseURL, challengeTTL)
	loginSvc := NewService(h.Tokens, h.Users, h.TTL)
	verifySvc := NewWalletVerificationService(challengeSvc, loginSvc, h.WalletIdentities)

	result, challenge, err := verifySvc.VerifyAndLogin(r.Context(), req.ChallengeID, req.Address, req.Signature)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidWalletAddress):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_wallet_address")
		case errors.Is(err, ErrInvalidWalletSignature):
			writeErrorJSON(w, http.StatusUnauthorized, "invalid_wallet_signature")
		case errors.Is(err, ErrWalletChallengeNotFound):
			writeErrorJSON(w, http.StatusNotFound, "wallet_challenge_not_found")
		case errors.Is(err, ErrChallengeExpired):
			writeErrorJSON(w, http.StatusUnauthorized, "wallet_challenge_expired")
		case errors.Is(err, ErrChallengeUsed):
			writeErrorJSON(w, http.StatusUnauthorized, "wallet_challenge_used")
		case errors.Is(err, ErrWalletChallengePurpose):
			writeErrorJSON(w, http.StatusConflict, "wallet_challenge_purpose_mismatch")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "wallet_verify_error")
		}
		return
	}

	userID := ""
	var user *usermod.User
	if result != nil && result.User != nil {
		userID = result.User.ID
		user = result.User
	}

	writeJSON(w, http.StatusOK, WalletVerifyResponse{
		AccessToken:   result.AccessToken,
		TokenType:     result.TokenType,
		ExpiresIn:     result.ExpiresIn,
		UserID:        userID,
		WalletID:      result.WalletID,
		WalletAddress: result.WalletAddress,
		Chain:         result.Chain,
		AuthMethod:    result.AuthMethod,
		User:          user,
		Challenge:     challenge,
	})
}

func (h HTTPHandlers) WalletLinkChallenge(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	var req WalletLinkChallengeRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	response, err := h.Application().CreateWalletLinkChallenge(r.Context(), claims.UserID, req.Address, req.Chain)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, ErrInvalidWalletAddress):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_wallet_address")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "wallet_link_challenge_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletLinkVerify(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	var req WalletLinkVerifyRequest
	if !decodeRequest(w, r, &req, 8<<10) {
		return
	}

	response, err := h.Application().VerifyWalletLink(r.Context(), claims.UserID, req.ChallengeID, req.Address, req.Signature)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, ErrInvalidWalletAddress):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_wallet_address")
		case errors.Is(err, ErrInvalidWalletSignature):
			writeErrorJSON(w, http.StatusUnauthorized, "invalid_wallet_signature")
		case errors.Is(err, ErrWalletChallengeNotFound):
			writeErrorJSON(w, http.StatusNotFound, "wallet_challenge_not_found")
		case errors.Is(err, ErrChallengeExpired):
			writeErrorJSON(w, http.StatusUnauthorized, "wallet_challenge_expired")
		case errors.Is(err, ErrChallengeUsed):
			writeErrorJSON(w, http.StatusUnauthorized, "wallet_challenge_used")
		case errors.Is(err, ErrWalletChallengePurpose):
			writeErrorJSON(w, http.StatusConflict, "wallet_challenge_purpose_mismatch")
		case errors.Is(err, ErrWalletLinkChallengeMismatch):
			writeErrorJSON(w, http.StatusForbidden, "wallet_link_challenge_user_mismatch")
		case errors.Is(err, ErrWalletIdentityAlreadyLinked):
			writeErrorJSON(w, http.StatusConflict, "wallet_identity_already_linked")
		case errors.Is(err, ErrWalletAlreadyLinkedToUser):
			writeErrorJSON(w, http.StatusConflict, "wallet_identity_already_linked_to_user")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "wallet_link_verify_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletAccountMergeChallenge(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	var req WalletAccountMergeChallengeRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	response, err := h.Application().CreateWalletAccountMergeChallenge(r.Context(), claims.UserID, req.Address, req.Chain)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, ErrInvalidWalletAddress):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_wallet_address")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "wallet_account_merge_challenge_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletAccountMergeVerify(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	var req WalletAccountMergeVerifyRequest
	if !decodeRequest(w, r, &req, 8<<10) {
		return
	}

	response, err := h.Application().VerifyWalletAccountMerge(r.Context(), claims.UserID, req.ChallengeID, req.Address, req.Signature)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, ErrInvalidWalletAddress):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_wallet_address")
		case errors.Is(err, ErrInvalidWalletSignature):
			writeErrorJSON(w, http.StatusUnauthorized, "invalid_wallet_signature")
		case errors.Is(err, ErrWalletChallengeNotFound):
			writeErrorJSON(w, http.StatusNotFound, "wallet_challenge_not_found")
		case errors.Is(err, ErrChallengeExpired):
			writeErrorJSON(w, http.StatusUnauthorized, "wallet_challenge_expired")
		case errors.Is(err, ErrChallengeUsed):
			writeErrorJSON(w, http.StatusUnauthorized, "wallet_challenge_used")
		case errors.Is(err, ErrWalletChallengePurpose):
			writeErrorJSON(w, http.StatusConflict, "wallet_challenge_purpose_mismatch")
		case errors.Is(err, ErrWalletLinkChallengeMismatch):
			writeErrorJSON(w, http.StatusForbidden, "wallet_account_merge_user_mismatch")
		case errors.Is(err, ErrWalletMergeSourceNotLinked):
			writeErrorJSON(w, http.StatusConflict, "wallet_account_merge_source_not_linked")
		case errors.Is(err, ErrWalletMergeSameUser):
			writeErrorJSON(w, http.StatusConflict, "wallet_account_merge_not_required")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "wallet_account_merge_verify_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletDetachCheck(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	var req WalletDetachCheckRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	response, err := h.Application().CheckWalletDetach(r.Context(), claims.UserID, req.Address)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, ErrInvalidWalletAddress):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_wallet_address")
		case errors.Is(err, ErrWalletIdentityNotFound):
			writeErrorJSON(w, http.StatusNotFound, "wallet_identity_not_found")
		case errors.Is(err, ErrWalletNotOwnedByUser):
			writeErrorJSON(w, http.StatusForbidden, "wallet_identity_not_owned_by_user")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "wallet_detach_check_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletDetach(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	var req WalletDetachExecuteRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	response, err := h.Application().ExecuteWalletDetach(r.Context(), claims.UserID, req.Address)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, ErrInvalidWalletAddress):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_wallet_address")
		case errors.Is(err, ErrWalletIdentityNotFound):
			writeErrorJSON(w, http.StatusNotFound, "wallet_identity_not_found")
		case errors.Is(err, ErrWalletNotOwnedByUser):
			writeErrorJSON(w, http.StatusForbidden, "wallet_identity_not_owned_by_user")
		case errors.Is(err, ErrWalletDetachNotEligible):
			writeErrorJSON(w, http.StatusConflict, "wallet_detach_not_eligible", map[string]any{"check": response.Check})
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "wallet_detach_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletSetPrimary(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	var req WalletPrimarySetRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	response, err := h.Application().SetPrimaryWallet(r.Context(), claims.UserID, req.Address)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, ErrInvalidWalletAddress):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_wallet_address")
		case errors.Is(err, ErrWalletIdentityNotFound):
			writeErrorJSON(w, http.StatusNotFound, "wallet_identity_not_found")
		case errors.Is(err, ErrWalletNotOwnedByUser):
			writeErrorJSON(w, http.StatusForbidden, "wallet_identity_not_owned_by_user")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "wallet_primary_set_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}
