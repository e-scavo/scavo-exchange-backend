package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

func mustTokenService(t *testing.T) *coreauth.TokenService {
	t.Helper()

	ts, err := coreauth.NewTokenService("dev_dev_dev_dev_dev_dev_dev_dev", "scavo-exchange-backend", time.Hour)
	if err != nil {
		t.Fatalf("NewTokenService error: %v", err)
	}

	return ts
}

func TestHTTPHandlers_Login_Success(t *testing.T) {
	h := HTTPHandlers{
		Tokens: mustTokenService(t),
		TTL:    time.Hour,
		Users:  usermod.NewService(nil),
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"test@example.com","password":"dev"}`))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload LoginResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if payload.AccessToken == "" {
		t.Fatal("expected access token")
	}

	if payload.UserID != "u_test_example_com" {
		t.Fatalf("unexpected user id: %q", payload.UserID)
	}
}

func TestHTTPHandlers_Login_InvalidBody(t *testing.T) {
	h := HTTPHandlers{
		Tokens: mustTokenService(t),
		TTL:    time.Hour,
		Users:  usermod.NewService(nil),
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":`))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func sessionClaims() *coreauth.Claims {
	now := time.Now().UTC()

	return &coreauth.Claims{
		UserID: "u_test_example_com",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "scavo-exchange-backend",
			Subject:   "u_test_example_com",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
}

func TestHTTPHandlers_Me_Success(t *testing.T) {
	ts := mustTokenService(t)

	h := HTTPHandlers{
		Tokens: ts,
		TTL:    time.Hour,
		Users:  usermod.NewService(nil),
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload MeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if payload.User == nil || payload.User.ID != "u_test_example_com" {
		t.Fatalf("unexpected user payload: %#v", payload.User)
	}
}

func TestHTTPHandlers_Me_MissingClaims(t *testing.T) {
	h := HTTPHandlers{
		Tokens: mustTokenService(t),
		TTL:    time.Hour,
		Users:  usermod.NewService(nil),
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func TestHTTPHandlers_Session_Success(t *testing.T) {
	h := HTTPHandlers{
		Tokens: mustTokenService(t),
		TTL:    time.Hour,
		Users:  usermod.NewService(nil),
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/session", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Session(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload SessionResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if payload.Session == nil {
		t.Fatal("expected session payload")
	}
	if !payload.Session.Authenticated {
		t.Fatal("expected authenticated session")
	}
	if payload.Session.UserID != "u_test_example_com" {
		t.Fatalf("unexpected session user id: %q", payload.Session.UserID)
	}
	if payload.Session.User == nil || payload.Session.User.Email != "test@example.com" {
		t.Fatalf("unexpected session user payload: %#v", payload.Session.User)
	}
}

func TestHTTPHandlers_Wallets_Success(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	address := testWalletAddress()

	identity, err := store.GetOrCreate(context.Background(), address)
	if err != nil {
		t.Fatalf("GetOrCreate error: %v", err)
	}

	attached, err := store.AttachUser(context.Background(), identity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("AttachUser error: %v", err)
	}

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/wallets", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Wallets(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(payload.Wallets) != 1 {
		t.Fatalf("expected 1 wallet, got %d", len(payload.Wallets))
	}
	if payload.Total != 1 {
		t.Fatalf("expected total=1, got %d", payload.Total)
	}
	if payload.Limit != 0 {
		t.Fatalf("expected default limit=0, got %d", payload.Limit)
	}
	if payload.Offset != 0 {
		t.Fatalf("expected default offset=0, got %d", payload.Offset)
	}
	if payload.Returned != 1 {
		t.Fatalf("expected returned=1, got %d", payload.Returned)
	}
	if payload.HasMore {
		t.Fatal("expected has_more=false by default")
	}
	if payload.Wallets[0].ID != identity.ID {
		t.Fatalf("unexpected wallet id: %q", payload.Wallets[0].ID)
	}
	if payload.Wallets[0].Address != address {
		t.Fatalf("unexpected wallet address: %q", payload.Wallets[0].Address)
	}
	if payload.Wallets[0].UserID != "u_test_example_com" {
		t.Fatalf("unexpected wallet user id: %q", payload.Wallets[0].UserID)
	}
	if !payload.Wallets[0].IsPrimary {
		t.Fatal("expected primary wallet")
	}
	if payload.Wallets[0].Status != "active" {
		t.Fatalf("unexpected wallet status: %q", payload.Wallets[0].Status)
	}
	if payload.Wallets[0].LinkedAt == nil {
		t.Fatal("expected linked_at")
	}
	if payload.Wallets[0].DetachedAt != nil {
		t.Fatalf("unexpected detached_at for active wallet: %#v", payload.Wallets[0].DetachedAt)
	}
	if attached.LinkedAt == nil || !payload.Wallets[0].LinkedAt.Equal(*attached.LinkedAt) {
		t.Fatalf("unexpected linked_at payload: %#v", payload.Wallets[0].LinkedAt)
	}
}

func mustSeedWalletIdentity(t *testing.T, store *InMemoryWalletIdentityStore, address, userID string, isPrimary bool, linkedAt time.Time) *WalletIdentity {
	t.Helper()

	identity, err := store.GetOrCreate(context.Background(), address)
	if err != nil {
		t.Fatalf("GetOrCreate error: %v", err)
	}

	attached, err := store.AttachUser(context.Background(), identity.ID, userID, isPrimary)
	if err != nil {
		t.Fatalf("AttachUser error: %v", err)
	}

	store.mu.Lock()
	if current := store.items[address]; current != nil {
		ts := linkedAt.UTC()
		current.LinkedAt = &ts
	}
	store.mu.Unlock()

	attached.LinkedAt = func() *time.Time {
		ts := linkedAt.UTC()
		return &ts
	}()

	return attached
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func TestHTTPHandlers_Wallets_ActionabilitySinglePrimary(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	now := time.Now().UTC()
	address := "0x7777777777777777777777777777777777777771"
	mustSeedWalletIdentity(t, store, address, "u_test_example_com", true, now.Add(-1*time.Hour))

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/wallets", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Wallets(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(payload.Wallets) != 1 {
		t.Fatalf("expected 1 wallet, got %d", len(payload.Wallets))
	}

	wallet := payload.Wallets[0]
	if wallet.CanSetPrimary {
		t.Fatal("expected can_set_primary=false for single primary wallet")
	}
	if wallet.CanDetach {
		t.Fatal("expected can_detach=false for single primary wallet")
	}
	if !containsString(wallet.DetachBlockReasons, WalletDetachReasonWalletIsPrimary) {
		t.Fatalf("expected detach block reason %q, got %#v", WalletDetachReasonWalletIsPrimary, wallet.DetachBlockReasons)
	}
	if !containsString(wallet.DetachBlockReasons, WalletDetachReasonUserWouldBeEmpty) {
		t.Fatalf("expected detach block reason %q, got %#v", WalletDetachReasonUserWouldBeEmpty, wallet.DetachBlockReasons)
	}
}

func TestHTTPHandlers_Wallets_ActionabilityTwoWalletInventory(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	now := time.Now().UTC()
	primaryAddress := "0x7777777777777777777777777777777777777772"
	secondaryAddress := "0x7777777777777777777777777777777777777773"
	mustSeedWalletIdentity(t, store, primaryAddress, "u_test_example_com", true, now.Add(-2*time.Hour))
	mustSeedWalletIdentity(t, store, secondaryAddress, "u_test_example_com", false, now.Add(-1*time.Hour))

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/wallets", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Wallets(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(payload.Wallets) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(payload.Wallets))
	}

	byAddress := map[string]*WalletReadModel{}
	for _, wallet := range payload.Wallets {
		byAddress[wallet.Address] = wallet
	}

	primary := byAddress[primaryAddress]
	if primary == nil {
		t.Fatalf("missing primary wallet %q", primaryAddress)
	}
	if primary.CanSetPrimary {
		t.Fatal("expected primary wallet can_set_primary=false")
	}
	if primary.CanDetach {
		t.Fatal("expected primary wallet can_detach=false")
	}
	if !containsString(primary.DetachBlockReasons, WalletDetachReasonWalletIsPrimary) {
		t.Fatalf("expected primary detach block reason %q, got %#v", WalletDetachReasonWalletIsPrimary, primary.DetachBlockReasons)
	}
	if containsString(primary.DetachBlockReasons, WalletDetachReasonUserWouldBeEmpty) {
		t.Fatalf("did not expect single-wallet block reason for primary in 2-wallet inventory: %#v", primary.DetachBlockReasons)
	}

	secondary := byAddress[secondaryAddress]
	if secondary == nil {
		t.Fatalf("missing secondary wallet %q", secondaryAddress)
	}
	if !secondary.CanSetPrimary {
		t.Fatal("expected secondary wallet can_set_primary=true")
	}
	if !secondary.CanDetach {
		t.Fatal("expected secondary wallet can_detach=true")
	}
	if len(secondary.DetachBlockReasons) != 0 {
		t.Fatalf("expected no detach block reasons for detachable secondary wallet, got %#v", secondary.DetachBlockReasons)
	}
}

func TestHTTPHandlers_WalletDetachCheck_ReadConsistencySinglePrimary(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	now := time.Now().UTC()
	address := "0x7777777777777777777777777777777777777781"
	mustSeedWalletIdentity(t, store, address, "u_test_example_com", true, now.Add(-1*time.Hour))

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	inventoryReq := httptest.NewRequest(http.MethodGet, "/auth/wallets", nil)
	inventoryReq = inventoryReq.WithContext(context.WithValue(inventoryReq.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	inventoryRec := httptest.NewRecorder()
	h.Wallets(inventoryRec, inventoryReq)

	if inventoryRec.Code != http.StatusOK {
		t.Fatalf("unexpected inventory status: %d body=%s", inventoryRec.Code, inventoryRec.Body.String())
	}

	var inventory WalletsResponse
	if err := json.Unmarshal(inventoryRec.Body.Bytes(), &inventory); err != nil {
		t.Fatalf("inventory decode error: %v", err)
	}
	if len(inventory.Wallets) != 1 {
		t.Fatalf("expected 1 wallet, got %d", len(inventory.Wallets))
	}

	wallet := inventory.Wallets[0]
	if wallet.Address != address {
		t.Fatalf("unexpected wallet address: %q", wallet.Address)
	}
	if wallet.CanDetach {
		t.Fatal("expected inventory can_detach=false for single primary wallet")
	}
	if !containsString(wallet.DetachBlockReasons, WalletDetachReasonWalletIsPrimary) {
		t.Fatalf("expected inventory block reason %q, got %#v", WalletDetachReasonWalletIsPrimary, wallet.DetachBlockReasons)
	}
	if !containsString(wallet.DetachBlockReasons, WalletDetachReasonUserWouldBeEmpty) {
		t.Fatalf("expected inventory block reason %q, got %#v", WalletDetachReasonUserWouldBeEmpty, wallet.DetachBlockReasons)
	}

	checkReq := httptest.NewRequest(http.MethodPost, "/auth/wallets/detach/check", strings.NewReader(`{"wallet_address":"`+address+`"}`))
	checkReq = checkReq.WithContext(context.WithValue(checkReq.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	checkRec := httptest.NewRecorder()
	h.WalletDetachCheck(checkRec, checkReq)

	if checkRec.Code != http.StatusOK {
		t.Fatalf("unexpected detach-check status: %d body=%s", checkRec.Code, checkRec.Body.String())
	}

	var check WalletDetachCheckResponse
	if err := json.Unmarshal(checkRec.Body.Bytes(), &check); err != nil {
		t.Fatalf("detach-check decode error: %v", err)
	}
	if check.Eligible {
		t.Fatalf("expected eligible=false, got reasons=%v", check.Reasons)
	}
	if !check.IsPrimary {
		t.Fatal("expected detach check is_primary=true")
	}
	if check.OwnedWalletCount != 1 {
		t.Fatalf("expected owned wallet count 1, got %d", check.OwnedWalletCount)
	}
	if !containsString(check.Reasons, WalletDetachReasonWalletIsPrimary) {
		t.Fatalf("expected detach check reason %q, got %#v", WalletDetachReasonWalletIsPrimary, check.Reasons)
	}
	if !containsString(check.Reasons, WalletDetachReasonUserWouldBeEmpty) {
		t.Fatalf("expected detach check reason %q, got %#v", WalletDetachReasonUserWouldBeEmpty, check.Reasons)
	}
}

func TestHTTPHandlers_WalletDetachCheck_ReadConsistencyTwoWalletInventory(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	now := time.Now().UTC()
	primaryAddress := "0x7777777777777777777777777777777777777782"
	secondaryAddress := "0x7777777777777777777777777777777777777783"
	mustSeedWalletIdentity(t, store, primaryAddress, "u_test_example_com", true, now.Add(-2*time.Hour))
	mustSeedWalletIdentity(t, store, secondaryAddress, "u_test_example_com", false, now.Add(-1*time.Hour))

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	inventoryReq := httptest.NewRequest(http.MethodGet, "/auth/wallets", nil)
	inventoryReq = inventoryReq.WithContext(context.WithValue(inventoryReq.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	inventoryRec := httptest.NewRecorder()
	h.Wallets(inventoryRec, inventoryReq)

	if inventoryRec.Code != http.StatusOK {
		t.Fatalf("unexpected inventory status: %d body=%s", inventoryRec.Code, inventoryRec.Body.String())
	}

	var inventory WalletsResponse
	if err := json.Unmarshal(inventoryRec.Body.Bytes(), &inventory); err != nil {
		t.Fatalf("inventory decode error: %v", err)
	}
	if len(inventory.Wallets) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(inventory.Wallets))
	}

	byAddress := map[string]*WalletReadModel{}
	for _, wallet := range inventory.Wallets {
		byAddress[wallet.Address] = wallet
	}

	primaryWallet := byAddress[primaryAddress]
	if primaryWallet == nil {
		t.Fatalf("missing primary inventory wallet %q", primaryAddress)
	}
	if primaryWallet.CanDetach {
		t.Fatal("expected primary inventory can_detach=false")
	}
	if !containsString(primaryWallet.DetachBlockReasons, WalletDetachReasonWalletIsPrimary) {
		t.Fatalf("expected primary inventory block reason %q, got %#v", WalletDetachReasonWalletIsPrimary, primaryWallet.DetachBlockReasons)
	}
	if containsString(primaryWallet.DetachBlockReasons, WalletDetachReasonUserWouldBeEmpty) {
		t.Fatalf("did not expect single-wallet reason for primary in two-wallet inventory: %#v", primaryWallet.DetachBlockReasons)
	}

	secondaryWallet := byAddress[secondaryAddress]
	if secondaryWallet == nil {
		t.Fatalf("missing secondary inventory wallet %q", secondaryAddress)
	}
	if !secondaryWallet.CanDetach {
		t.Fatal("expected secondary inventory can_detach=true")
	}
	if len(secondaryWallet.DetachBlockReasons) != 0 {
		t.Fatalf("expected no secondary detach block reasons, got %#v", secondaryWallet.DetachBlockReasons)
	}

	for _, tc := range []struct {
		name              string
		address           string
		expectedEligible  bool
		expectedIsPrimary bool
		expectedReason    string
		unexpectedReason  string
	}{
		{name: "primary", address: primaryAddress, expectedEligible: false, expectedIsPrimary: true, expectedReason: WalletDetachReasonWalletIsPrimary, unexpectedReason: WalletDetachReasonUserWouldBeEmpty},
		{name: "secondary", address: secondaryAddress, expectedEligible: true, expectedIsPrimary: false, expectedReason: "", unexpectedReason: WalletDetachReasonWalletIsPrimary},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checkReq := httptest.NewRequest(http.MethodPost, "/auth/wallets/detach/check", strings.NewReader(`{"wallet_address":"`+tc.address+`"}`))
			checkReq = checkReq.WithContext(context.WithValue(checkReq.Context(), coreauth.ClaimsContextKey, sessionClaims()))
			checkRec := httptest.NewRecorder()
			h.WalletDetachCheck(checkRec, checkReq)

			if checkRec.Code != http.StatusOK {
				t.Fatalf("unexpected detach-check status: %d body=%s", checkRec.Code, checkRec.Body.String())
			}

			var check WalletDetachCheckResponse
			if err := json.Unmarshal(checkRec.Body.Bytes(), &check); err != nil {
				t.Fatalf("detach-check decode error: %v", err)
			}
			if check.Eligible != tc.expectedEligible {
				t.Fatalf("unexpected eligible=%v want %v reasons=%v", check.Eligible, tc.expectedEligible, check.Reasons)
			}
			if check.IsPrimary != tc.expectedIsPrimary {
				t.Fatalf("unexpected is_primary=%v want %v", check.IsPrimary, tc.expectedIsPrimary)
			}
			if check.OwnedWalletCount != 2 {
				t.Fatalf("expected owned wallet count 2, got %d", check.OwnedWalletCount)
			}
			if tc.expectedReason != "" && !containsString(check.Reasons, tc.expectedReason) {
				t.Fatalf("expected detach check reason %q, got %#v", tc.expectedReason, check.Reasons)
			}
			if tc.unexpectedReason != "" && containsString(check.Reasons, tc.unexpectedReason) {
				t.Fatalf("did not expect detach check reason %q, got %#v", tc.unexpectedReason, check.Reasons)
			}
		})
	}
}

func TestHTTPHandlers_Wallets_FilterPrimary(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	primaryAddress := "0x1111111111111111111111111111111111111111"
	secondaryAddress := "0x2222222222222222222222222222222222222222"
	now := time.Now().UTC()

	mustSeedWalletIdentity(t, store, primaryAddress, "u_test_example_com", true, now.Add(-2*time.Hour))
	mustSeedWalletIdentity(t, store, secondaryAddress, "u_test_example_com", false, now.Add(-1*time.Hour))

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	for _, tc := range []struct {
		name            string
		query           string
		expectedCount   int
		expectedAddr    string
		expectedPrimary bool
	}{
		{name: "primary_true", query: "/auth/wallets?primary=true", expectedCount: 1, expectedAddr: primaryAddress, expectedPrimary: true},
		{name: "primary_false", query: "/auth/wallets?primary=false", expectedCount: 1, expectedAddr: secondaryAddress, expectedPrimary: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
			rec := httptest.NewRecorder()

			h.Wallets(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
			}

			var payload WalletsResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
				t.Fatalf("decode error: %v", err)
			}

			if len(payload.Wallets) != tc.expectedCount {
				t.Fatalf("expected %d wallets, got %d", tc.expectedCount, len(payload.Wallets))
			}
			if payload.Wallets[0].Address != tc.expectedAddr {
				t.Fatalf("unexpected wallet address: %q", payload.Wallets[0].Address)
			}
			if payload.Wallets[0].IsPrimary != tc.expectedPrimary {
				t.Fatalf("unexpected primary flag: %v", payload.Wallets[0].IsPrimary)
			}
		})
	}
}

func TestHTTPHandlers_Wallets_FilterStatus(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	address := "0x3333333333333333333333333333333333333333"
	mustSeedWalletIdentity(t, store, address, "u_test_example_com", true, time.Now().UTC())

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	for _, tc := range []struct {
		name          string
		query         string
		expectedCount int
	}{
		{name: "active", query: "/auth/wallets?status=active", expectedCount: 1},
		{name: "detached", query: "/auth/wallets?status=detached", expectedCount: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
			rec := httptest.NewRecorder()

			h.Wallets(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
			}

			var payload WalletsResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
				t.Fatalf("decode error: %v", err)
			}

			if len(payload.Wallets) != tc.expectedCount {
				t.Fatalf("expected %d wallets, got %d", tc.expectedCount, len(payload.Wallets))
			}
		})
	}
}

