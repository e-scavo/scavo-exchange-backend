package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

func testWalletAddress() string {
	d := mustHexBig("1")
	x, y := secp256k1.ScalarBaseMult(d.Bytes())
	return publicKeyToAddress(x, y)
}

func signWalletMessageForTest(t *testing.T, message string) (string, string) {
	t.Helper()

	d := mustHexBig("1")
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

	t.Fatal("failed to derive recovery id for test signature")
	return "", ""
}

func TestRecoverWalletAddress_Success(t *testing.T) {
	message := "test wallet sign-in message"
	address, signature := signWalletMessageForTest(t, message)

	recovered, err := recoverWalletAddress(message, signature)
	if err != nil {
		t.Fatalf("unexpected recover error: %v", err)
	}
	if recovered != address {
		t.Fatalf("unexpected recovered address: %s", recovered)
	}
}

func TestWalletVerificationService_VerifyAndLogin_Success(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)

	tokens, err := coreauth.NewTokenService("dev_dev_dev_dev_dev_dev_dev_dev", "scavo-exchange-backend", time.Hour)
	if err != nil {
		t.Fatalf("NewTokenService error: %v", err)
	}

	users := usermod.NewService(&stubUserRepo{})
	loginSvc := NewService(tokens, users, time.Hour)

	identityStore := NewInMemoryWalletIdentityStore()
	verifySvc := NewWalletVerificationService(challengeSvc, loginSvc, identityStore)

	address := testWalletAddress()
	challenge, err := challengeSvc.Create(context.Background(), address, "scavium")
	if err != nil {
		t.Fatalf("challenge create error: %v", err)
	}

	_, signature := signWalletMessageForTest(t, challenge.Message)

	result, usedChallenge, err := verifySvc.VerifyAndLogin(context.Background(), challenge.ID, address, signature)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if result == nil || result.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if usedChallenge == nil || usedChallenge.UsedAt == nil {
		t.Fatal("expected used challenge")
	}
	if result.AuthMethod != "wallet_evm" {
		t.Fatalf("unexpected auth method: %q", result.AuthMethod)
	}

	if result.User == nil {
		t.Fatal("expected linked user")
	}
	if result.User.ID != walletUserID(address) {
		t.Fatalf("unexpected linked user id: %q", result.User.ID)
	}
	if result.User.Email != walletUserEmail(address) {
		t.Fatalf("unexpected linked user email: %q", result.User.Email)
	}

	claims, err := tokens.Parse(result.AccessToken)
	if err != nil {
		t.Fatalf("token parse error: %v", err)
	}
	if claims.WalletAddress != address {
		t.Fatalf("unexpected wallet address in claims: %q", claims.WalletAddress)
	}
	if claims.AuthMethod != "wallet_evm" {
		t.Fatalf("unexpected auth method in claims: %q", claims.AuthMethod)
	}

	wallets, err := identityStore.ListByUser(context.Background(), result.User.ID)
	if err != nil {
		t.Fatalf("ListByUser error: %v", err)
	}
	if len(wallets) != 1 {
		t.Fatalf("expected 1 linked wallet, got %d", len(wallets))
	}
	if !wallets[0].IsPrimary {
		t.Fatal("expected primary wallet")
	}
	if wallets[0].LinkedAt == nil {
		t.Fatal("expected linked_at")
	}
}

func TestWalletVerificationService_VerifyAndLogin_RejectsReplay(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)

	users := usermod.NewService(&stubUserRepo{})
	loginSvc := NewService(newTokenServiceForTest(t), users, time.Hour)

	identityStore := NewInMemoryWalletIdentityStore()
	verifySvc := NewWalletVerificationService(challengeSvc, loginSvc, identityStore)

	address := testWalletAddress()
	challenge, err := challengeSvc.Create(context.Background(), address, "scavium")
	if err != nil {
		t.Fatalf("challenge create error: %v", err)
	}

	_, signature := signWalletMessageForTest(t, challenge.Message)

	if _, _, err := verifySvc.VerifyAndLogin(context.Background(), challenge.ID, address, signature); err != nil {
		t.Fatalf("first verify error: %v", err)
	}

	if _, _, err := verifySvc.VerifyAndLogin(context.Background(), challenge.ID, address, signature); err != ErrChallengeUsed {
		t.Fatalf("expected ErrChallengeUsed, got %v", err)
	}
}

func TestWalletVerificationService_VerifyAndLogin_RejectsWalletLinkChallengePurpose(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)

	users := usermod.NewService(&stubUserRepo{})
	loginSvc := NewService(newTokenServiceForTest(t), users, time.Hour)

	identityStore := NewInMemoryWalletIdentityStore()
	verifySvc := NewWalletVerificationService(challengeSvc, loginSvc, identityStore)

	address := testWalletAddress()
	challenge, err := challengeSvc.CreateWithOptions(context.Background(), address, "scavium", WalletChallengeOptions{
		Purpose: WalletChallengePurposeLinkWallet,
	})
	if err != nil {
		t.Fatalf("challenge create error: %v", err)
	}

	_, signature := signWalletMessageForTest(t, challenge.Message)

	if _, _, err := verifySvc.VerifyAndLogin(context.Background(), challenge.ID, address, signature); err != ErrWalletChallengePurpose {
		t.Fatalf("expected ErrWalletChallengePurpose, got %v", err)
	}
}

func TestWalletVerificationService_VerifyAndLogin_RejectsAccountMergeChallengePurpose(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)

	users := usermod.NewService(&stubUserRepo{})
	loginSvc := NewService(newTokenServiceForTest(t), users, time.Hour)

	identityStore := NewInMemoryWalletIdentityStore()
	verifySvc := NewWalletVerificationService(challengeSvc, loginSvc, identityStore)

	address := testWalletAddress()
	challenge, err := challengeSvc.CreateWithOptions(context.Background(), address, "scavium", WalletChallengeOptions{
		Purpose: WalletChallengePurposeAccountMerge,
	})
	if err != nil {
		t.Fatalf("challenge create error: %v", err)
	}

	_, signature := signWalletMessageForTest(t, challenge.Message)

	if _, _, err := verifySvc.VerifyAndLogin(context.Background(), challenge.ID, address, signature); err != ErrWalletChallengePurpose {
		t.Fatalf("expected ErrWalletChallengePurpose, got %v", err)
	}
}

func TestInMemoryWalletIdentityStore_AttachUser_RejectsReassign(t *testing.T) {
	store := NewInMemoryWalletIdentityStore()
	address := testWalletAddress()

	identity, err := store.GetOrCreate(context.Background(), address)
	if err != nil {
		t.Fatalf("GetOrCreate error: %v", err)
	}

	_, err = store.AttachUser(context.Background(), identity.ID, "u_wallet_owner_a", true)
	if err != nil {
		t.Fatalf("first AttachUser error: %v", err)
	}

	_, err = store.AttachUser(context.Background(), identity.ID, "u_wallet_owner_b", true)
	if err != ErrWalletIdentityAlreadyLinked {
		t.Fatalf("expected ErrWalletIdentityAlreadyLinked, got %v", err)
	}
}
