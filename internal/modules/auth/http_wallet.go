package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
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
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_challenge_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, WalletChallengeResponse{Challenge: challenge})
}

func (h HTTPHandlers) WalletVerify(w http.ResponseWriter, r *http.Request) {
	var req WalletVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
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
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		case errors.Is(err, ErrInvalidWalletSignature):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid_wallet_signature"})
		case errors.Is(err, ErrWalletChallengeNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "wallet_challenge_not_found"})
		case errors.Is(err, ErrChallengeExpired):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "wallet_challenge_expired"})
		case errors.Is(err, ErrChallengeUsed):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "wallet_challenge_used"})
		case errors.Is(err, ErrWalletChallengePurpose):
			writeJSON(w, http.StatusConflict, map[string]any{"error": "wallet_challenge_purpose_mismatch"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_verify_error"})
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
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok || claims == nil || claims.UserID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	var req WalletLinkChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	challengeTTL := h.ChallengeTTL
	if challengeTTL <= 0 {
		challengeTTL = 5 * time.Minute
	}

	challengeSvc := NewWalletChallengeService(h.Challenges, h.PublicBaseURL, challengeTTL)
	linkSvc := NewWalletLinkingService(challengeSvc, h.WalletIdentities)

	challenge, err := linkSvc.CreateChallenge(r.Context(), claims.UserID, req.Address, req.Chain)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case errors.Is(err, ErrInvalidWalletAddress):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_link_challenge_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, WalletLinkChallengeResponse{Challenge: challenge})
}

func (h HTTPHandlers) WalletLinkVerify(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok || claims == nil || claims.UserID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	var req WalletLinkVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	challengeTTL := h.ChallengeTTL
	if challengeTTL <= 0 {
		challengeTTL = 5 * time.Minute
	}

	challengeSvc := NewWalletChallengeService(h.Challenges, h.PublicBaseURL, challengeTTL)
	linkSvc := NewWalletLinkingService(challengeSvc, h.WalletIdentities)

	result, err := linkSvc.VerifyAndLink(r.Context(), claims.UserID, req.ChallengeID, req.Address, req.Signature)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case errors.Is(err, ErrInvalidWalletAddress):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		case errors.Is(err, ErrInvalidWalletSignature):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid_wallet_signature"})
		case errors.Is(err, ErrWalletChallengeNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "wallet_challenge_not_found"})
		case errors.Is(err, ErrChallengeExpired):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "wallet_challenge_expired"})
		case errors.Is(err, ErrChallengeUsed):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "wallet_challenge_used"})
		case errors.Is(err, ErrWalletChallengePurpose):
			writeJSON(w, http.StatusConflict, map[string]any{"error": "wallet_challenge_purpose_mismatch"})
		case errors.Is(err, ErrWalletLinkChallengeMismatch):
			writeJSON(w, http.StatusForbidden, map[string]any{"error": "wallet_link_challenge_user_mismatch"})
		case errors.Is(err, ErrWalletIdentityAlreadyLinked):
			writeJSON(w, http.StatusConflict, map[string]any{"error": "wallet_identity_already_linked"})
		case errors.Is(err, ErrWalletAlreadyLinkedToUser):
			writeJSON(w, http.StatusConflict, map[string]any{"error": "wallet_identity_already_linked_to_user"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_link_verify_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, WalletLinkVerifyResponse{
		LinkedWallet: result.Linked,
		Wallets:      result.Wallets,
		Challenge:    result.Challenge,
	})
}

func (h HTTPHandlers) WalletAccountMergeChallenge(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok || claims == nil || claims.UserID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	var req WalletAccountMergeChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	challengeTTL := h.ChallengeTTL
	if challengeTTL <= 0 {
		challengeTTL = 5 * time.Minute
	}

	challengeSvc := NewWalletChallengeService(h.Challenges, h.PublicBaseURL, challengeTTL)
	mergeSvc := NewWalletAccountMergeService(challengeSvc, h.WalletIdentities)

	challenge, err := mergeSvc.CreateChallenge(r.Context(), claims.UserID, req.Address, req.Chain)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case errors.Is(err, ErrInvalidWalletAddress):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_account_merge_challenge_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, WalletAccountMergeChallengeResponse{Challenge: challenge})
}