func TestHTTPHandlers_Wallets_SortLinkedAt(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	earlierAddress := "0x4444444444444444444444444444444444444444"
	laterAddress := "0x5555555555555555555555555555555555555555"
	now := time.Now().UTC()
	earlier := now.Add(-2 * time.Hour)
	later := now.Add(-30 * time.Minute)

	mustSeedWalletIdentity(t, store, earlierAddress, "u_test_example_com", false, earlier)
	mustSeedWalletIdentity(t, store, laterAddress, "u_test_example_com", true, later)

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/wallets?sort=linked_at&order=desc", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Wallets(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(payload.Wallets) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(payload.Wallets))
	}
	if payload.Wallets[0].Address != laterAddress {
		t.Fatalf("unexpected first wallet: %q", payload.Wallets[0].Address)
	}
	if payload.Wallets[1].Address != earlierAddress {
		t.Fatalf("unexpected second wallet: %q", payload.Wallets[1].Address)
	}
}

func TestHTTPHandlers_Wallets_SortLinkedAt_DefaultOrderAsc(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	earlierAddress := "0x4545454545454545454545454545454545454545"
	laterAddress := "0x5656565656565656565656565656565656565656"
	now := time.Now().UTC()
	earlier := now.Add(-2 * time.Hour)
	later := now.Add(-30 * time.Minute)

	mustSeedWalletIdentity(t, store, earlierAddress, "u_test_example_com", false, earlier)
	mustSeedWalletIdentity(t, store, laterAddress, "u_test_example_com", true, later)

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/wallets?sort=linked_at", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Wallets(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(payload.Wallets) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(payload.Wallets))
	}
	if payload.Wallets[0].Address != earlierAddress {
		t.Fatalf("unexpected first wallet: %q", payload.Wallets[0].Address)
	}
	if payload.Wallets[1].Address != laterAddress {
		t.Fatalf("unexpected second wallet: %q", payload.Wallets[1].Address)
	}
}

