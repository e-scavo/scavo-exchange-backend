package errs

import (
	"net/http"
	"strings"
)

type AppErrorSpec struct {
	LegacyKey string
	Code      string
	Message   string
	Status    int
	Category  Category
}

var authErrorCatalog = map[string]AppErrorSpec{
	"auth_not_configured":                    {LegacyKey: "auth_not_configured", Code: "AUTH_NOT_CONFIGURED", Message: "authentication is not configured", Status: http.StatusInternalServerError, Category: CategoryAuth},
	"auth_missing_bearer_token":              {LegacyKey: "auth_missing_bearer_token", Code: "AUTH_MISSING_BEARER_TOKEN", Message: "missing bearer token", Status: http.StatusUnauthorized, Category: CategoryAuth},
	"auth_service_error":                     {LegacyKey: "auth_service_error", Code: "AUTH_SERVICE_ERROR", Message: "authentication service error", Status: http.StatusInternalServerError, Category: CategoryAuth},
	"bad_request":                            {LegacyKey: "bad_request", Code: "BAD_REQUEST", Message: "invalid request payload", Status: http.StatusBadRequest, Category: CategoryGeneric},
	"display_name_too_long":                  {LegacyKey: "display_name_too_long", Code: "DISPLAY_NAME_TOO_LONG", Message: "display name is too long", Status: http.StatusBadRequest, Category: CategoryAuth},
	"internal_error":                         {LegacyKey: "internal_error", Code: "INTERNAL_ERROR", Message: "internal server error", Status: http.StatusInternalServerError, Category: CategoryInternal},
	"invalid_credentials":                    {LegacyKey: "invalid_credentials", Code: "AUTH_INVALID_CREDENTIALS", Message: "invalid credentials", Status: http.StatusUnauthorized, Category: CategoryAuth},
	"invalid_display_name":                   {LegacyKey: "invalid_display_name", Code: "INVALID_DISPLAY_NAME", Message: "invalid display name", Status: http.StatusBadRequest, Category: CategoryAuth},
	"invalid_limit":                          {LegacyKey: "invalid_limit", Code: "INVALID_LIMIT", Message: "invalid limit", Status: http.StatusBadRequest, Category: CategoryWallet},
	"invalid_offset":                         {LegacyKey: "invalid_offset", Code: "INVALID_OFFSET", Message: "invalid offset", Status: http.StatusBadRequest, Category: CategoryWallet},
	"invalid_order":                          {LegacyKey: "invalid_order", Code: "INVALID_ORDER", Message: "invalid order", Status: http.StatusBadRequest, Category: CategoryWallet},
	"invalid_order_requires_sort":            {LegacyKey: "invalid_order_requires_sort", Code: "INVALID_ORDER_REQUIRES_SORT", Message: "order requires sort", Status: http.StatusBadRequest, Category: CategoryWallet},
	"invalid_preferences":                    {LegacyKey: "invalid_preferences", Code: "SETTINGS_INVALID_PAYLOAD", Message: "invalid preferences payload", Status: http.StatusBadRequest, Category: CategorySettings},
	"invalid_primary":                        {LegacyKey: "invalid_primary", Code: "INVALID_PRIMARY", Message: "invalid primary filter", Status: http.StatusBadRequest, Category: CategoryWallet},
	"invalid_sort":                           {LegacyKey: "invalid_sort", Code: "INVALID_SORT", Message: "invalid sort", Status: http.StatusBadRequest, Category: CategoryWallet},
	"invalid_status":                         {LegacyKey: "invalid_status", Code: "INVALID_STATUS", Message: "invalid status", Status: http.StatusBadRequest, Category: CategoryWallet},
	"invalid_wallet_address":                 {LegacyKey: "invalid_wallet_address", Code: "WALLET_INVALID_ADDRESS", Message: "invalid wallet address", Status: http.StatusBadRequest, Category: CategoryWallet},
	"invalid_wallet_signature":               {LegacyKey: "invalid_wallet_signature", Code: "WALLET_INVALID_SIGNATURE", Message: "invalid wallet signature", Status: http.StatusUnauthorized, Category: CategoryWallet},
	"timeout":                                {LegacyKey: "timeout", Code: "TIMEOUT", Message: "request timed out", Status: http.StatusServiceUnavailable, Category: CategoryGeneric},
	"unauthorized":                           {LegacyKey: "unauthorized", Code: "AUTH_UNAUTHORIZED", Message: "authentication required", Status: http.StatusUnauthorized, Category: CategoryAuth},
	"user_not_found":                         {LegacyKey: "user_not_found", Code: "AUTH_USER_NOT_FOUND", Message: "user not found", Status: http.StatusNotFound, Category: CategoryAuth},
	"wallet_account_merge_challenge_error":   {LegacyKey: "wallet_account_merge_challenge_error", Code: "WALLET_ACCOUNT_MERGE_CHALLENGE_ERROR", Message: "wallet account merge challenge error", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_account_merge_not_required":      {LegacyKey: "wallet_account_merge_not_required", Code: "WALLET_ACCOUNT_MERGE_NOT_REQUIRED", Message: "wallet account merge is not required", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_account_merge_source_not_linked": {LegacyKey: "wallet_account_merge_source_not_linked", Code: "WALLET_ACCOUNT_MERGE_SOURCE_NOT_LINKED", Message: "wallet account merge source is not linked", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_account_merge_user_mismatch":     {LegacyKey: "wallet_account_merge_user_mismatch", Code: "WALLET_ACCOUNT_MERGE_USER_MISMATCH", Message: "wallet account merge user mismatch", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_account_merge_verify_error":      {LegacyKey: "wallet_account_merge_verify_error", Code: "WALLET_ACCOUNT_MERGE_VERIFY_ERROR", Message: "wallet account merge verification error", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_challenge_error":                 {LegacyKey: "wallet_challenge_error", Code: "WALLET_CHALLENGE_ERROR", Message: "wallet challenge error", Status: http.StatusInternalServerError, Category: CategoryWallet},
	"wallet_challenge_expired":               {LegacyKey: "wallet_challenge_expired", Code: "WALLET_CHALLENGE_EXPIRED", Message: "wallet challenge expired", Status: http.StatusUnauthorized, Category: CategoryWallet},
	"wallet_challenge_not_found":             {LegacyKey: "wallet_challenge_not_found", Code: "WALLET_CHALLENGE_NOT_FOUND", Message: "wallet challenge not found", Status: http.StatusNotFound, Category: CategoryWallet},
	"wallet_challenge_purpose_mismatch":      {LegacyKey: "wallet_challenge_purpose_mismatch", Code: "WALLET_CHALLENGE_PURPOSE_MISMATCH", Message: "wallet challenge purpose mismatch", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_challenge_used":                  {LegacyKey: "wallet_challenge_used", Code: "WALLET_CHALLENGE_USED", Message: "wallet challenge already used", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_detach_check_error":              {LegacyKey: "wallet_detach_check_error", Code: "WALLET_DETACH_CHECK_ERROR", Message: "wallet detach check error", Status: http.StatusInternalServerError, Category: CategoryWallet},
	"wallet_detach_error":                    {LegacyKey: "wallet_detach_error", Code: "WALLET_DETACH_ERROR", Message: "wallet detach error", Status: http.StatusInternalServerError, Category: CategoryWallet},
	"wallet_detach_not_eligible":             {LegacyKey: "wallet_detach_not_eligible", Code: "WALLET_CANNOT_DETACH", Message: "wallet cannot be detached under current ownership rules", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_identity_already_linked":         {LegacyKey: "wallet_identity_already_linked", Code: "WALLET_ALREADY_LINKED", Message: "wallet identity already linked", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_identity_already_linked_to_user": {LegacyKey: "wallet_identity_already_linked_to_user", Code: "WALLET_ALREADY_LINKED_TO_USER", Message: "wallet identity already linked to user", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_identity_error":                  {LegacyKey: "wallet_identity_error", Code: "WALLET_IDENTITY_ERROR", Message: "wallet identity error", Status: http.StatusInternalServerError, Category: CategoryWallet},
	"wallet_identity_not_found":              {LegacyKey: "wallet_identity_not_found", Code: "WALLET_NOT_FOUND", Message: "wallet identity not found", Status: http.StatusNotFound, Category: CategoryWallet},
	"wallet_identity_not_owned_by_user": {
		LegacyKey: "wallet_identity_not_owned_by_user",
		Code:      "WALLET_NOT_OWNED_BY_USER",
		Message:   "wallet identity not owned by user",
		Status:    http.StatusForbidden,
		Category:  CategoryWallet,
	},
	"wallet_link_challenge_error":         {LegacyKey: "wallet_link_challenge_error", Code: "WALLET_LINK_CHALLENGE_ERROR", Message: "wallet link challenge error", Status: http.StatusInternalServerError, Category: CategoryWallet},
	"wallet_link_challenge_user_mismatch": {LegacyKey: "wallet_link_challenge_user_mismatch", Code: "WALLET_LINK_CHALLENGE_USER_MISMATCH", Message: "wallet link challenge user mismatch", Status: http.StatusConflict, Category: CategoryWallet},
	"wallet_link_verify_error":            {LegacyKey: "wallet_link_verify_error", Code: "WALLET_LINK_VERIFY_ERROR", Message: "wallet link verification error", Status: http.StatusInternalServerError, Category: CategoryWallet},
	"wallet_primary_set_error":            {LegacyKey: "wallet_primary_set_error", Code: "WALLET_PRIMARY_SET_ERROR", Message: "wallet primary set error", Status: http.StatusInternalServerError, Category: CategoryWallet},
	"wallet_verify_error":                 {LegacyKey: "wallet_verify_error", Code: "WALLET_VERIFY_ERROR", Message: "wallet verification error", Status: http.StatusInternalServerError, Category: CategoryWallet},
}

func LookupAuthErrorSpec(legacyKey string) (AppErrorSpec, bool) {
	spec, ok := authErrorCatalog[strings.TrimSpace(legacyKey)]
	return spec, ok
}

func NormalizeLegacyAuthError(legacyKey string) AppErrorSpec {
	if spec, ok := LookupAuthErrorSpec(legacyKey); ok {
		return spec
	}

	trimmed := strings.TrimSpace(legacyKey)
	message := strings.ReplaceAll(trimmed, "_", " ")
	if message == "" {
		message = "unexpected error"
	}

	code := strings.ToUpper(strings.ReplaceAll(trimmed, "-", "_"))
	code = strings.ReplaceAll(code, " ", "_")
	if code == "" {
		code = "INTERNAL_ERROR"
	}

	return AppErrorSpec{
		LegacyKey: trimmed,
		Code:      code,
		Message:   message,
		Status:    http.StatusInternalServerError,
		Category:  CategoryInternal,
	}
}

func AppErrorFromLegacyAuthKey(legacyKey string, details map[string]any) *AppError {
	spec := NormalizeLegacyAuthError(legacyKey)
	err := New(spec.Code, spec.Message, spec.Status, spec.Category)
	if len(details) > 0 {
		err.Details = details
	}
	return err
}
