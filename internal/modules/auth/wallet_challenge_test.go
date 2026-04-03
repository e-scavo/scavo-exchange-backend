package auth

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestWalletChallengeService_Create_Success(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	svc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)

	challenge, err := svc.Create(context.Background(), "0x1111111111111111111111111111111111111111", "scavium")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if challenge == nil {
		t.Fatal("expected challenge")
	}
	if challenge.ID == "" {
		t.Fatal("expected challenge id")
	}
	if challenge.Nonce == "" {
		t.Fatal("expected nonce")
	}
	if challenge.Address != "0x1111111111111111111111111111111111111111" {
		t.Fatalf("unexpected address: %q", challenge.Address)
	}
	if challenge.Chain != "scavium" {
		t.Fatalf("unexpected chain: %q", challenge.Chain)
	}
	if !strings.Contains(challenge.Message, challenge.Nonce) {
		t.Fatal("expected message to contain nonce")
	}

	got, err := store.GetByID(context.Background(), challenge.ID)
	if err != nil {
		t.Fatalf("expected saved challenge, got error: %v", err)
	}
	if got.ID != challenge.ID {
		t.Fatalf("unexpected stored id: %q", got.ID)
	}
}

func TestWalletChallengeService_Create_InvalidAddress(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	svc := NewWalletChallengeService(store, "", 5*time.Minute)

	_, err := svc.Create(context.Background(), "invalid-address", "scavium")
	if err == nil {
		t.Fatal("expected invalid wallet address error")
	}
	if err != ErrInvalidWalletAddress {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWalletChallengeService_Create_DefaultChain(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	svc := NewWalletChallengeService(store, "", 5*time.Minute)

	challenge, err := svc.Create(context.Background(), "0x1111111111111111111111111111111111111111", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if challenge.Chain != "scavium" {
		t.Fatalf("expected default chain scavium, got %q", challenge.Chain)
	}
}

func TestWalletChallengeService_CreateWithOptions_RejectsUnknownPurpose(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	svc := NewWalletChallengeService(store, "", 5*time.Minute)

	_, err := svc.CreateWithOptions(context.Background(), "0x1111111111111111111111111111111111111111", "scavium", WalletChallengeOptions{
		Purpose: "legacy_bootstrap_typo",
	})
	if err != ErrWalletChallengePurpose {
		t.Fatalf("expected ErrWalletChallengePurpose, got %v", err)
	}
}

func TestWalletChallengeService_Get_DoesNotDefaultUnknownPurposeToBootstrap(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	svc := NewWalletChallengeService(store, "", 5*time.Minute)
	now := time.Now().UTC()

	challenge := &WalletChallenge{
		ID:        "ch_invalid_purpose",
		Address:   "0x1111111111111111111111111111111111111111",
		Chain:     "scavium",
		Nonce:     "nonce",
		Message:   "invalid purpose challenge",
		Purpose:   "legacy_bootstrap_typo",
		IssuedAt:  now,
		ExpiresAt: now.Add(5 * time.Minute),
	}
	if err := store.Save(context.Background(), challenge); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	got, err := svc.Get(context.Background(), challenge.ID)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if got.Purpose != "legacy_bootstrap_typo" {
		t.Fatalf("expected unknown purpose to remain unchanged, got %q", got.Purpose)
	}
}