func TestHTTPHandlers_Wallets_OffsetWithoutLimitKeepsUnboundedWindow(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	now := time.Now().UTC()

	mustSeedWalletIdentity(t, store, "0x6868686868686868686868686868686868686861", "u_test_example_com", true, now.Add(-3*time.Hour))
	mustSeedWalletIdentity(t, store, "0x6868686868686868686868686868686868686862", "u_test_example_com", false, now.Add(-2*time.Hour))
	mustSeedWalletIdentity(t, store, "0x6868686868686868686868686868686868686863", "u_test_example_com", false, now.Add(-1*time.Hour))

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/wallets?offset=1", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Wallets(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if payload.Total != 3 {
		t.Fatalf("expected total=3, got %d", payload.Total)
	}
	if payload.Limit != 0 {
		t.Fatalf("expected limit=0, got %d", payload.Limit)
	}
	if payload.Offset != 1 {
		t.Fatalf("expected offset=1, got %d", payload.Offset)
	}
	if payload.Returned != 2 {
		t.Fatalf("expected returned=2, got %d", payload.Returned)
	}
	if payload.HasMore {
		t.Fatal("expected has_more=false without explicit limit")
	}
	if payload.NextOffset != nil {
		t.Fatalf("expected next_offset=nil without explicit limit, got %#v", payload.NextOffset)
	}
	if payload.PreviousOffset != nil {
		t.Fatalf("expected previous_offset=nil without explicit limit, got %#v", payload.PreviousOffset)
	}
}

