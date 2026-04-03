package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

func signWalletMessageForScalar(t *testing.T, message, scalar string) (string, string) {
	t.Helper()

	d, ok := new(big.Int).SetString(scalar, 16)
	if !ok {
		t.Fatalf("invalid scalar: %s", scalar)
	}
	x, y := secp256k1.ScalarBaseMult(d.Bytes())
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: secp256k1, X: x, Y: y},
		D:         d,
	}

	hash := ethereumMessageHash(message)
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash)
	if err != nil {
		t.Fatalf("ecdsa.Sign error: %v", err)
	}

	address := publicKeyToAddress(x, y)
	sig := make([]byte, 65)
	r.FillBytes(sig[:32])
	s.FillBytes(sig[32:64])

	for v := 0; v < 4; v++ {
		sig[64] = byte(v)
		recoveredAddress, err := recoverWalletAddress(message, "0x"+hex.EncodeToString(sig))
		if err == nil && recoveredAddress == address {
			return address, "0x" + hex.EncodeToString(sig)
		}
	}

	t.Fatal("failed to derive recovery id for scalar test signature")
	return "", ""
}

func TestWalletLinkingService_VerifyAndLink_Success(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)
	identityStore := NewInMemoryWalletIdentityStore()
	linkSvc := NewWalletLinkingService(challengeSvc, identityStore)

	primaryAddress, _ := signWalletMessageForScalar(t, "bootstrap", "1")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_test_example_com", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "bootstrap-2", "2")
	challenge, err := linkSvc.CreateChallenge(context.Background(), "u_test_example_com", secondaryAddress, "scavium")
	if err != nil {
		t.Fatalf("CreateChallenge error: %v", err)
	}
	if challenge.Purpose != WalletChallengePurposeLinkWallet {
		t.Fatalf("unexpected challenge purpose: %q", challenge.Purpose)
	}
	if challenge.RequestedByUserID != "u_test_example_com" {
		t.Fatalf("unexpected requested user id: %q", challenge.RequestedByUserID)
	}

	_, signature := signWalletMessageForScalar(t, challenge.Message, "2")
	result, err := linkSvc.VerifyAndLink(context.Background(), "u_test_example_com", challenge.ID, secondaryAddress, signature)
	if err != nil {
		t.Fatalf("VerifyAndLink error: %v", err)
	}
	if result.Linked == nil {
		t.Fatal("expected linked wallet")
	}
	if result.Linked.UserID != "u_test_example_com" {
		t.Fatalf("unexpected linked wallet user id: %q", result.Linked.UserID)
	}
	if result.Linked.IsPrimary {
		t.Fatal("expected linked wallet to be secondary")
	}
	if result.Linked.LinkedAt == nil {
		t.Fatal("expected linked_at on linked wallet")
	}
	if len(result.Wallets) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(result.Wallets))
	}
	if !result.Wallets[0].IsPrimary {
		t.Fatal("expected first wallet to remain primary")
	}
}

