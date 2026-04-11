package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"
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

type UpdateMeSettingsRequest struct {
	Preferences json.RawMessage `json:"preferences"`
}

type MeSettingsResponse struct {
	Settings usersettingsmod.View `json:"settings"`
}

type HTTPHandlers struct {
	Tokens           *coreauth.TokenService
	TTL              time.Duration
	Users            *usermod.Service
	UserSettings     *usersettingsmod.Service
	PublicBaseURL    string
	ChallengeTTL     time.Duration
	Challenges       WalletChallengeStore
	WalletIdentities WalletIdentityStore
}

func (h HTTPHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSONBody(r, &req, 4<<10); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	resp, err := h.Application().Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid_credentials"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h HTTPHandlers) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	resp, err := h.Application().GetMe(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
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
	if err := decodeJSONBody(r, &req, 4<<10); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	updatedUser, err := h.Users.UpdateDisplayName(r.Context(), claims.UserID, req.DisplayName)
	if err != nil {
		switch {
		case errors.Is(err, usermod.ErrEmptyUserID):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case errors.Is(err, usermod.ErrEmptyDisplayName):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_display_name"})
		case errors.Is(err, usermod.ErrDisplayNameTooLong):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "display_name_too_long"})
		case errors.Is(err, usermod.ErrUserNotFound):
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "user_not_found"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
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

func (h HTTPHandlers) UpdateMeSettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	if h.UserSettings == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		return
	}

	var req UpdateMeSettingsRequest
	if err := decodeJSONBody(r, &req, 32<<10); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	if len(req.Preferences) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_preferences"})
		return
	}

	var patch map[string]any
	if err := json.Unmarshal(req.Preferences, &patch); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_preferences"})
		return
	}

	settings, err := h.UserSettings.UpdatePreferences(r.Context(), claims.UserID, patch)
	if err != nil {
		switch {
		case errors.Is(err, usersettingsmod.ErrUserIDRequired):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		case errors.Is(err, usersettingsmod.ErrInvalidPreferences),
			errors.Is(err, usersettingsmod.ErrNullPreferenceValue),
			errors.Is(err, usersettingsmod.ErrInvalidPreferenceValue),
			errors.Is(err, usersettingsmod.ErrIncompatiblePreference):
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid_preferences"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	view := usersettingsmod.ToView(settings)

	writeJSON(w, http.StatusOK, MeSettingsResponse{
		Settings: view,
	})
}

func (h HTTPHandlers) MeSettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	if h.UserSettings == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		return
	}

	settings, err := h.UserSettings.GetOrDefault(r.Context(), claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, usersettingsmod.ErrUserIDRequired):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	view := usersettingsmod.ToView(settings)

	writeJSON(w, http.StatusOK, MeSettingsResponse{
		Settings: view,
	})
}

func (h HTTPHandlers) Session(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	resp, err := h.Application().GetSession(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "auth_service_error"})
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func decodeJSONBody(r *http.Request, dst any, maxBodyBytes int64) error {
	if r == nil {
		return errors.New("nil request")
	}

	body := r.Body
	if maxBodyBytes > 0 {
		body = http.MaxBytesReader(nil, r.Body, maxBodyBytes)
	}

	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}

	if err := dec.Decode(&struct{}{}); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return errors.New("unexpected trailing data")
	}

	return errors.New("unexpected trailing data")
}