func TestHTTPHandlers_Wallets_InvalidQueryParams(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	mustSeedWalletIdentity(t, store, "0x6666666666666666666666666666666666666666", "u_test_example_com", true, time.Now().UTC())

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	for _, tc := range []struct {
		name          string
		query         string
		expectedError string
	}{
		{name: "invalid_status", query: "/auth/wallets?status=whatever", expectedError: "invalid_status"},
		{name: "invalid_primary", query: "/auth/wallets?primary=maybe", expectedError: "invalid_primary"},
		{name: "invalid_sort", query: "/auth/wallets?sort=address", expectedError: "invalid_sort"},
		{name: "invalid_order", query: "/auth/wallets?sort=linked_at&order=sideways", expectedError: "invalid_order"},
		{name: "order_without_sort", query: "/auth/wallets?order=desc", expectedError: "invalid_order_requires_sort"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
			rec := httptest.NewRecorder()

			h.Wallets(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
			}

			var payload map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
				t.Fatalf("decode error: %v", err)
			}

			if payload["error"] != tc.expectedError {
				t.Fatalf("unexpected error payload: %#v", payload)
			}
		})
	}
}

func TestHTTPHandlers_Wallets_Pagination(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	now := time.Now().UTC()

	mustSeedWalletIdentity(t, store, "0x7777777777777777777777777777777777777771", "u_test_example_com", true, now.Add(-3*time.Hour))
	mustSeedWalletIdentity(t, store, "0x7777777777777777777777777777777777777772", "u_test_example_com", false, now.Add(-2*time.Hour))
	mustSeedWalletIdentity(t, store, "0x7777777777777777777777777777777777777773", "u_test_example_com", false, now.Add(-1*time.Hour))

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	for _, tc := range []struct {
		name                   string
		query                  string
		expectedAddrs          []string
		expectedTotal          int
		expectedLimit          int
		expectedOffset         int
		expectedReturned       int
		expectedHasMore        bool
		expectedNextOffset     *int
		expectedPreviousOffset *int
	}{
		{name: "limit_only", query: "/auth/wallets?limit=2", expectedAddrs: []string{"0x7777777777777777777777777777777777777771", "0x7777777777777777777777777777777777777772"}, expectedTotal: 3, expectedLimit: 2, expectedOffset: 0, expectedReturned: 2, expectedHasMore: true, expectedNextOffset: intPtr(2), expectedPreviousOffset: nil},
		{name: "offset_only", query: "/auth/wallets?offset=1", expectedAddrs: []string{"0x7777777777777777777777777777777777777772", "0x7777777777777777777777777777777777777773"}, expectedTotal: 3, expectedLimit: 0, expectedOffset: 1, expectedReturned: 2, expectedHasMore: false, expectedNextOffset: nil, expectedPreviousOffset: nil},
		{name: "limit_and_offset", query: "/auth/wallets?limit=1&offset=1", expectedAddrs: []string{"0x7777777777777777777777777777777777777772"}, expectedTotal: 3, expectedLimit: 1, expectedOffset: 1, expectedReturned: 1, expectedHasMore: true, expectedNextOffset: intPtr(2), expectedPreviousOffset: intPtr(0)},
		{name: "window_empty", query: "/auth/wallets?limit=2&offset=10", expectedAddrs: []string{}, expectedTotal: 3, expectedLimit: 2, expectedOffset: 10, expectedReturned: 0, expectedHasMore: false, expectedNextOffset: nil, expectedPreviousOffset: intPtr(8)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
			rec := httptest.NewRecorder()

			h.Wallets(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
			}

			var payload WalletsResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
				t.Fatalf("decode error: %v", err)
			}

			if payload.Total != tc.expectedTotal {
				t.Fatalf("expected total=%d, got %d", tc.expectedTotal, payload.Total)
			}
			if payload.Limit != tc.expectedLimit {
				t.Fatalf("expected limit=%d, got %d", tc.expectedLimit, payload.Limit)
			}
			if payload.Offset != tc.expectedOffset {
				t.Fatalf("expected offset=%d, got %d", tc.expectedOffset, payload.Offset)
			}
			if payload.Returned != tc.expectedReturned {
				t.Fatalf("expected returned=%d, got %d", tc.expectedReturned, payload.Returned)
			}
			if payload.HasMore != tc.expectedHasMore {
				t.Fatalf("expected has_more=%v, got %v", tc.expectedHasMore, payload.HasMore)
			}
			if !equalOptionalInt(payload.NextOffset, tc.expectedNextOffset) {
				t.Fatalf("unexpected next_offset: got=%#v want=%#v", payload.NextOffset, tc.expectedNextOffset)
			}
			if !equalOptionalInt(payload.PreviousOffset, tc.expectedPreviousOffset) {
				t.Fatalf("unexpected previous_offset: got=%#v want=%#v", payload.PreviousOffset, tc.expectedPreviousOffset)
			}
			if len(payload.Wallets) != len(tc.expectedAddrs) {
				t.Fatalf("expected %d wallets, got %d", len(tc.expectedAddrs), len(payload.Wallets))
			}
			for i, addr := range tc.expectedAddrs {
				if payload.Wallets[i].Address != addr {
					t.Fatalf("unexpected wallet address at %d: %q", i, payload.Wallets[i].Address)
				}
			}
		})
	}
}