func TestWalletLinkingService_VerifyAndLink_RejectsAlreadyLinkedToOtherUser(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)
	identityStore := NewInMemoryWalletIdentityStore()
	linkSvc := NewWalletLinkingService(challengeSvc, identityStore)

	address, _ := signWalletMessageForScalar(t, "bootstrap", "1")
	identity, err := identityStore.GetOrCreate(context.Background(), address)
	if err != nil {
		t.Fatalf("GetOrCreate error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), identity.ID, "u_other", true)
	if err != nil {
		t.Fatalf("AttachUser error: %v", err)
	}

	challenge, err := linkSvc.CreateChallenge(context.Background(), "u_test_example_com", address, "scavium")
	if err != nil {
		t.Fatalf("CreateChallenge error: %v", err)
	}
	_, signature := signWalletMessageForScalar(t, challenge.Message, "1")

	_, err = linkSvc.VerifyAndLink(context.Background(), "u_test_example_com", challenge.ID, address, signature)
	if err != ErrWalletIdentityAlreadyLinked {
		t.Fatalf("expected ErrWalletIdentityAlreadyLinked, got %v", err)
	}
}

func TestWalletAccountMergeService_VerifyAndMerge_Success(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)
	identityStore := NewInMemoryWalletIdentityStore()
	mergeSvc := NewWalletAccountMergeService(challengeSvc, identityStore)

	targetPrimaryAddress, _ := signWalletMessageForScalar(t, "target-primary", "1")
	targetPrimaryIdentity, err := identityStore.GetOrCreate(context.Background(), targetPrimaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate target primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), targetPrimaryIdentity.ID, "u_target", true)
	if err != nil {
		t.Fatalf("AttachUser target primary error: %v", err)
	}

	sourcePrimaryAddress, _ := signWalletMessageForScalar(t, "source-primary", "2")
	sourcePrimaryIdentity, err := identityStore.GetOrCreate(context.Background(), sourcePrimaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate source primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), sourcePrimaryIdentity.ID, "u_source", true)
	if err != nil {
		t.Fatalf("AttachUser source primary error: %v", err)
	}

	sourceSecondaryAddress, _ := signWalletMessageForScalar(t, "source-secondary", "3")
	sourceSecondaryIdentity, err := identityStore.GetOrCreate(context.Background(), sourceSecondaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate source secondary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), sourceSecondaryIdentity.ID, "u_source", false)
	if err != nil {
		t.Fatalf("AttachUser source secondary error: %v", err)
	}

	challenge, err := mergeSvc.CreateChallenge(context.Background(), "u_target", sourcePrimaryAddress, "scavium")
	if err != nil {
		t.Fatalf("CreateChallenge error: %v", err)
	}
	if challenge.Purpose != WalletChallengePurposeAccountMerge {
		t.Fatalf("unexpected challenge purpose: %q", challenge.Purpose)
	}
	_, signature := signWalletMessageForScalar(t, challenge.Message, "2")

	result, err := mergeSvc.VerifyAndMerge(context.Background(), "u_target", challenge.ID, sourcePrimaryAddress, signature)
	if err != nil {
		t.Fatalf("VerifyAndMerge error: %v", err)
	}
	if result.SourceUserID != "u_source" {
		t.Fatalf("unexpected source user id: %q", result.SourceUserID)
	}
	if result.TargetUserID != "u_target" {
		t.Fatalf("unexpected target user id: %q", result.TargetUserID)
	}
	if result.MergedWallet == nil {
		t.Fatal("expected merged wallet")
	}
	if result.MergedWallet.UserID != "u_target" {
		t.Fatalf("unexpected merged wallet user id: %q", result.MergedWallet.UserID)
	}
	if len(result.Wallets) != 3 {
		t.Fatalf("expected 3 wallets after merge, got %d", len(result.Wallets))
	}
	if !result.Wallets[0].IsPrimary || result.Wallets[0].Address != targetPrimaryAddress {
		t.Fatal("expected target primary wallet to remain primary")
	}

	sourceWallets, err := identityStore.ListByUser(context.Background(), "u_source")
	if err != nil {
		t.Fatalf("ListByUser source error: %v", err)
	}
	if len(sourceWallets) != 0 {
		t.Fatalf("expected source user to have 0 wallets after merge, got %d", len(sourceWallets))
	}
}

func TestWalletAccountMergeService_VerifyAndMerge_RequiresLinkedSource(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)
	identityStore := NewInMemoryWalletIdentityStore()
	mergeSvc := NewWalletAccountMergeService(challengeSvc, identityStore)

	address, _ := signWalletMessageForScalar(t, "unlinked-source", "4")
	_, err := identityStore.GetOrCreate(context.Background(), address)
	if err != nil {
		t.Fatalf("GetOrCreate error: %v", err)
	}

	challenge, err := mergeSvc.CreateChallenge(context.Background(), "u_target", address, "scavium")
	if err != nil {
		t.Fatalf("CreateChallenge error: %v", err)
	}
	_, signature := signWalletMessageForScalar(t, challenge.Message, "4")

	_, err = mergeSvc.VerifyAndMerge(context.Background(), "u_target", challenge.ID, address, signature)
	if err != ErrWalletMergeSourceNotLinked {
		t.Fatalf("expected ErrWalletMergeSourceNotLinked, got %v", err)
	}
}

