package errs

import (
	"net/http"
	"testing"
)

func TestNormalizeLegacyAuthError_KnownMappings(t *testing.T) {
	tests := []struct {
		name       string
		legacyKey  string
		wantCode   string
		wantStatus int
		wantCat    Category
	}{
		{name: "unauthorized", legacyKey: "unauthorized", wantCode: "AUTH_UNAUTHORIZED", wantStatus: http.StatusUnauthorized, wantCat: CategoryAuth},
		{name: "bad request", legacyKey: "bad_request", wantCode: "BAD_REQUEST", wantStatus: http.StatusBadRequest, wantCat: CategoryGeneric},
		{name: "settings invalid payload", legacyKey: "invalid_preferences", wantCode: "SETTINGS_INVALID_PAYLOAD", wantStatus: http.StatusBadRequest, wantCat: CategorySettings},
		{name: "wallet not owned", legacyKey: "wallet_identity_not_owned_by_user", wantCode: "WALLET_NOT_OWNED_BY_USER", wantStatus: http.StatusForbidden, wantCat: CategoryWallet},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := NormalizeLegacyAuthError(tt.legacyKey)
			if spec.Code != tt.wantCode {
				t.Fatalf("unexpected code: got=%q want=%q", spec.Code, tt.wantCode)
			}
			if spec.Status != tt.wantStatus {
				t.Fatalf("unexpected status: got=%d want=%d", spec.Status, tt.wantStatus)
			}
			if spec.Category != tt.wantCat {
				t.Fatalf("unexpected category: got=%q want=%q", spec.Category, tt.wantCat)
			}
		})
	}
}

func TestNormalizeLegacyAuthError_UnknownFallback(t *testing.T) {
	spec := NormalizeLegacyAuthError("some_new_error")
	if spec.Code != "SOME_NEW_ERROR" {
		t.Fatalf("unexpected code: %q", spec.Code)
	}
	if spec.Status != http.StatusInternalServerError {
		t.Fatalf("unexpected status: %d", spec.Status)
	}
	if spec.Category != CategoryInternal {
		t.Fatalf("unexpected category: %q", spec.Category)
	}
}

func TestAppError_ToResponseError_PreservesFields(t *testing.T) {
	appErr := WalletCannotDetach(map[string]any{"check": map[string]any{"eligible": false}})
	resp := appErr.ToResponseError()
	if resp.Code != "WALLET_CANNOT_DETACH" {
		t.Fatalf("unexpected code: %q", resp.Code)
	}
	if resp.Message != "wallet cannot be detached under current ownership rules" {
		t.Fatalf("unexpected message: %q", resp.Message)
	}
	if check, ok := resp.Details["check"].(map[string]any); !ok || check["eligible"] != false {
		t.Fatalf("unexpected details: %#v", resp.Details)
	}
}

func TestAppError_WithDetails_MergesWithoutMutation(t *testing.T) {
	base := AuthServiceError(nil).WithDetails(map[string]any{"phase": "0.8.4"})
	next := base.WithDetails(map[string]any{"step": "contract-tests"})

	if _, ok := base.Details["step"]; ok {
		t.Fatalf("base details mutated: %#v", base.Details)
	}
	if next.Details["phase"] != "0.8.4" || next.Details["step"] != "contract-tests" {
		t.Fatalf("unexpected merged details: %#v", next.Details)
	}
}

func TestAppError_NilToResponseErrorFallsBackToInternal(t *testing.T) {
	var appErr *AppError
	resp := appErr.ToResponseError()
	if resp.Code != "INTERNAL_ERROR" {
		t.Fatalf("unexpected code: %q", resp.Code)
	}
	if resp.Message != "internal server error" {
		t.Fatalf("unexpected message: %q", resp.Message)
	}
}

func TestAppError_WithContext_AddsSingleDetailWithoutMutation(t *testing.T) {
	base := AuthServiceError(nil).WithDetails(map[string]any{"phase": "0.14.3"})
	next := base.WithContext("component", "auth")

	if _, ok := base.Details["component"]; ok {
		t.Fatalf("base details mutated: %#v", base.Details)
	}
	if next.Details["phase"] != "0.14.3" || next.Details["component"] != "auth" {
		t.Fatalf("unexpected enriched details: %#v", next.Details)
	}
}

func TestAppError_WithRequestID_UsesStandardRequestIDKey(t *testing.T) {
	err := AuthServiceError(nil).WithRequestID("req-123")

	if err.Details["request_id"] != "req-123" {
		t.Fatalf("unexpected request_id detail: %#v", err.Details)
	}
}

func TestAppError_WithRequestID_EmptyRequestIDDoesNotMutate(t *testing.T) {
	base := AuthServiceError(nil).WithDetails(map[string]any{"phase": "0.14.3"})
	next := base.WithRequestID("")

	if next != base {
		t.Fatalf("empty request_id should keep original error instance")
	}
	if next.Details["phase"] != "0.14.3" {
		t.Fatalf("unexpected details: %#v", next.Details)
	}
}

func TestAppError_PublicDetails_ReturnsCopy(t *testing.T) {
	err := AuthServiceError(nil).WithDetails(map[string]any{"phase": "0.14.3"})
	details := err.PublicDetails()
	details["phase"] = "mutated"

	if err.Details["phase"] != "0.14.3" {
		t.Fatalf("public details mutation leaked into app error: %#v", err.Details)
	}
}

func TestAppError_ToResponseError_UsesPublicDetailsCopy(t *testing.T) {
	err := AuthServiceError(nil).WithRequestID("req-123")
	resp := err.ToResponseError()
	resp.Details["request_id"] = "mutated"

	if err.Details["request_id"] != "req-123" {
		t.Fatalf("response details mutation leaked into app error: %#v", err.Details)
	}
}
