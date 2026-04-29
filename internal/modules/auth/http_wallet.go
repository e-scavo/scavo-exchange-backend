package auth

import (
	"errors"
	"net/http"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
	authapp "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/app"
	authmappers "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/mappers"
	authreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/readmodels"
	authwritemodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/writemodels"
	usermappers "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/mappers"
	userreadmodels "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/readmodels"
)

// =====================================================
// WALLET BOOTSTRAP AUTH (PUBLIC - NO AUTH REQUIRED)
// =====================================================

type WalletChallengeRequest = authwritemodels.AuthWalletChallengeWriteModel

type WalletChallengeResponse struct {
	Challenge *authreadmodels.AuthWalletChallengeReadModel `json:"challenge"`
}

type WalletVerifyRequest = authwritemodels.AuthWalletVerifyWriteModel

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

func (h HTTPHandlers) WalletChallenge(w http.ResponseWriter, r *http.Request) {
	var req WalletChallengeRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	challengeTTL := h.ChallengeTTL
	if challengeTTL <= 0 {
		challengeTTL = 5 * time.Minute
	}

	input := authmappers.WalletChallengeWriteToDomainInput(req)

	svc := NewWalletChallengeService(h.Challenges, h.PublicBaseURL, challengeTTL)
	challenge, err := svc.Create(r.Context(), input.Address, input.Chain)
	if err != nil {
		writeWalletChallengeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, WalletChallengeResponse{Challenge: authmappers.WalletChallengeToReadModel(challenge)})
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

	input := authmappers.WalletVerifyWriteToDomainInput(req)

	result, challenge, err := verifySvc.VerifyAndLogin(r.Context(), input.ChallengeID, input.Address, input.Signature)
	if err != nil {
		writeWalletVerifyError(w, err)
		return
	}

	userID := ""
	var user *userreadmodels.UserReadModel
	if result != nil && result.User != nil {
		userID = result.User.ID
		user = usermappers.UserToReadModel(result.User)
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
		Challenge:     authmappers.WalletChallengeToReadModel(challenge),
	})
}

// =====================================================
// WALLET MANAGEMENT (AUTHENTICATED)
// =====================================================

type WalletLinkChallengeRequest = authwritemodels.AuthWalletLinkChallengeWriteModel

type WalletLinkChallengeResponse = authapp.WalletLinkChallengeResponse

type WalletLinkVerifyRequest = authwritemodels.AuthWalletLinkVerifyWriteModel

type WalletLinkVerifyResponse = authapp.WalletLinkVerifyResponse

type WalletAccountMergeVerifyRequest = authwritemodels.AuthWalletAccountMergeVerifyWriteModel

type WalletAccountMergeVerifyResponse = authapp.WalletAccountMergeVerifyResponse

type WalletDetachCheckRequest = authwritemodels.AuthWalletDetachCheckWriteModel

type WalletDetachCheckResponse = authapp.WalletDetachCheckResponse

type WalletDetachExecuteRequest = authwritemodels.AuthWalletDetachExecuteWriteModel

type WalletDetachExecuteResponse = authapp.WalletDetachExecuteResponse

type WalletPrimarySetRequest = authwritemodels.AuthWalletPrimarySetWriteModel

type WalletPrimarySetResponse = authapp.WalletPrimarySetResponse

// (handlers autenticados quedan EXACTAMENTE como los tenías, ya están bien)

type WalletAccountMergeChallengeRequest = authwritemodels.AuthWalletAccountMergeChallengeWriteModel

type WalletAccountMergeChallengeResponse = authapp.WalletAccountMergeChallengeResponse