func TestWalletPrimaryService_SetPrimary_Success(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()
	primarySvc := NewWalletPrimaryService(identityStore)

	primaryAddress, _ := signWalletMessageForScalar(t, "primary", "10")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_target", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "secondary", "11")
	secondaryIdentity, err := identityStore.GetOrCreate(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate secondary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), secondaryIdentity.ID, "u_target", false)
	if err != nil {
		t.Fatalf("AttachUser secondary error: %v", err)
	}

	result, err := primarySvc.SetPrimary(context.Background(), "u_target", secondaryAddress)
	if err != nil {
		t.Fatalf("SetPrimary error: %v", err)
	}
	if result.Primary == nil || result.Primary.Address != secondaryAddress {
		t.Fatalf("unexpected primary wallet: %#v", result.Primary)
	}
	if !result.Primary.IsPrimary {
		t.Fatal("expected switched wallet to be primary")
	}
	if len(result.Wallets) != 2 {
		t.Fatalf("expected 2 wallets, got %d", len(result.Wallets))
	}
	if result.Wallets[0].Address != secondaryAddress || !result.Wallets[0].IsPrimary {
		t.Fatal("expected switched wallet to be first and primary")
	}
	if result.Wallets[1].Address != primaryAddress || result.Wallets[1].IsPrimary {
		t.Fatal("expected old primary to become secondary")
	}
}