func (h HTTPHandlers) WalletAccountMergeVerify(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok || claims == nil || claims.UserID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	var req WalletAccountMergeVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	challengeTTL := h.ChallengeTTL
	if challengeTTL <= 0 {
		challengeTTL = 5 * time.Minute
	}

	challengeSvc := NewWalletChallengeService(h.Challenges, h.PublicBaseURL, challengeTTL)
	mergeSvc := NewWalletAccountMergeService(challengeSvc, h.WalletIdentities)

	result, err := mergeSvc.VerifyAndMerge(r.Context(), claims.UserID, req.ChallengeID, req.Address, req.Signature)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case errors.Is(err, ErrInvalidWalletAddress):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		case errors.Is(err, ErrInvalidWalletSignature):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid_wallet_signature"})
		case errors.Is(err, ErrWalletChallengeNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "wallet_challenge_not_found"})
		case errors.Is(err, ErrChallengeExpired):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "wallet_challenge_expired"})
		case errors.Is(err, ErrChallengeUsed):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "wallet_challenge_used"})
		case errors.Is(err, ErrWalletChallengePurpose):
			writeJSON(w, http.StatusConflict, map[string]any{"error": "wallet_challenge_purpose_mismatch"})
		case errors.Is(err, ErrWalletLinkChallengeMismatch):
			writeJSON(w, http.StatusForbidden, map[string]any{"error": "wallet_account_merge_user_mismatch"})
		case errors.Is(err, ErrWalletMergeSourceNotLinked):
			writeJSON(w, http.StatusConflict, map[string]any{"error": "wallet_account_merge_source_not_linked"})
		case errors.Is(err, ErrWalletMergeSameUser):
			writeJSON(w, http.StatusConflict, map[string]any{"error": "wallet_account_merge_not_required"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_account_merge_verify_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, WalletAccountMergeVerifyResponse{
		MergedWallet: result.MergedWallet,
		Wallets:      result.Wallets,
		Challenge:    result.Challenge,
		SourceUserID: result.SourceUserID,
		TargetUserID: result.TargetUserID,
	})
}

func (h HTTPHandlers) WalletDetachCheck(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok || claims == nil || claims.UserID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	var req WalletDetachCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	svc := NewWalletDetachService(h.WalletIdentities)
	result, err := svc.CheckEligibility(r.Context(), claims.UserID, req.Address)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case errors.Is(err, ErrInvalidWalletAddress):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		case errors.Is(err, ErrWalletIdentityNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "wallet_identity_not_found"})
		case errors.Is(err, ErrWalletNotOwnedByUser):
			writeJSON(w, http.StatusForbidden, map[string]any{"error": "wallet_identity_not_owned_by_user"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_detach_check_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, WalletDetachCheckResponse{
		WalletAddress:    result.WalletAddress,
		Eligible:         result.Eligible,
		IsPrimary:        result.IsPrimary,
		OwnedWalletCount: result.OwnedWalletCount,
		Reasons:          result.Reasons,
	})
}

func (h HTTPHandlers) WalletDetach(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok || claims == nil || claims.UserID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	var req WalletDetachExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	svc := NewWalletDetachService(h.WalletIdentities)
	result, err := svc.Execute(r.Context(), claims.UserID, req.Address)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case errors.Is(err, ErrInvalidWalletAddress):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		case errors.Is(err, ErrWalletIdentityNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "wallet_identity_not_found"})
		case errors.Is(err, ErrWalletNotOwnedByUser):
			writeJSON(w, http.StatusForbidden, map[string]any{"error": "wallet_identity_not_owned_by_user"})
		case errors.Is(err, ErrWalletDetachNotEligible):
			var check *WalletDetachCheckResponse
			if result != nil && result.Check != nil {
				check = &WalletDetachCheckResponse{
					WalletAddress:    result.Check.WalletAddress,
					Eligible:         result.Check.Eligible,
					IsPrimary:        result.Check.IsPrimary,
					OwnedWalletCount: result.Check.OwnedWalletCount,
					Reasons:          result.Check.Reasons,
				}
			}
			writeJSON(w, http.StatusConflict, map[string]any{"error": "wallet_detach_not_eligible", "check": check})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_detach_error"})
		}
		return
	}

	var check *WalletDetachCheckResponse
	if result != nil && result.Check != nil {
		check = &WalletDetachCheckResponse{
			WalletAddress:    result.Check.WalletAddress,
			Eligible:         result.Check.Eligible,
			IsPrimary:        result.Check.IsPrimary,
			OwnedWalletCount: result.Check.OwnedWalletCount,
			Reasons:          result.Check.Reasons,
		}
	}

	writeJSON(w, http.StatusOK, WalletDetachExecuteResponse{
		DetachedWallet: result.Detached,
		Wallets:        result.Wallets,
		Check:          check,
	})
}

func (h HTTPHandlers) WalletSetPrimary(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok || claims == nil || claims.UserID == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	var req WalletPrimarySetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	svc := NewWalletPrimaryService(h.WalletIdentities)
	result, err := svc.SetPrimary(r.Context(), claims.UserID, req.Address)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case errors.Is(err, ErrInvalidWalletAddress):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_wallet_address"})
		case errors.Is(err, ErrWalletIdentityNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "wallet_identity_not_found"})
		case errors.Is(err, ErrWalletNotOwnedByUser):
			writeJSON(w, http.StatusForbidden, map[string]any{"error": "wallet_identity_not_owned_by_user"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_primary_set_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, WalletPrimarySetResponse{
		PrimaryWallet: result.Primary,
		Wallets:       result.Wallets,
	})
}