func (h HTTPHandlers) WalletLinkChallenge(w http.ResponseWriter, r *http.Request) {
	claims, req, ok := decodeAuthenticatedWalletRequest[WalletLinkChallengeRequest](h, w, r, 4<<10)
	if !ok {
		return
	}

	input := authmappers.WalletLinkChallengeWriteToDomainInput(req)

	response, err := h.Application().CreateWalletLinkChallenge(r.Context(), claims.UserID, input.Address, input.Chain)
	if err != nil {
		writeWalletLinkChallengeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletLinkVerify(w http.ResponseWriter, r *http.Request) {
	claims, req, ok := decodeAuthenticatedWalletRequest[WalletLinkVerifyRequest](h, w, r, 8<<10)
	if !ok {
		return
	}

	input := authmappers.WalletLinkVerifyWriteToDomainInput(req)

	response, err := h.Application().VerifyWalletLink(r.Context(), claims.UserID, input.ChallengeID, input.Address, input.Signature)
	if err != nil {
		writeWalletLinkVerifyError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletAccountMergeChallenge(w http.ResponseWriter, r *http.Request) {
	claims, req, ok := decodeAuthenticatedWalletRequest[WalletAccountMergeChallengeRequest](h, w, r, 4<<10)
	if !ok {
		return
	}

	input := authmappers.WalletAccountMergeChallengeWriteToDomainInput(req)

	response, err := h.Application().CreateWalletAccountMergeChallenge(r.Context(), claims.UserID, input.Address, input.Chain)
	if err != nil {
		writeWalletAccountMergeChallengeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletAccountMergeVerify(w http.ResponseWriter, r *http.Request) {
	claims, req, ok := decodeAuthenticatedWalletRequest[WalletAccountMergeVerifyRequest](h, w, r, 8<<10)
	if !ok {
		return
	}

	input := authmappers.WalletAccountMergeVerifyWriteToDomainInput(req)

	response, err := h.Application().VerifyWalletAccountMerge(r.Context(), claims.UserID, input.ChallengeID, input.Address, input.Signature)
	if err != nil {
		writeWalletAccountMergeVerifyError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletDetachCheck(w http.ResponseWriter, r *http.Request) {
	claims, req, ok := decodeAuthenticatedWalletRequest[WalletDetachCheckRequest](h, w, r, 4<<10)
	if !ok {
		return
	}

	input := authmappers.WalletDetachCheckWriteToDomainInput(req)

	response, err := h.Application().CheckWalletDetach(r.Context(), claims.UserID, input.Address)
	if err != nil {
		writeWalletDetachCheckError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletDetach(w http.ResponseWriter, r *http.Request) {
	claims, req, ok := decodeAuthenticatedWalletRequest[WalletDetachExecuteRequest](h, w, r, 4<<10)
	if !ok {
		return
	}

	input := authmappers.WalletDetachExecuteWriteToDomainInput(req)

	response, err := h.Application().ExecuteWalletDetach(r.Context(), claims.UserID, input.Address)
	if err != nil {
		writeWalletDetachError(w, err, response)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func (h HTTPHandlers) WalletSetPrimary(w http.ResponseWriter, r *http.Request) {
	claims, req, ok := decodeAuthenticatedWalletRequest[WalletPrimarySetRequest](h, w, r, 4<<10)
	if !ok {
		return
	}

	input := authmappers.WalletPrimarySetWriteToDomainInput(req)

	response, err := h.Application().SetPrimaryWallet(r.Context(), claims.UserID, input.Address)
	if err != nil {
		writeWalletPrimarySetError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, response)
}

func decodeAuthenticatedWalletRequest[T any](h HTTPHandlers, w http.ResponseWriter, r *http.Request, maxBodyBytes int64) (*coreauth.Claims, T, bool) {
	var req T

	claims, ok := h.requireClaims(w, r)
	if !ok {
		return nil, req, false
	}

	if !decodeRequest(w, r, &req, maxBodyBytes) {
		return nil, req, false
	}

	return claims, req, true
}

func writeWalletChallengeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidWalletAddress):
		writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
	default:
		writeAppErrorJSON(w, coreerrs.WalletChallengeError(nil))
	}
}

func writeWalletVerifyError(w http.ResponseWriter, err error) {
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
}

func writeWalletLinkChallengeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUnauthorized):
		writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
	case errors.Is(err, ErrInvalidWalletAddress):
		writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
	default:
		writeAppErrorJSON(w, coreerrs.WalletLinkChallengeError(nil))
	}
}

func writeWalletLinkVerifyError(w http.ResponseWriter, err error) {
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
}

func writeWalletAccountMergeChallengeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUnauthorized):
		writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
	case errors.Is(err, ErrInvalidWalletAddress):
		writeAppErrorJSON(w, coreerrs.WalletInvalidAddress())
	default:
		writeAppErrorJSON(w, coreerrs.WalletAccountMergeChallengeError(nil))
	}
}

func writeWalletAccountMergeVerifyError(w http.ResponseWriter, err error) {
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
}

func writeWalletDetachCheckError(w http.ResponseWriter, err error) {
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
}

func writeWalletDetachError(w http.ResponseWriter, err error, response WalletDetachExecuteResponse) {
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
}

func writeWalletPrimarySetError(w http.ResponseWriter, err error) {
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
}
