package auth

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
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

type errorEnvelope struct {
	Error coreerrs.ResponseError `json:"error"`
}

func writeErrorJSON(w http.ResponseWriter, status int, errCode string, extras ...map[string]any) {
	details := map[string]any{}
	for _, extra := range extras {
		for key, value := range extra {
			details[key] = value
		}
	}

	spec := normalizeAuthError(errCode)

	payload := errorEnvelope{
		Error: coreerrs.NewResponseError(spec.Code, spec.Message, details),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

type authErrorSpec struct {
	Code    string
	Message string
}

var authErrorCatalog = map[string]authErrorSpec{
	"auth_service_error":                     {Code: "AUTH_SERVICE_ERROR", Message: "authentication service error"},
	"bad_request":                            {Code: "BAD_REQUEST", Message: "invalid request payload"},
	"display_name_too_long":                  {Code: "DISPLAY_NAME_TOO_LONG", Message: "display name is too long"},
	"invalid_credentials":                    {Code: "AUTH_INVALID_CREDENTIALS", Message: "invalid credentials"},
	"invalid_display_name":                   {Code: "INVALID_DISPLAY_NAME", Message: "invalid display name"},
	"invalid_limit":                          {Code: "INVALID_LIMIT", Message: "invalid limit"},
	"invalid_offset":                         {Code: "INVALID_OFFSET", Message: "invalid offset"},
	"invalid_order":                          {Code: "INVALID_ORDER", Message: "invalid order"},
	"invalid_order_requires_sort":            {Code: "INVALID_ORDER_REQUIRES_SORT", Message: "order requires sort"},
	"invalid_preferences":                    {Code: "SETTINGS_INVALID_PAYLOAD", Message: "invalid preferences payload"},
	"invalid_primary":                        {Code: "INVALID_PRIMARY", Message: "invalid primary filter"},
	"invalid_sort":                           {Code: "INVALID_SORT", Message: "invalid sort"},
	"invalid_status":                         {Code: "INVALID_STATUS", Message: "invalid status"},
	"invalid_wallet_address":                 {Code: "WALLET_INVALID_ADDRESS", Message: "invalid wallet address"},
	"invalid_wallet_signature":               {Code: "WALLET_INVALID_SIGNATURE", Message: "invalid wallet signature"},
	"unauthorized":                           {Code: "AUTH_UNAUTHORIZED", Message: "authentication required"},
	"user_not_found":                         {Code: "AUTH_USER_NOT_FOUND", Message: "user not found"},
	"wallet_account_merge_challenge_error":   {Code: "WALLET_ACCOUNT_MERGE_CHALLENGE_ERROR", Message: "wallet account merge challenge error"},
	"wallet_account_merge_not_required":      {Code: "WALLET_ACCOUNT_MERGE_NOT_REQUIRED", Message: "wallet account merge is not required"},
	"wallet_account_merge_source_not_linked": {Code: "WALLET_ACCOUNT_MERGE_SOURCE_NOT_LINKED", Message: "wallet account merge source is not linked"},
	"wallet_account_merge_user_mismatch":     {Code: "WALLET_ACCOUNT_MERGE_USER_MISMATCH", Message: "wallet account merge user mismatch"},
	"wallet_account_merge_verify_error":      {Code: "WALLET_ACCOUNT_MERGE_VERIFY_ERROR", Message: "wallet account merge verification error"},
	"wallet_challenge_error":                 {Code: "WALLET_CHALLENGE_ERROR", Message: "wallet challenge error"},
	"wallet_challenge_expired":               {Code: "WALLET_CHALLENGE_EXPIRED", Message: "wallet challenge expired"},
	"wallet_challenge_not_found":             {Code: "WALLET_CHALLENGE_NOT_FOUND", Message: "wallet challenge not found"},
	"wallet_challenge_purpose_mismatch":      {Code: "WALLET_CHALLENGE_PURPOSE_MISMATCH", Message: "wallet challenge purpose mismatch"},
	"wallet_challenge_used":                  {Code: "WALLET_CHALLENGE_USED", Message: "wallet challenge already used"},
	"wallet_detach_check_error":              {Code: "WALLET_DETACH_CHECK_ERROR", Message: "wallet detach check error"},
	"wallet_detach_error":                    {Code: "WALLET_DETACH_ERROR", Message: "wallet detach error"},
	"wallet_detach_not_eligible":             {Code: "WALLET_CANNOT_DETACH", Message: "wallet cannot be detached under current ownership rules"},
	"wallet_identity_already_linked":         {Code: "WALLET_ALREADY_LINKED", Message: "wallet identity already linked"},
	"wallet_identity_already_linked_to_user": {Code: "WALLET_ALREADY_LINKED_TO_USER", Message: "wallet identity already linked to user"},
	"wallet_identity_error":                  {Code: "WALLET_IDENTITY_ERROR", Message: "wallet identity error"},
	"wallet_identity_not_found":              {Code: "WALLET_NOT_FOUND", Message: "wallet identity not found"},
	"wallet_identity_not_owned_by_user":      {Code: "WALLET_NOT_OWNED_BY_USER", Message: "wallet identity not owned by user"},
	"wallet_link_challenge_error":            {Code: "WALLET_LINK_CHALLENGE_ERROR", Message: "wallet link challenge error"},
	"wallet_link_challenge_user_mismatch":    {Code: "WALLET_LINK_CHALLENGE_USER_MISMATCH", Message: "wallet link challenge user mismatch"},
	"wallet_link_verify_error":               {Code: "WALLET_LINK_VERIFY_ERROR", Message: "wallet link verification error"},
	"wallet_primary_set_error":               {Code: "WALLET_PRIMARY_SET_ERROR", Message: "wallet primary set error"},
	"wallet_verify_error":                    {Code: "WALLET_VERIFY_ERROR", Message: "wallet verification error"},
}

func normalizeAuthError(errCode string) authErrorSpec {
	if spec, ok := authErrorCatalog[errCode]; ok {
		return spec
	}

	message := strings.ReplaceAll(strings.TrimSpace(errCode), "_", " ")
	if message == "" {
		message = "unexpected error"
	}

	code := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(errCode), "-", "_"))
	code = strings.ReplaceAll(code, " ", "_")
	if code == "" {
		code = "INTERNAL_ERROR"
	}

	return authErrorSpec{Code: code, Message: message}
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