func TestWalletPrimaryService_SetPrimary_RejectsWalletNotOwnedByUser(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()
	primarySvc := NewWalletPrimaryService(identityStore)

	otherAddress, _ := signWalletMessageForScalar(t, "other", "12")
	otherIdentity, err := identityStore.GetOrCreate(context.Background(), otherAddress)
	if err != nil {
		t.Fatalf("GetOrCreate other error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), otherIdentity.ID, "u_other", true)
	if err != nil {
		t.Fatalf("AttachUser other error: %v", err)
	}

	_, err = primarySvc.SetPrimary(context.Background(), "u_target", otherAddress)
	if err != ErrWalletNotOwnedByUser {
		t.Fatalf("expected ErrWalletNotOwnedByUser, got %v", err)
	}
}

func TestWalletDetachService_CheckEligibility_RejectsPrimaryOnlyWallet(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()
	detachSvc := NewWalletDetachService(identityStore)

	address, _ := signWalletMessageForScalar(t, "detach-primary-only", "30")
	identity, err := identityStore.GetOrCreate(context.Background(), address)
	if err != nil {
		t.Fatalf("GetOrCreate error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), identity.ID, "u_target", true)
	if err != nil {
		t.Fatalf("AttachUser error: %v", err)
	}

	result, err := detachSvc.CheckEligibility(context.Background(), "u_target", address)
	if err != nil {
		t.Fatalf("CheckEligibility error: %v", err)
	}
	if result.Eligible {
		t.Fatal("expected ineligible detach result")
	}
	if !result.IsPrimary {
		t.Fatal("expected wallet to be primary")
	}
	if result.OwnedWalletCount != 1 {
		t.Fatalf("expected owned wallet count 1, got %d", result.OwnedWalletCount)
	}
	if len(result.Reasons) != 2 {
		t.Fatalf("expected 2 reasons, got %d", len(result.Reasons))
	}
	if result.Reasons[0] != WalletDetachReasonWalletIsPrimary {
		t.Fatalf("unexpected first reason: %q", result.Reasons[0])
	}
	if result.Reasons[1] != WalletDetachReasonUserWouldBeEmpty {
		t.Fatalf("unexpected second reason: %q", result.Reasons[1])
	}
}

func TestWalletDetachService_CheckEligibility_SuccessForOwnedSecondaryWallet(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()
	detachSvc := NewWalletDetachService(identityStore)

	primaryAddress, _ := signWalletMessageForScalar(t, "detach-primary", "31")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_target", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "detach-secondary", "32")
	secondaryIdentity, err := identityStore.GetOrCreate(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate secondary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), secondaryIdentity.ID, "u_target", false)
	if err != nil {
		t.Fatalf("AttachUser secondary error: %v", err)
	}

	result, err := detachSvc.CheckEligibility(context.Background(), "u_target", secondaryAddress)
	if err != nil {
		t.Fatalf("CheckEligibility error: %v", err)
	}
	if !result.Eligible {
		t.Fatalf("expected eligible detach result, got reasons=%v", result.Reasons)
	}
	if result.IsPrimary {
		t.Fatal("expected secondary wallet to be non-primary")
	}
	if result.OwnedWalletCount != 2 {
		t.Fatalf("expected owned wallet count 2, got %d", result.OwnedWalletCount)
	}
	if len(result.Reasons) != 0 {
		t.Fatalf("expected no reasons, got %v", result.Reasons)
	}
}

func TestWalletDetachService_Execute_SuccessForOwnedSecondaryWallet(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()
	detachSvc := NewWalletDetachService(identityStore)

	primaryAddress, _ := signWalletMessageForScalar(t, "detach-exec-primary", "36")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_target", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "detach-exec-secondary", "37")
	secondaryIdentity, err := identityStore.GetOrCreate(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate secondary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), secondaryIdentity.ID, "u_target", false)
	if err != nil {
		t.Fatalf("AttachUser secondary error: %v", err)
	}

	result, err := detachSvc.Execute(context.Background(), "u_target", secondaryAddress)
	if err != nil {
		t.Fatalf("Execute error: %v", err)
	}
	if result.Detached == nil || result.Detached.Address != secondaryAddress {
		t.Fatalf("unexpected detached wallet: %#v", result.Detached)
	}
	if result.Detached.UserID != "" {
		t.Fatalf("expected detached wallet to have empty user id, got %q", result.Detached.UserID)
	}
	if result.Detached.IsPrimary {
		t.Fatal("expected detached wallet to be non-primary")
	}
	if result.Detached.LinkedAt != nil {
		t.Fatal("expected detached wallet linked_at to be cleared")
	}
	if len(result.Wallets) != 1 {
		t.Fatalf("expected 1 remaining wallet, got %d", len(result.Wallets))
	}
	if result.Wallets[0].Address != primaryAddress || !result.Wallets[0].IsPrimary {
		t.Fatal("expected original primary wallet to remain attached and primary")
	}

	storedDetached, err := identityStore.GetByAddress(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetByAddress detached error: %v", err)
	}
	if storedDetached.UserID != "" || storedDetached.IsPrimary || storedDetached.LinkedAt != nil {
		t.Fatalf("unexpected stored detached wallet state: %#v", storedDetached)
	}
}

func TestWalletLinkingService_VerifyAndLink_AllowsDetachedWalletReattachment(t *testing.T) {
	challengeStore := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(challengeStore, "https://api.scavo.exchange", 5*time.Minute)
	identityStore := NewInMemoryWalletIdentityStore()
	linkSvc := NewWalletLinkingService(challengeSvc, identityStore)
	detachSvc := NewWalletDetachService(identityStore)

	primaryAddress, _ := signWalletMessageForScalar(t, "reattach-primary", "42")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_target", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "reattach-secondary", "43")
	secondaryIdentity, err := identityStore.GetOrCreate(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate secondary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), secondaryIdentity.ID, "u_target", false)
	if err != nil {
		t.Fatalf("AttachUser secondary error: %v", err)
	}

	detachResult, err := detachSvc.Execute(context.Background(), "u_target", secondaryAddress)
	if err != nil {
		t.Fatalf("Execute detach error: %v", err)
	}
	if detachResult == nil || detachResult.Detached == nil {
		t.Fatal("expected detached wallet result")
	}
	if len(detachResult.Wallets) != 1 {
		t.Fatalf("expected 1 remaining wallet after detach, got %d", len(detachResult.Wallets))
	}

	challenge, err := linkSvc.CreateChallenge(context.Background(), "u_target", secondaryAddress, "scavium")
	if err != nil {
		t.Fatalf("CreateChallenge error: %v", err)
	}
	_, signature := signWalletMessageForScalar(t, challenge.Message, "43")

	result, err := linkSvc.VerifyAndLink(context.Background(), "u_target", challenge.ID, secondaryAddress, signature)
	if err != nil {
		t.Fatalf("VerifyAndLink error: %v", err)
	}
	if result.Linked == nil {
		t.Fatal("expected linked wallet")
	}
	if result.Linked.UserID != "u_target" {
		t.Fatalf("unexpected relinked user id: %q", result.Linked.UserID)
	}
	if result.Linked.IsPrimary {
		t.Fatal("expected reattached wallet to remain non-primary")
	}
	if result.Linked.LinkedAt == nil {
		t.Fatal("expected reattached wallet to get linked_at again")
	}
	if len(result.Wallets) != 2 {
		t.Fatalf("expected 2 wallets after reattachment, got %d", len(result.Wallets))
	}

	stored, err := identityStore.GetByAddress(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetByAddress relinked error: %v", err)
	}
	if stored.UserID != "u_target" {
		t.Fatalf("unexpected stored relinked user id: %q", stored.UserID)
	}
	if stored.LinkedAt == nil {
		t.Fatal("expected stored relinked wallet linked_at")
	}
	if stored.DetachedAt == nil {
		t.Fatal("expected relinked wallet to preserve detached_at")
	}
	if stored.IsPrimary {
		t.Fatal("expected stored relinked wallet to remain non-primary")
	}
}

func TestWalletVerificationService_VerifyAndLogin_RebindsDetachedWalletToWalletUser(t *testing.T) {
	challengeStore := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(challengeStore, "https://api.scavo.exchange", 5*time.Minute)
	identityStore := NewInMemoryWalletIdentityStore()
	detachSvc := NewWalletDetachService(identityStore)
	loginSvc := NewService(newTokenServiceForTest(t), usermod.NewService(&stubUserRepo{}), time.Hour)
	verifySvc := NewWalletVerificationService(challengeSvc, loginSvc, identityStore)

	primaryAddress, _ := signWalletMessageForScalar(t, "wallet-login-primary", "44")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_existing_owner", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	secondaryAddress, _ := signWalletMessageForScalar(t, "wallet-login-secondary", "45")
	secondaryIdentity, err := identityStore.GetOrCreate(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate secondary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), secondaryIdentity.ID, "u_existing_owner", false)
	if err != nil {
		t.Fatalf("AttachUser secondary error: %v", err)
	}

	detachResult, err := detachSvc.Execute(context.Background(), "u_existing_owner", secondaryAddress)
	if err != nil {
		t.Fatalf("Execute secondary detach error: %v", err)
	}
	if detachResult == nil || detachResult.Detached == nil || detachResult.Detached.DetachedAt == nil {
		t.Fatal("expected detached wallet metadata with detached_at")
	}

	challenge, err := challengeSvc.Create(context.Background(), secondaryAddress, "scavium")
	if err != nil {
		t.Fatalf("Create challenge error: %v", err)
	}
	_, signature := signWalletMessageForScalar(t, challenge.Message, "45")

	result, usedChallenge, err := verifySvc.VerifyAndLogin(context.Background(), challenge.ID, secondaryAddress, signature)
	if err != nil {
		t.Fatalf("VerifyAndLogin error: %v", err)
	}
	if result == nil || result.User == nil {
		t.Fatal("expected wallet login result with user")
	}
	if result.User.ID != walletUserID(secondaryAddress) {
		t.Fatalf("expected detached wallet login to resolve wallet-owned user, got %q", result.User.ID)
	}
	if result.AuthMethod != "wallet_evm" {
		t.Fatalf("unexpected auth method: %q", result.AuthMethod)
	}
	if usedChallenge == nil || usedChallenge.UsedAt == nil {
		t.Fatal("expected used challenge metadata")
	}

	stored, err := identityStore.GetByAddress(context.Background(), secondaryAddress)
	if err != nil {
		t.Fatalf("GetByAddress relogin error: %v", err)
	}
	if stored.UserID == "" {
		t.Fatal("expected wallet login to rebind detached wallet")
	}
	if stored.UserID != result.User.ID {
		t.Fatalf("unexpected rebound user id: %q vs %q", stored.UserID, result.User.ID)
	}
	if stored.LinkedAt == nil {
		t.Fatal("expected rebound wallet linked_at")
	}
	if stored.DetachedAt == nil {
		t.Fatal("expected rebound wallet to preserve detached_at")
	}
	if !stored.IsPrimary {
		t.Fatal("expected rebound wallet to become primary for wallet-owned user")
	}
}

