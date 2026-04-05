package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	UserID      string `json:"user_id"`
}

type UpdateMeRequest struct {
	DisplayName string `json:"display_name"`
}

type MeResponse struct {
	User    *usermod.User `json:"user"`
	Profile *ProfileView  `json:"profile,omitempty"`
}

type SessionResponse struct {
	Session *SessionView `json:"session"`
}

type HTTPHandlers struct {
	Tokens           *coreauth.TokenService
	TTL              time.Duration
	Users            *usermod.Service
	PublicBaseURL    string
	ChallengeTTL     time.Duration
	Challenges       WalletChallengeStore
	WalletIdentities WalletIdentityStore
}

func (h HTTPHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	svc := NewService(h.Tokens, h.Users, h.TTL)
	result, err := svc.LoginDev(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid_credentials"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	userID := ""
	if result.User != nil {
		userID = result.User.ID
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresIn:   result.ExpiresIn,
		UserID:      userID,
	})
}

func (h HTTPHandlers) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	profile, err := buildProfileView(r.Context(), claims, h.Users, h.WalletIdentities)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, MeResponse{
		User:    profile.User,
		Profile: profile,
	})
}

func (h HTTPHandlers) UpdateMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}
	if h.Users == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		return
	}

	var req UpdateMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	updatedUser, err := h.Users.UpdateDisplayName(r.Context(), claims.UserID, req.DisplayName)
	if err != nil {
		switch err.Error() {
		case "empty user id":
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case "empty display name":
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_display_name"})
		case "display name too long":
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "display_name_too_long"})
		default:
			switch {
			case errors.Is(err, usermod.ErrUserNotFound):
				writeJSON(w, http.StatusNotFound, map[string]any{"error": "user_not_found"})
			default:
				writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
			}
		}
		return
	}

	profile, err := buildProfileViewWithUser(r.Context(), claims, updatedUser, h.WalletIdentities)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		return
	}

	writeJSON(w, http.StatusOK, MeResponse{
		User:    profile.User,
		Profile: profile,
	})
}

func (h HTTPHandlers) Session(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	svc := NewService(h.Tokens, h.Users, h.TTL)
	session, err := svc.ResolveSessionClaims(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, SessionResponse{Session: session})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
