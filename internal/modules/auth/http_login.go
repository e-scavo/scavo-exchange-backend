package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
	authmappers "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/mappers"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"
)

type HTTPHandlers struct {
	Provider         AuthProvider
	Tokens           *coreauth.TokenService
	TTL              time.Duration
	Users            *usermod.Service
	UserSettings     *usersettingsmod.Service
	PublicBaseURL    string
	ChallengeTTL     time.Duration
	Challenges       WalletChallengeStore
	WalletIdentities WalletIdentityStore
}

// NewHTTPHandlers creates the runtime HTTP handler surface from the provider
// boundary. Legacy fields remain on HTTPHandlers only for transitional tests and
// fallback compatibility; production wiring should provide Provider explicitly.
func NewHTTPHandlers(provider AuthProvider) HTTPHandlers {
	return HTTPHandlers{Provider: provider}
}

func writeAppErrorJSON(w http.ResponseWriter, appErr *coreerrs.AppError) {
	if appErr == nil {
		appErr = coreerrs.InternalError(nil)
	}

	payload := coreerrs.NewErrorEnvelope(appErr.ToResponseError())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(appErr.Status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeErrorJSON(w http.ResponseWriter, status int, errCode string, extras ...map[string]any) {
	details := map[string]any{}
	for _, extra := range extras {
		for key, value := range extra {
			details[key] = value
		}
	}

	appErr := coreerrs.AppErrorFromLegacyAuthKey(errCode, details)
	if appErr == nil {
		appErr = coreerrs.InternalError(nil)
	}
	appErr.Status = status

	writeAppErrorJSON(w, appErr)
}

func (h HTTPHandlers) requireClaims(w http.ResponseWriter, r *http.Request) (*coreauth.Claims, bool) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok || claims == nil || claims.UserID == "" {
		writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		return nil, false
	}
	return claims, true
}

func decodeRequest(w http.ResponseWriter, r *http.Request, dst any, maxBodyBytes int64) bool {
	if err := decodeJSONBody(r, dst, maxBodyBytes); err != nil {
		writeAppErrorJSON(w, coreerrs.BadRequestInvalidPayload())
		return false
	}
	return true
}

func (h HTTPHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	input := authmappers.LoginWriteToDomainInput(req)

	resp, err := h.AuthProvider().Login(r.Context(), input.Email, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			writeAppErrorJSON(w, coreerrs.AuthInvalidCredentials())
		default:
			writeAppErrorJSON(w, coreerrs.AuthServiceError(nil))
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

	resp, err := h.AuthProvider().GetMe(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		default:
			writeAppErrorJSON(w, coreerrs.AuthServiceError(nil))
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

	var req UpdateMeRequest
	if !decodeRequest(w, r, &req, 4<<10) {
		return
	}

	input := authmappers.ProfileUpdateWriteToDomainInput(req)

	resp, err := h.AuthProvider().UpdateProfile(r.Context(), claims, input)
	if err != nil {
		switch {
		case errors.Is(err, usermod.ErrEmptyUserID):
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, usermod.ErrEmptyDisplayName):
			writeAppErrorJSON(w, coreerrs.InvalidDisplayName())
		case errors.Is(err, usermod.ErrDisplayNameTooLong):
			writeAppErrorJSON(w, coreerrs.DisplayNameTooLong())
		case errors.Is(err, usermod.ErrUserNotFound):
			writeAppErrorJSON(w, coreerrs.AuthUserNotFound())
		default:
			writeAppErrorJSON(w, coreerrs.AuthServiceError(nil))
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
func (h HTTPHandlers) UpdateMeSettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	var req UpdateMeSettingsRequest
	if !decodeRequest(w, r, &req, 32<<10) {
		return
	}

	if len(req.Preferences) == 0 {
		writeAppErrorJSON(w, coreerrs.SettingsInvalidPayload())
		return
	}

	input, err := authmappers.SettingsUpdateWriteToDomainInput(req)
	if err != nil {
		writeAppErrorJSON(w, coreerrs.SettingsInvalidPayload())
		return
	}

	resp, err := h.AuthProvider().UpdateSettings(r.Context(), claims, input)
	if err != nil {
		switch {
		case errors.Is(err, usersettingsmod.ErrUserIDRequired):
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, usersettingsmod.ErrInvalidPreferences),
			errors.Is(err, usersettingsmod.ErrNullPreferenceValue),
			errors.Is(err, usersettingsmod.ErrInvalidPreferenceValue),
			errors.Is(err, usersettingsmod.ErrIncompatiblePreference):
			writeAppErrorJSON(w, coreerrs.SettingsInvalidPayload())
		default:
			writeAppErrorJSON(w, coreerrs.AuthServiceError(nil))
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
func (h HTTPHandlers) MeSettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	resp, err := h.AuthProvider().GetSettings(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, usersettingsmod.ErrUserIDRequired):
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		default:
			writeAppErrorJSON(w, coreerrs.AuthServiceError(nil))
		}
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
func (h HTTPHandlers) Session(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	resp, err := h.AuthProvider().GetSession(r.Context(), claims)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		default:
			writeAppErrorJSON(w, coreerrs.AuthServiceError(nil))
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