func TestHTTPHandlers_Wallets_NavigationMetadataWithFiltersAndSorting(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	now := time.Now().UTC()

	mustSeedWalletIdentity(t, store, "0x9999999999999999999999999999999999999991", "u_test_example_com", true, now.Add(-4*time.Hour))
	mustSeedWalletIdentity(t, store, "0x9999999999999999999999999999999999999992", "u_test_example_com", false, now.Add(-3*time.Hour))
	mustSeedWalletIdentity(t, store, "0x9999999999999999999999999999999999999993", "u_test_example_com", false, now.Add(-2*time.Hour))
	mustSeedWalletIdentity(t, store, "0x9999999999999999999999999999999999999994", "u_test_example_com", false, now.Add(-1*time.Hour))

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/wallets?primary=false&sort=linked_at&order=desc&limit=2&offset=1", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Wallets(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if payload.Total != 3 {
		t.Fatalf("expected total=3, got %d", payload.Total)
	}
	if payload.Limit != 2 {
		t.Fatalf("expected limit=2, got %d", payload.Limit)
	}
	if payload.Offset != 1 {
		t.Fatalf("expected offset=1, got %d", payload.Offset)
	}
	if payload.Returned != 2 {
		t.Fatalf("expected returned=2, got %d", payload.Returned)
	}
	if payload.HasMore {
		t.Fatal("expected has_more=false on final filtered window")
	}
	if payload.NextOffset != nil {
		t.Fatalf("expected next_offset=nil on final filtered window, got %#v", payload.NextOffset)
	}
	if payload.PreviousOffset == nil || *payload.PreviousOffset != 0 {
		t.Fatalf("expected previous_offset=0 on final filtered window, got %#v", payload.PreviousOffset)
	}
	if len(payload.Wallets) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(payload.Wallets))
	}
	if payload.Wallets[0].Address != "0x9999999999999999999999999999999999999993" {
		t.Fatalf("unexpected first wallet: %q", payload.Wallets[0].Address)
	}
	if payload.Wallets[1].Address != "0x9999999999999999999999999999999999999992" {
		t.Fatalf("unexpected second wallet: %q", payload.Wallets[1].Address)
	}
}

