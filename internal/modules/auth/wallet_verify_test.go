package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
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
	loginSvc := NewService(tokens, nil, time.Hour)
	verifySvc := NewWalletVerificationService(challengeSvc, loginSvc)

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
}

func TestWalletVerificationService_VerifyAndLogin_RejectsReplay(t *testing.T) {
	store := NewInMemoryWalletChallengeStore()
	challengeSvc := NewWalletChallengeService(store, "https://api.scavo.exchange", 5*time.Minute)
	loginSvc := NewService(newTokenServiceForTest(t), nil, time.Hour)
	verifySvc := NewWalletVerificationService(challengeSvc, loginSvc)

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
