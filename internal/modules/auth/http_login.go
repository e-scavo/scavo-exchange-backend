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

func writeErrorJSON(w http.ResponseWriter, code int, errCode string, extras ...map[string]any) {
	payload := map[string]any{"error": errCode}
	for _, extra := range extras {
		for key, value := range extra {
			payload[key] = value
		}
	}
	writeJSON(w, code, payload)
}

func (h HTTPHandlers) requireClaims(w http.ResponseWriter, r *http.Request) (*coreauth.Claims, bool) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok || claims == nil || claims.UserID == "" {
		writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		return nil, false
	}
	return claims, true
}

func decodeRequest(w http.ResponseWriter, r *http.Request, dst any, maxBodyBytes int64) bool {
	if err := decodeJSONBody(r, dst, maxBodyBytes); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "bad_request")
		return false
	}
	return true
}

func (h HTTPHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	resp, err := h.Application().Login(r.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			writeErrorJSON(w, http.StatusUnauthorized, "invalid_credentials")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "auth_service_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h HTTPHandlers) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	resp, err := h.Application().GetMe(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "auth_service_error")
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h HTTPHandlers) UpdateMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}
	if h.Users == nil {
		writeErrorJSON(w, http.StatusInternalServerError, "auth_service_error")
		return
	}

	var req UpdateMeRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	updatedUser, err := h.Users.UpdateDisplayName(r.Context(), claims.UserID, req.DisplayName)
	if err != nil {
		switch {
		case errors.Is(err, usermod.ErrEmptyUserID):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, usermod.ErrEmptyDisplayName):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_display_name")
		case errors.Is(err, usermod.ErrDisplayNameTooLong):
			writeErrorJSON(w, http.StatusBadRequest, "display_name_too_long")
		case errors.Is(err, usermod.ErrUserNotFound):
			writeErrorJSON(w, http.StatusNotFound, "user_not_found")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "auth_service_error")
		}
		return
	}

	profile, err := buildProfileViewWithUser(r.Context(), claims, updatedUser, h.WalletIdentities)
	if err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "auth_service_error")
		return
	}

	writeJSON(w, http.StatusOK, MeResponse{
		User:    profile.User,
		Profile: profile,
	})
}

func (h HTTPHandlers) UpdateMeSettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	if h.UserSettings == nil {
		writeErrorJSON(w, http.StatusInternalServerError, "auth_service_error")
		return
	}

	var req UpdateMeSettingsRequest
	if !decodeRequest(w, r, &req, 32<<10) {
		return
	}

	if len(req.Preferences) == 0 {
		writeErrorJSON(w, http.StatusBadRequest, "invalid_preferences")
		return
	}

	var patch map[string]any
	if err := json.Unmarshal(req.Preferences, &patch); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, "invalid_preferences")
		return
	}

	settings, err := h.UserSettings.UpdatePreferences(r.Context(), claims.UserID, patch)
	if err != nil {
		switch {
		case errors.Is(err, usersettingsmod.ErrUserIDRequired):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		case errors.Is(err, usersettingsmod.ErrInvalidPreferences),
			errors.Is(err, usersettingsmod.ErrNullPreferenceValue),
			errors.Is(err, usersettingsmod.ErrInvalidPreferenceValue),
			errors.Is(err, usersettingsmod.ErrIncompatiblePreference):
			writeErrorJSON(w, http.StatusBadRequest, "invalid_preferences")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "auth_service_error")
		}
		return
	}

	view := usersettingsmod.ToView(settings)

	writeJSON(w, http.StatusOK, MeSettingsResponse{
		Settings: view,
	})
}

func (h HTTPHandlers) MeSettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	if h.UserSettings == nil {
		writeErrorJSON(w, http.StatusInternalServerError, "auth_service_error")
		return
	}

	settings, err := h.UserSettings.GetOrDefault(r.Context(), claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, usersettingsmod.ErrUserIDRequired):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "auth_service_error")
		}
		return
	}

	view := usersettingsmod.ToView(settings)

	writeJSON(w, http.StatusOK, MeSettingsResponse{
		Settings: view,
	})
}

func (h HTTPHandlers) Session(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	resp, err := h.Application().GetSession(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeErrorJSON(w, http.StatusUnauthorized, "unauthorized")
		default:
			writeErrorJSON(w, http.StatusInternalServerError, "auth_service_error")
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