func TestHTTPHandlers_Wallets_InvalidPaginationParams(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	mustSeedWalletIdentity(t, store, "0x8888888888888888888888888888888888888888", "u_test_example_com", true, time.Now().UTC())

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	for _, tc := range []struct {
		name          string
		query         string
		expectedError string
	}{
		{name: "invalid_limit_text", query: "/auth/wallets?limit=abc", expectedError: "invalid_limit"},
		{name: "invalid_limit_zero", query: "/auth/wallets?limit=0", expectedError: "invalid_limit"},
		{name: "invalid_limit_negative", query: "/auth/wallets?limit=-1", expectedError: "invalid_limit"},
		{name: "invalid_offset_text", query: "/auth/wallets?offset=abc", expectedError: "invalid_offset"},
		{name: "invalid_offset_negative", query: "/auth/wallets?offset=-1", expectedError: "invalid_offset"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
			rec := httptest.NewRecorder()

			h.Wallets(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
			}

			var payload map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
				t.Fatalf("decode error: %v", err)
			}

			if payload["error"] != tc.expectedError {
				t.Fatalf("unexpected error payload: %#v", payload)
			}
		})
	}
}

func intPtr(v int) *int {
	return &v
}

func equalOptionalInt(left, right *int) bool {
	switch {
	case left == nil && right == nil:
		return true
	case left == nil || right == nil:
		return false
	default:
		return *left == *right
	}
}

func TestHTTPHandlers_WalletLinkChallenge_Success(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()

	h := HTTPHandlers{
		Tokens:       mustTokenService(t),
		TTL:          time.Hour,
		Users:        usermod.NewService(nil),
		Challenges:   store,
		ChallengeTTL: 5 * time.Minute,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/wallets/link/challenge", strings.NewReader(`{"address":"0x1111111111111111111111111111111111111111","chain":"scavium"}`))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletLinkChallenge(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletLinkChallengeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Challenge == nil {
		t.Fatal("expected challenge payload")
	}
	if payload.Challenge.Purpose != WalletChallengePurposeLinkWallet {
		t.Fatalf("unexpected challenge purpose: %q", payload.Challenge.Purpose)
	}
	if payload.Challenge.RequestedByUserID != "u_test_example_com" {
		t.Fatalf("unexpected requested user id: %q", payload.Challenge.RequestedByUserID)
	}
}

func TestHTTPHandlers_WalletLinkVerify_Success(t *testing.T) {
	challengeStore := NewInMemoryWalletChallengeStore()
	identityStore := NewInMemoryWalletIdentityStore()

	primaryAddress, _ := signWalletMessageForScalar(t, "bootstrap", "1")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "bootstrap2", "2")
	challengeSvc := NewWalletChallengeService(challengeStore, "https://api.scavo.exchange", 5*time.Minute)
	challenge, err := challengeSvc.CreateWithOptions(context.Background(), secondaryAddress, "scavium", WalletChallengeOptions{
		Purpose:           WalletChallengePurposeLinkWallet,
		RequestedByUserID: "u_test_example_com",
	})
	if err != nil {
		t.Fatalf("CreateWithOptions error: %v", err)
	}
	_, signature := signWalletMessageForScalar(t, challenge.Message, "2")

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		Challenges:       challengeStore,
		WalletIdentities: identityStore,
		ChallengeTTL:     5 * time.Minute,
		PublicBaseURL:    "https://api.scavo.exchange",
	}

	body := `{"challenge_id":"` + challenge.ID + `","address":"` + secondaryAddress + `","signature":"` + signature + `"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/wallets/link/verify", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletLinkVerify(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletLinkVerifyResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.LinkedWallet == nil {
		t.Fatal("expected linked wallet")
	}
	if payload.LinkedWallet.UserID != "u_test_example_com" {
		t.Fatalf("unexpected linked wallet user id: %q", payload.LinkedWallet.UserID)
	}
	if payload.LinkedWallet.IsPrimary {
		t.Fatal("expected linked wallet to remain secondary")
	}
	if len(payload.Wallets) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(payload.Wallets))
	}
}

func TestHTTPHandlers_WalletLinkVerify_AllowsDetachedWalletReattachment(t *testing.T) {
	challengeStore := NewInMemoryWalletChallengeStore()
	identityStore := NewInMemoryWalletIdentityStore()
	detachSvc := NewWalletDetachService(identityStore)

	primaryAddress, _ := signWalletMessageForScalar(t, "handler-reattach-primary", "46")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "handler-reattach-secondary", "47")
	secondaryIdentity, err := identityStore.GetOrCreate(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate secondary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), secondaryIdentity.ID, "u_test_example_com", false)
	if err != nil {
		t.Fatalf("AttachUser secondary error: %v", err)
	}

	_, err = detachSvc.Execute(context.Background(), "u_test_example_com", secondaryAddress)
	if err != nil {
		t.Fatalf("Execute detach error: %v", err)
	}

	challengeSvc := NewWalletChallengeService(challengeStore, "https://api.scavo.exchange", 5*time.Minute)
	challenge, err := challengeSvc.CreateWithOptions(context.Background(), secondaryAddress, "scavium", WalletChallengeOptions{
		Purpose:           WalletChallengePurposeLinkWallet,
		RequestedByUserID: "u_test_example_com",
	})
	if err != nil {
		t.Fatalf("CreateWithOptions error: %v", err)
	}
	_, signature := signWalletMessageForScalar(t, challenge.Message, "47")

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		Challenges:       challengeStore,
		WalletIdentities: identityStore,
		ChallengeTTL:     5 * time.Minute,
		PublicBaseURL:    "https://api.scavo.exchange",
	}

	body := `{"challenge_id":"` + challenge.ID + `","address":"` + secondaryAddress + `","signature":"` + signature + `"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/wallets/link/verify", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletLinkVerify(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletLinkVerifyResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.LinkedWallet == nil {
		t.Fatal("expected linked wallet")
	}
	if payload.LinkedWallet.UserID != "u_test_example_com" {
		t.Fatalf("unexpected relinked wallet user id: %q", payload.LinkedWallet.UserID)
	}
	if payload.LinkedWallet.IsPrimary {
		t.Fatal("expected relinked wallet to remain secondary")
	}
	if payload.LinkedWallet.LinkedAt == nil {
		t.Fatal("expected relinked wallet linked_at")
	}
	if len(payload.Wallets) != 2 {
		t.Fatalf("expected 2 wallets after reattachment, got %d", len(payload.Wallets))
	}
}

