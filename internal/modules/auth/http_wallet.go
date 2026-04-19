package auth

import (
	"errors"
	"net/http"
	"time"

	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
	authapp "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/app"
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

type WalletLinkChallengeResponse = authapp.WalletLinkChallengeResponse

type WalletLinkVerifyRequest struct {
	ChallengeID string `json:"challenge_id"`
	Address     string `json:"address"`
	Signature   string `json:"signature"`
}

type WalletLinkVerifyResponse = authapp.WalletLinkVerifyResponse

type WalletAccountMergeChallengeRequest struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
}

type WalletAccountMergeChallengeResponse = authapp.WalletAccountMergeChallengeResponse

type WalletAccountMergeVerifyRequest struct {
	ChallengeID string `json:"challenge_id"`
	Address     string `json:"address"`
	Signature   string `json:"signature"`
}

type WalletAccountMergeVerifyResponse = authapp.WalletAccountMergeVerifyResponse

type WalletDetachCheckRequest struct {
	Address string `json:"wallet_address"`
}

type WalletDetachCheckResponse = authapp.WalletDetachCheckResponse

type WalletDetachExecuteRequest struct {
	Address string `json:"wallet_address"`
}

type WalletDetachExecuteResponse = authapp.WalletDetachExecuteResponse

type WalletPrimarySetRequest struct {
	Address string `json:"wallet_address"`
}

type WalletPrimarySetResponse = authapp.WalletPrimarySetResponse

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
			writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
		default:
			writeAppErrorJSON(w, coreerrs.WalletChallengeError(nil))
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
			writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
		case errors.Is(err, ErrInvalidWalletSignature):
			writeAppErrorJSON(w, coreerrs.WalletInvalidSignature())
		case errors.Is(err, ErrWalletChallengeNotFound):
			writeAppErrorJSON(w, coreerrs.WalletChallengeNotFound())
		case errors.Is(err, ErrChallengeExpired):
			writeAppErrorJSON(w, coreerrs.WalletChallengeExpired())
		case errors.Is(err, ErrChallengeUsed):
			writeAppErrorJSON(w, coreerrs.WalletChallengeUsed())
		case errors.Is(err, ErrWalletChallengePurpose):
			writeAppErrorJSON(w, coreerrs.WalletChallengePurposeMismatch())
		default:
			writeAppErrorJSON(w, coreerrs.WalletVerifyError(nil))
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
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, ErrInvalidWalletAddress):
			writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
		default:
			writeAppErrorJSON(w, coreerrs.WalletLinkChallengeError(nil))
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
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, ErrInvalidWalletAddress):
			writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
		case errors.Is(err, ErrInvalidWalletSignature):
			writeAppErrorJSON(w, coreerrs.WalletInvalidSignature())
		case errors.Is(err, ErrWalletChallengeNotFound):
			writeAppErrorJSON(w, coreerrs.WalletChallengeNotFound())
		case errors.Is(err, ErrChallengeExpired):
			writeAppErrorJSON(w, coreerrs.WalletChallengeExpired())
		case errors.Is(err, ErrChallengeUsed):
			writeAppErrorJSON(w, coreerrs.WalletChallengeUsed())
		case errors.Is(err, ErrWalletChallengePurpose):
			writeAppErrorJSON(w, coreerrs.WalletChallengePurposeMismatch())
		case errors.Is(err, ErrWalletLinkChallengeMismatch):
			writeAppErrorJSON(w, coreerrs.WalletLinkChallengeUserMismatch())
		case errors.Is(err, ErrWalletIdentityAlreadyLinked):
			writeAppErrorJSON(w, coreerrs.WalletAlreadyLinked())
		case errors.Is(err, ErrWalletAlreadyLinkedToUser):
			writeAppErrorJSON(w, coreerrs.WalletAlreadyLinkedToUser())
		default:
			writeAppErrorJSON(w, coreerrs.WalletLinkVerifyError(nil))
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
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, ErrInvalidWalletAddress):
			writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
		default:
			writeAppErrorJSON(w, coreerrs.WalletAccountMergeChallengeError(nil))
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
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, ErrInvalidWalletAddress):
			writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
		case errors.Is(err, ErrInvalidWalletSignature):
			writeAppErrorJSON(w, coreerrs.WalletInvalidSignature())
		case errors.Is(err, ErrWalletChallengeNotFound):
			writeAppErrorJSON(w, coreerrs.WalletChallengeNotFound())
		case errors.Is(err, ErrChallengeExpired):
			writeAppErrorJSON(w, coreerrs.WalletChallengeExpired())
		case errors.Is(err, ErrChallengeUsed):
			writeAppErrorJSON(w, coreerrs.WalletChallengeUsed())
		case errors.Is(err, ErrWalletChallengePurpose):
			writeAppErrorJSON(w, coreerrs.WalletChallengePurposeMismatch())
		case errors.Is(err, ErrWalletLinkChallengeMismatch):
			writeAppErrorJSON(w, coreerrs.WalletAccountMergeUserMismatch())
		case errors.Is(err, ErrWalletMergeSourceNotLinked):
			writeAppErrorJSON(w, coreerrs.WalletAccountMergeSourceNotLinked())
		case errors.Is(err, ErrWalletMergeSameUser):
			writeAppErrorJSON(w, coreerrs.WalletAccountMergeNotRequired())
		default:
			writeAppErrorJSON(w, coreerrs.WalletAccountMergeVerifyError(nil))
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
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, ErrInvalidWalletAddress):
			writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
		case errors.Is(err, ErrWalletIdentityNotFound):
			writeAppErrorJSON(w, coreerrs.WalletNotFound())
		case errors.Is(err, ErrWalletNotOwnedByUser):
			writeAppErrorJSON(w, coreerrs.WalletNotOwnedByUser())
		default:
			writeAppErrorJSON(w, coreerrs.WalletDetachCheckError(nil))
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
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, ErrInvalidWalletAddress):
			writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
		case errors.Is(err, ErrWalletIdentityNotFound):
			writeAppErrorJSON(w, coreerrs.WalletNotFound())
		case errors.Is(err, ErrWalletNotOwnedByUser):
			writeAppErrorJSON(w, coreerrs.WalletNotOwnedByUser())
		case errors.Is(err, ErrWalletDetachNotEligible):
			writeAppErrorJSON(w, coreerrs.WalletCannotDetach(map[string]any{"check": response.Check}))
		default:
			writeAppErrorJSON(w, coreerrs.WalletDetachError(nil))
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
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, ErrInvalidWalletAddress):
			writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
		case errors.Is(err, ErrWalletIdentityNotFound):
			writeAppErrorJSON(w, coreerrs.WalletNotFound())
		case errors.Is(err, ErrWalletNotOwnedByUser):
			writeAppErrorJSON(w, coreerrs.WalletNotOwnedByUser())
		default:
			writeAppErrorJSON(w, coreerrs.WalletPrimarySetError(nil))
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}
