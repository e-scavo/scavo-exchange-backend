package errs

func AuthNotConfigured() *AppError {
	return AppErrorFromLegacyAuthKey("auth_not_configured", nil)
}

func AuthMissingBearerToken() *AppError {
	return AppErrorFromLegacyAuthKey("auth_missing_bearer_token", nil)
}

func AuthUnauthorized() *AppError {
	return AppErrorFromLegacyAuthKey("unauthorized", nil)
}

func AuthInvalidCredentials() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_credentials", nil)
}

func AuthServiceError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("auth_service_error", nil)
	err.Cause = cause
	return err
}

func AuthUserNotFound() *AppError {
	return AppErrorFromLegacyAuthKey("user_not_found", nil)
}

func BadRequestInvalidPayload() *AppError {
	return AppErrorFromLegacyAuthKey("bad_request", nil)
}

func InternalError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("internal_error", nil)
	err.Cause = cause
	return err
}

func Timeout() *AppError {
	return AppErrorFromLegacyAuthKey("timeout", nil)
}

func SettingsInvalidPayload() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_preferences", nil)
}

func WalletCannotDetach(details map[string]any) *AppError {
	return AppErrorFromLegacyAuthKey("wallet_detach_not_eligible", details)
}
