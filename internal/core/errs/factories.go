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

func InvalidDisplayName() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_display_name", nil)
}

func DisplayNameTooLong() *AppError {
	return AppErrorFromLegacyAuthKey("display_name_too_long", nil)
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

func InvalidLimit() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_limit", nil)
}

func InvalidOffset() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_offset", nil)
}

func InvalidOrder() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_order", nil)
}

func InvalidOrderRequiresSort() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_order_requires_sort", nil)
}

func InvalidPrimary() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_primary", nil)
}

func InvalidSort() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_sort", nil)
}

func InvalidStatus() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_status", nil)
}

func WalletInvalidAddress() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_wallet_address", nil)
}

func WalletInvalidSignature() *AppError {
	return AppErrorFromLegacyAuthKey("invalid_wallet_signature", nil)
}

func WalletChallengeError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("wallet_challenge_error", nil)
	err.Cause = cause
	return err
}

func WalletChallengeNotFound() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_challenge_not_found", nil)
}

func WalletChallengeExpired() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_challenge_expired", nil)
}

func WalletChallengeUsed() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_challenge_used", nil)
}

func WalletChallengePurposeMismatch() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_challenge_purpose_mismatch", nil)
}

func WalletVerifyError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("wallet_verify_error", nil)
	err.Cause = cause
	return err
}

func WalletLinkChallengeError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("wallet_link_challenge_error", nil)
	err.Cause = cause
	return err
}

func WalletLinkChallengeUserMismatch() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_link_challenge_user_mismatch", nil)
}

func WalletLinkVerifyError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("wallet_link_verify_error", nil)
	err.Cause = cause
	return err
}

func WalletAccountMergeChallengeError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("wallet_account_merge_challenge_error", nil)
	err.Cause = cause
	return err
}

func WalletAccountMergeUserMismatch() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_account_merge_user_mismatch", nil)
}

func WalletAccountMergeSourceNotLinked() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_account_merge_source_not_linked", nil)
}

func WalletAccountMergeNotRequired() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_account_merge_not_required", nil)
}

func WalletAccountMergeVerifyError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("wallet_account_merge_verify_error", nil)
	err.Cause = cause
	return err
}

func WalletIdentityError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("wallet_identity_error", nil)
	err.Cause = cause
	return err
}

func WalletNotFound() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_identity_not_found", nil)
}

func WalletNotOwnedByUser() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_identity_not_owned_by_user", nil)
}

func WalletDetachCheckError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("wallet_detach_check_error", nil)
	err.Cause = cause
	return err
}

func WalletCannotDetach(details map[string]any) *AppError {
	return AppErrorFromLegacyAuthKey("wallet_detach_not_eligible", details)
}

func WalletDetachError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("wallet_detach_error", nil)
	err.Cause = cause
	return err
}

func WalletAlreadyLinked() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_identity_already_linked", nil)
}

func WalletAlreadyLinkedToUser() *AppError {
	return AppErrorFromLegacyAuthKey("wallet_identity_already_linked_to_user", nil)
}

func WalletPrimarySetError(cause error) *AppError {
	err := AppErrorFromLegacyAuthKey("wallet_primary_set_error", nil)
	err.Cause = cause
	return err
}