func TestHTTPHandlers_WalletAccountMergeChallenge_Success(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()

	h := HTTPHandlers{
		Tokens:       mustTokenService(t),
		TTL:          time.Hour,
		Users:        usermod.NewService(nil),
		Challenges:   store,
		ChallengeTTL: 5 * time.Minute,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/account/merge/wallet/challenge", strings.NewReader(`{"address":"0x1111111111111111111111111111111111111111","chain":"scavium"}`))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletAccountMergeChallenge(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletAccountMergeChallengeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Challenge == nil {
		t.Fatal("expected challenge payload")
	}
	if payload.Challenge.Purpose != WalletChallengePurposeAccountMerge {
		t.Fatalf("unexpected challenge purpose: %q", payload.Challenge.Purpose)
	}
}

func TestHTTPHandlers_WalletAccountMergeVerify_Success(t *testing.T) {
	challengeStore := NewInMemoryWalletChallengeStore()
	identityStore := NewInMemoryWalletIdentityStore()

	targetPrimaryAddress, _ := signWalletMessageForScalar(t, "handler-target-primary", "1")
	targetPrimaryIdentity, err := identityStore.GetOrCreate(context.Background(), targetPrimaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate target primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), targetPrimaryIdentity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("AttachUser target primary error: %v", err)
	}

	sourcePrimaryAddress, _ := signWalletMessageForScalar(t, "handler-source-primary", "2")
	sourcePrimaryIdentity, err := identityStore.GetOrCreate(context.Background(), sourcePrimaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate source primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), sourcePrimaryIdentity.ID, "u_wallet_source", true)
	if err != nil {
		t.Fatalf("AttachUser source primary error: %v", err)
	}

	challengeSvc := NewWalletChallengeService(challengeStore, "https://api.scavo.exchange", 5*time.Minute)
	challenge, err := challengeSvc.CreateWithOptions(context.Background(), sourcePrimaryAddress, "scavium", WalletChallengeOptions{
		Purpose:           WalletChallengePurposeAccountMerge,
		RequestedByUserID: "u_test_example_com",
	})
	if err != nil {
		t.Fatalf("CreateWithOptions error: %v", err)
	}
	_, signature := signWalletMessageForScalar(t, challenge.Message, "2")

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		Challenges:       challengeStore,
		WalletIdentities: identityStore,
		ChallengeTTL:     5 * time.Minute,
		PublicBaseURL:    "https://api.scavo.exchange",
	}

	body := `{"challenge_id":"` + challenge.ID + `","address":"` + sourcePrimaryAddress + `","signature":"` + signature + `"}`
	req := httptest.NewRequest(http.MethodPost, "/auth/account/merge/wallet/verify", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletAccountMergeVerify(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletAccountMergeVerifyResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.MergedWallet == nil {
		t.Fatal("expected merged wallet")
	}
	if payload.MergedWallet.UserID != "u_test_example_com" {
		t.Fatalf("unexpected merged wallet user id: %q", payload.MergedWallet.UserID)
	}
	if payload.SourceUserID != "u_wallet_source" {
		t.Fatalf("unexpected source user id: %q", payload.SourceUserID)
	}
	if payload.TargetUserID != "u_test_example_com" {
		t.Fatalf("unexpected target user id: %q", payload.TargetUserID)
	}
	if len(payload.Wallets) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(payload.Wallets))
	}
}

func TestHTTPHandlers_WalletSetPrimary_Success(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()

	primaryAddress, _ := signWalletMessageForScalar(t, "handler-primary", "21")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "handler-secondary", "22")
	secondaryIdentity, err := identityStore.GetOrCreate(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate secondary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), secondaryIdentity.ID, "u_test_example_com", false)
	if err != nil {
		t.Fatalf("AttachUser secondary error: %v", err)
	}

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: identityStore,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/wallets/primary", strings.NewReader(`{"wallet_address":"`+secondaryAddress+`"}`))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletSetPrimary(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletPrimarySetResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.PrimaryWallet == nil || payload.PrimaryWallet.Address != secondaryAddress {
		t.Fatalf("unexpected primary wallet payload: %#v", payload.PrimaryWallet)
	}
	if len(payload.Wallets) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(payload.Wallets))
	}
	if payload.Wallets[0].Address != secondaryAddress || !payload.Wallets[0].IsPrimary {
		t.Fatal("expected switched wallet to be first and primary")
	}
}