func TestWalletDetachService_Execute_RejectsPrimaryWallet(t *testing.T) {
	identityStore := NewInMemoryWalletIdentityStore()
	detachSvc := NewWalletDetachService(identityStore)

	primaryAddress, _ := signWalletMessageForScalar(t, "detach-exec-primary-only", "38")
	primaryIdentity, err := identityStore.GetOrCreate(context.Background(), primaryAddress)
	if err != nil {
		t.Fatalf("GetOrCreate primary error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), primaryIdentity.ID, "u_target", true)
	if err != nil {
		t.Fatalf("AttachUser primary error: %v", err)
	}

	_, err = detachSvc.Execute(context.Background(), "u_target", primaryAddress)
	if err != ErrWalletDetachNotEligible {
		t.Fatalf("expected ErrWalletDetachNotEligible, got %v", err)
	}
}

func TestWalletLinkingService_VerifyAndLink_RejectsUnknownChallengePurpose(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)
	identityStore := NewInMemoryWalletIdentityStore()
	linkSvc := NewWalletLinkingService(challengeSvc, identityStore)

	address, _ := signWalletMessageForScalar(t, "unknown-link", "12")
	now := time.Now().UTC()
	challenge := &WalletChallenge{
		ID:                "link_unknown_purpose",
		Address:           address,
		Chain:             "scavium",
		Nonce:             "nonce_unknown_link",
		Message:           "SCAVO Exchange legacy link purpose",
		Purpose:           "legacy_link_typo",
		RequestedByUserID: "u_target",
		IssuedAt:          now,
		ExpiresAt:         now.Add(5 * time.Minute),
	}
	if err := store.Save(context.Background(), challenge); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	_, signature := signWalletMessageForScalar(t, challenge.Message, "12")

	_, err := linkSvc.VerifyAndLink(context.Background(), "u_target", challenge.ID, address, signature)
	if err != ErrWalletChallengePurpose {
		t.Fatalf("expected ErrWalletChallengePurpose, got %v", err)
	}
}

func TestWalletAccountMergeService_VerifyAndMerge_RejectsUnknownChallengePurpose(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)
	identityStore := NewInMemoryWalletIdentityStore()
	mergeSvc := NewWalletAccountMergeService(challengeSvc, identityStore)

	sourceAddress, _ := signWalletMessageForScalar(t, "unknown-merge", "13")
	sourceIdentity, err := identityStore.GetOrCreate(context.Background(), sourceAddress)
	if err != nil {
		t.Fatalf("GetOrCreate error: %v", err)
	}
	_, err = identityStore.AttachUser(context.Background(), sourceIdentity.ID, "u_source", true)
	if err != nil {
		t.Fatalf("AttachUser error: %v", err)
	}

	now := time.Now().UTC()
	challenge := &WalletChallenge{
		ID:                "merge_unknown_purpose",
		Address:           sourceAddress,
		Chain:             "scavium",
		Nonce:             "nonce_unknown_merge",
		Message:           "SCAVO Exchange legacy merge purpose",
		Purpose:           "legacy_merge_typo",
		RequestedByUserID: "u_target",
		IssuedAt:          now,
		ExpiresAt:         now.Add(5 * time.Minute),
	}
	if err := store.Save(context.Background(), challenge); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	_, signature := signWalletMessageForScalar(t, challenge.Message, "13")

	_, err = mergeSvc.VerifyAndMerge(context.Background(), "u_target", challenge.ID, sourceAddress, signature)
	if err != ErrWalletChallengePurpose {
		t.Fatalf("expected ErrWalletChallengePurpose, got %v", err)
	}
}
