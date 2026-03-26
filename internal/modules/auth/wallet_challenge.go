package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidWalletAddress        = errors.New("invalid wallet address")
	ErrChallengeStore              = errors.New("wallet challenge store error")
	ErrChallengeExpired            = errors.New("wallet challenge expired")
	ErrChallengeUsed               = errors.New("wallet challenge already used")
	ErrWalletChallengeNotFound     = errors.New("wallet challenge not found")
	ErrWalletIdentityNotFound      = errors.New("wallet identity not found")
	ErrWalletIdentityAlreadyLinked = errors.New("wallet identity already linked to another user")
	ErrWalletAlreadyLinkedToUser   = errors.New("wallet identity already linked to current user")
	ErrWalletLinkChallengeMismatch = errors.New("wallet link challenge does not belong to current user")
	ErrWalletChallengePurpose      = errors.New("wallet challenge purpose mismatch")
	ErrWalletMergeSourceNotLinked  = errors.New("wallet merge source wallet is not linked to another user")
	ErrWalletMergeSameUser         = errors.New("wallet merge source already belongs to current user")
)

const (
	WalletChallengePurposeAuthBootstrap = "auth_bootstrap"
	WalletChallengePurposeLinkWallet    = "wallet_link"
	WalletChallengePurposeAccountMerge  = "account_merge"
)

var evmAddressRE = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

type WalletChallenge struct {
	ID                string     `json:"id"`
	Address           string     `json:"address"`
	Chain             string     `json:"chain"`
	Nonce             string     `json:"nonce"`
	Message           string     `json:"message"`
	Purpose           string     `json:"purpose"`
	RequestedByUserID string     `json:"requested_by_user_id,omitempty"`
	IssuedAt          time.Time  `json:"issued_at"`
	ExpiresAt         time.Time  `json:"expires_at"`
	UsedAt            *time.Time `json:"used_at,omitempty"`
}

type WalletChallengeStore interface {
	Save(ctx context.Context, challenge *WalletChallenge) error
	GetByID(ctx context.Context, id string) (*WalletChallenge, error)
	Use(ctx context.Context, id string, usedAt time.Time) (*WalletChallenge, error)
}

type WalletChallengeOptions struct {
	Purpose           string
	RequestedByUserID string
}

type WalletChallengeService struct {
	store         WalletChallengeStore
	publicBaseURL string
	ttl           time.Duration
}

func NewWalletChallengeService(store WalletChallengeStore, publicBaseURL string, ttl time.Duration) *WalletChallengeService {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}

	return &WalletChallengeService{
		store:         store,
		publicBaseURL: strings.TrimSpace(publicBaseURL),
		ttl:           ttl,
	}
}

func (s *WalletChallengeService) Create(ctx context.Context, address, chain string) (*WalletChallenge, error) {
	return s.CreateWithOptions(ctx, address, chain, WalletChallengeOptions{
		Purpose: WalletChallengePurposeAuthBootstrap,
	})
}

func (s *WalletChallengeService) CreateWithOptions(ctx context.Context, address, chain string, options WalletChallengeOptions) (*WalletChallenge, error) {
	address = normalizeWalletAddress(address)
	if !evmAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}

	chain = normalizeChain(chain)
	purpose := normalizeWalletChallengePurpose(options.Purpose)
	now := time.Now().UTC()
	expiresAt := now.Add(s.ttl)

	nonce, err := randomToken(24)
	if err != nil {
		return nil, err
	}

	challenge := &WalletChallenge{
		ID:                uuid.NewString(),
		Address:           address,
		Chain:             chain,
		Nonce:             nonce,
		Purpose:           purpose,
		RequestedByUserID: strings.TrimSpace(options.RequestedByUserID),
		IssuedAt:          now,
		ExpiresAt:         expiresAt,
	}
	challenge.Message = s.buildMessage(challenge)

	if s.store != nil {
		if err := s.store.Save(ctx, challenge); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrChallengeStore, err)
		}
	}

	return challenge, nil
}

