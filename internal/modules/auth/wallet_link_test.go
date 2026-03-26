package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"
	"time"
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