func TestHTTPHandlers_WalletSetPrimary_RejectsWalletNotOwnedByUser(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()

	otherAddress, _ := signWalletMessageForScalar(t, "handler-other", "23")
	otherIdentity, err := identityStore.GetOrCreate(context.Background(), otherAddress)
	if err != nil {
		t.Fatalf("GetOrCreate other error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), otherIdentity.ID, "u_other", true)
	if err != nil {
		t.Fatalf("AttachUser other error: %v", err)
	}

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: identityStore,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/wallets/primary", strings.NewReader(`{"wallet_address":"`+otherAddress+`"}`))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletSetPrimary(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestHTTPHandlers_WalletDetachCheck_Success(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()

	primaryAddress, _ := signWalletMessageForScalar(t, "handler-detach-primary", "33")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "handler-detach-secondary", "34")
	secondaryIdentity, err := identityStore.GetOrCreate(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate secondary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), secondaryIdentity.ID, "u_test_example_com", false)
	if err != nil {
		t.Fatalf("AttachUser secondary error: %v", err)
	}

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: identityStore,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/wallets/detach/check", strings.NewReader(`{"wallet_address":"`+secondaryAddress+`"}`))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletDetachCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletDetachCheckResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !payload.Eligible {
		t.Fatalf("expected eligible payload, got reasons=%v", payload.Reasons)
	}
	if payload.OwnedWalletCount != 2 {
		t.Fatalf("expected owned wallet count 2, got %d", payload.OwnedWalletCount)
	}
}

func TestHTTPHandlers_WalletDetachCheck_RejectsWalletNotOwnedByUser(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()

	otherAddress, _ := signWalletMessageForScalar(t, "handler-detach-other", "35")
	otherIdentity, err := identityStore.GetOrCreate(context.Background(), otherAddress)
	if err != nil {
		t.Fatalf("GetOrCreate other error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), otherIdentity.ID, "u_other", true)
	if err != nil {
		t.Fatalf("AttachUser other error: %v", err)
	}

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: identityStore,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/wallets/detach/check", strings.NewReader(`{"wallet_address":"`+otherAddress+`"}`))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletDetachCheck(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}
}

type walletDetachConflictPayload struct {
	Error string                     `json:"error"`
	Check *WalletDetachCheckResponse `json:"check,omitempty"`
}

func TestHTTPHandlers_WalletDetach_Success(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()

	primaryAddress, _ := signWalletMessageForScalar(t, "handler-detach-exec-primary", "39")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "handler-detach-exec-secondary", "40")
	secondaryIdentity, err := identityStore.GetOrCreate(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate secondary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), secondaryIdentity.ID, "u_test_example_com", false)
	if err != nil {
		t.Fatalf("AttachUser secondary error: %v", err)
	}

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: identityStore,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/wallets/detach", strings.NewReader(`{"wallet_address":"`+secondaryAddress+`"}`))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletDetach(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletDetachExecuteResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.DetachedWallet == nil || payload.DetachedWallet.Address != secondaryAddress {
		t.Fatalf("unexpected detached wallet: %#v", payload.DetachedWallet)
	}
	if payload.DetachedWallet.UserID != "" {
		t.Fatalf("expected detached wallet user to be empty, got %q", payload.DetachedWallet.UserID)
	}
	if payload.DetachedWallet.DetachedAt == nil {
		t.Fatal("expected detached wallet detached_at metadata")
	}
	if len(payload.Wallets) != 1 {
		t.Fatalf("expected 1 remaining wallet, got %d", len(payload.Wallets))
	}
	if payload.Wallets[0].Address != primaryAddress || !payload.Wallets[0].IsPrimary {
		t.Fatal("expected original primary wallet to remain attached and primary")
	}
	if payload.Check == nil || !payload.Check.Eligible {
		t.Fatalf("expected successful eligibility snapshot, got %#v", payload.Check)
	}
}

func TestHTTPHandlers_WalletDetach_RejectsPrimaryWallet(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()

	primaryAddress, _ := signWalletMessageForScalar(t, "handler-detach-exec-primary-only", "41")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: identityStore,
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/wallets/detach", strings.NewReader(`{"wallet_address":"`+primaryAddress+`"}`))
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.WalletDetach(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload walletDetachConflictPayload
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if payload.Error != "wallet_detach_not_eligible" {
		t.Fatalf("unexpected error code: %q", payload.Error)
	}
	if payload.Check == nil || payload.Check.Eligible {
		t.Fatalf("expected ineligible check payload, got %#v", payload.Check)
	}
	if len(payload.Check.Reasons) == 0 || payload.Check.Reasons[0] != WalletDetachReasonWalletIsPrimary {
		t.Fatalf("unexpected detach reasons: %#v", payload.Check)
	}
}
func TestHTTPHandlers_Wallets_ReattachedWalletPreservesDetachedAt(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	address := testWalletAddress()

	identity, err := store.GetOrCreate(context.Background(), address)
	if err != nil {
		t.Fatalf("GetOrCreate error: %v", err)
	}

	_, err = store.AttachUser(context.Background(), identity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("initial AttachUser error: %v", err)
	}

	_, _, err = store.DetachUser(context.Background(), "u_test_example_com", address)
	if err != nil {
		t.Fatalf("DetachUser error: %v", err)
	}

	relinked, err := store.AttachUser(context.Background(), identity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("reattach AttachUser error: %v", err)
	}

	h := HTTPHandlers{
		Tokens:           mustTokenService(t),
		TTL:              time.Hour,
		Users:            usermod.NewService(nil),
		WalletIdentities: store,
	}

	req := httptest.NewRequest(http.MethodGet, "/auth/wallets", nil)
	req = req.WithContext(context.WithValue(req.Context(), coreauth.ClaimsContextKey, sessionClaims()))
	rec := httptest.NewRecorder()

	h.Wallets(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var payload WalletsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if len(payload.Wallets) != 1 {
		t.Fatalf("expected 1 wallet, got %d", len(payload.Wallets))
	}
	if payload.Wallets[0].Status != "active" {
		t.Fatalf("unexpected wallet status after reattach: %q", payload.Wallets[0].Status)
	}
	if payload.Wallets[0].DetachedAt == nil {
		t.Fatal("expected detached_at to remain visible after reattach")
	}
	if relinked.DetachedAt == nil || !payload.Wallets[0].DetachedAt.Equal(*relinked.DetachedAt) {
		t.Fatalf("unexpected detached_at payload: %#v", payload.Wallets[0].DetachedAt)
	}
}