func (s *WalletChallengeService) Get(ctx context.Context, id string) (*WalletChallenge, error) {
	if s == nil || s.store == nil {
		return nil, ErrChallengeStore
	}

	challenge, err := s.store.GetByID(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrWalletChallengeNotFound
	}
	if challenge.UsedAt != nil {
		return nil, ErrChallengeUsed
	}
	if time.Now().UTC().After(challenge.ExpiresAt) {
		return nil, ErrChallengeExpired
	}

	challenge.Purpose = normalizeWalletChallengePurpose(challenge.Purpose)
	challenge.RequestedByUserID = strings.TrimSpace(challenge.RequestedByUserID)

	return challenge, nil
}

func (s *WalletChallengeService) MarkUsed(ctx context.Context, id string, usedAt time.Time) (*WalletChallenge, error) {
	if s == nil || s.store == nil {
		return nil, ErrChallengeStore
	}

	challenge, err := s.store.Use(ctx, strings.TrimSpace(id), usedAt.UTC())
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrWalletChallengeNotFound
	}

	challenge.Purpose = normalizeWalletChallengePurpose(challenge.Purpose)
	challenge.RequestedByUserID = strings.TrimSpace(challenge.RequestedByUserID)

	return challenge, nil
}

func (s *WalletChallengeService) buildMessage(ch *WalletChallenge) string {
	domain := "SCAVO Exchange"
	uri := "http://localhost"

	if s.publicBaseURL != "" {
		if u, err := url.Parse(s.publicBaseURL); err == nil {
			if host := strings.TrimSpace(u.Host); host != "" {
				domain = host
			}
			uri = s.publicBaseURL
		}
	}

	purposeLine := "Purpose: SCAVO Exchange wallet authentication bootstrap."
	switch normalizeWalletChallengePurpose(ch.Purpose) {
	case WalletChallengePurposeLinkWallet:
		purposeLine = "Purpose: SCAVO Exchange authenticated wallet linking confirmation."
	case WalletChallengePurposeAccountMerge:
		purposeLine = "Purpose: SCAVO Exchange authenticated account merge confirmation."
	}

	lines := []string{
		fmt.Sprintf("%s wants you to sign in with your wallet.", domain),
		"",
		fmt.Sprintf("Address: %s", ch.Address),
		fmt.Sprintf("Chain: %s", ch.Chain),
		fmt.Sprintf("Nonce: %s", ch.Nonce),
		fmt.Sprintf("Issued At: %s", ch.IssuedAt.Format(time.RFC3339)),
		fmt.Sprintf("Expiration Time: %s", ch.ExpiresAt.Format(time.RFC3339)),
		fmt.Sprintf("URI: %s", uri),
		"",
		purposeLine,
	}

	if requestedBy := strings.TrimSpace(ch.RequestedByUserID); requestedBy != "" {
		lines = append(lines, fmt.Sprintf("Requested By User ID: %s", requestedBy))
	}

	return strings.Join(lines, "\n")
}

func normalizeWalletChallengePurpose(purpose string) string {
	switch strings.TrimSpace(strings.ToLower(purpose)) {
	case "", WalletChallengePurposeAuthBootstrap:
		return WalletChallengePurposeAuthBootstrap
	case WalletChallengePurposeLinkWallet:
		return WalletChallengePurposeLinkWallet
	case WalletChallengePurposeAccountMerge:
		return WalletChallengePurposeAccountMerge
	default:
		return WalletChallengePurposeAuthBootstrap
	}
}

func normalizeChain(chain string) string {
	chain = strings.TrimSpace(strings.ToLower(chain))
	if chain == "" {
		return "scavium"
	}
	return chain
}

func randomToken(numBytes int) (string, error) {
	if numBytes <= 0 {
		numBytes = 24
	}

	b := make([]byte, numBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
