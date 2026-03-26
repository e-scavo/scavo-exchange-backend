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
