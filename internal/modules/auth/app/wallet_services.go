package app

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

	authdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/domain"
)

var (
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
	ErrWalletNotOwnedByUser        = errors.New("wallet identity does not belong to current user")
	ErrWalletDetachNotEligible     = errors.New("wallet detach not eligible under current ownership guardrails")
	ErrInvalidWalletSignature      = errors.New("invalid wallet signature")
)

const (
	WalletChallengePurposeAuthBootstrap = "auth_bootstrap"
	WalletChallengePurposeLinkWallet    = "wallet_link"
	WalletChallengePurposeAccountMerge  = "account_merge"
)

var walletEVMAddressRE = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

type WalletChallengeService struct {
	store         authdomain.WalletChallengeStore
	publicBaseURL string
	ttl           time.Duration
}

type WalletLinkingService struct {
	challenges WalletChallengeService
	identities authdomain.WalletIdentityStore
}

type WalletLinkResult struct {
	Challenge *authdomain.WalletChallenge  `json:"challenge,omitempty"`
	Linked    *authdomain.WalletIdentity   `json:"linked_wallet,omitempty"`
	Wallets   []*authdomain.WalletIdentity `json:"wallets"`
}

type WalletAccountMergeService struct {
	challenges WalletChallengeService
	identities authdomain.WalletIdentityStore
}

type WalletAccountMergeResult struct {
	Challenge    *authdomain.WalletChallenge  `json:"challenge,omitempty"`
	MergedWallet *authdomain.WalletIdentity   `json:"merged_wallet,omitempty"`
	Wallets      []*authdomain.WalletIdentity `json:"wallets"`
	SourceUserID string                       `json:"source_user_id"`
	TargetUserID string                       `json:"target_user_id"`
}

type WalletPrimaryService struct {
	identities authdomain.WalletIdentityStore
}

type WalletPrimaryResult struct {
	Primary *authdomain.WalletIdentity   `json:"primary_wallet,omitempty"`
	Wallets []*authdomain.WalletIdentity `json:"wallets"`
}

type WalletDetachCheckResult struct {
	WalletAddress    string   `json:"wallet_address"`
	Eligible         bool     `json:"eligible"`
	IsPrimary        bool     `json:"is_primary"`
	OwnedWalletCount int      `json:"owned_wallet_count"`
	Reasons          []string `json:"reasons"`
}

type WalletDetachExecuteResult struct {
	Detached *authdomain.WalletIdentity   `json:"detached,omitempty"`
	Wallets  []*authdomain.WalletIdentity `json:"wallets"`
	Check    *WalletDetachCheckResult     `json:"check,omitempty"`
}

type WalletDetachService struct {
	identities authdomain.WalletIdentityStore
}

func NewWalletChallengeService(store authdomain.WalletChallengeStore, publicBaseURL string, ttl time.Duration) *WalletChallengeService {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &WalletChallengeService{store: store, publicBaseURL: strings.TrimSpace(publicBaseURL), ttl: ttl}
}

func (s *WalletChallengeService) Create(ctx context.Context, address, chain string) (*authdomain.WalletChallenge, error) {
	return s.CreateWithOptions(ctx, address, chain, authdomain.WalletChallengeOptions{Purpose: WalletChallengePurposeAuthBootstrap})
}

func (s *WalletChallengeService) CreateWithOptions(ctx context.Context, address, chain string, options authdomain.WalletChallengeOptions) (*authdomain.WalletChallenge, error) {
	if s == nil || s.store == nil {
		return nil, ErrChallengeStore
	}
	address = normalizeWalletAddress(address)
	if !walletEVMAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}
	chain = normalizeChain(chain)
	purpose, err := resolveWalletChallengePurposeForCreate(options.Purpose)
	if err != nil {
		return nil, err
	}
	challenge, err := s.store.CreateWithOptions(ctx, address, chain, authdomain.WalletChallengeOptions{
		Purpose:           purpose,
		RequestedByUserID: strings.TrimSpace(options.RequestedByUserID),
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrChallengeStore, err)
	}
	return challenge, nil
}

func (s *WalletChallengeService) Get(ctx context.Context, id string) (*authdomain.WalletChallenge, error) {
	if s == nil || s.store == nil {
		return nil, ErrChallengeStore
	}
	challenge, err := s.store.Get(ctx, strings.TrimSpace(id))
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
	normalizeWalletChallengeLoaded(challenge)
	return challenge, nil
}

func (s *WalletChallengeService) MarkUsed(ctx context.Context, id string, usedAt time.Time) (*authdomain.WalletChallenge, error) {
	if s == nil || s.store == nil {
		return nil, ErrChallengeStore
	}
	challenge, err := s.store.Get(ctx, strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrWalletChallengeNotFound
	}
	if err := s.store.MarkUsed(ctx, strings.TrimSpace(id), usedAt.UTC()); err != nil {
		return nil, err
	}
	challenge.UsedAt = ptrTime(usedAt.UTC())
	normalizeWalletChallengeLoaded(challenge)
	return challenge, nil
}

func NewWalletLinkingService(challenges *WalletChallengeService, identities authdomain.WalletIdentityStore) *WalletLinkingService {
	if challenges == nil {
		return &WalletLinkingService{identities: identities}
	}
	return &WalletLinkingService{challenges: *challenges, identities: identities}
}

func (s *WalletLinkingService) CreateChallenge(ctx context.Context, userID, address, chain string) (*authdomain.WalletChallenge, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" || s == nil {
		return nil, ErrUnauthorized
	}
	return s.challenges.CreateWithOptions(ctx, address, chain, authdomain.WalletChallengeOptions{
		Purpose: WalletChallengePurposeLinkWallet, RequestedByUserID: userID,
	})
}

func (s *WalletLinkingService) VerifyAndLink(ctx context.Context, userID, challengeID, address, signature string) (*WalletLinkResult, error) {
	userID = strings.TrimSpace(userID)
	if s == nil || s.identities == nil || userID == "" {
		return nil, ErrUnauthorized
	}
	address = normalizeWalletAddress(address)
	if !walletEVMAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}
	challenge, err := s.challenges.Get(ctx, strings.TrimSpace(challengeID))
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrWalletChallengeNotFound
	}
	if purpose, ok := canonicalWalletChallengePurpose(challenge.Purpose); !ok || purpose != WalletChallengePurposeLinkWallet {
		return nil, ErrWalletChallengePurpose
	}
	if strings.TrimSpace(challenge.RequestedByUserID) != userID {
		return nil, ErrWalletLinkChallengeMismatch
	}
	if normalizeWalletAddress(challenge.Address) != address {
		return nil, ErrInvalidWalletSignature
	}
	recoveredAddress, err := recoverWalletAddress(challenge.Message, signature)
	if err != nil {
		return nil, err
	}
	if normalizeWalletAddress(recoveredAddress) != address {
		return nil, ErrInvalidWalletSignature
	}

	identity, err := s.identities.GetOrCreate(ctx, address)
	if err != nil {
		return nil, err
	}
	if identity == nil {
		return nil, ErrWalletIdentityNotFound
	}
	if strings.TrimSpace(identity.UserID) != "" {
		if strings.TrimSpace(identity.UserID) == userID {
			return nil, ErrWalletAlreadyLinkedToUser
		}
		return nil, ErrWalletIdentityAlreadyLinked
	}
	challenge, err = s.challenges.MarkUsed(ctx, challenge.ID, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	wallets, err := s.identities.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	attachPrimary := len(wallets) == 0
	linked, err := s.identities.AttachUser(ctx, identity.ID, userID, attachPrimary)
	if err != nil {
		return nil, err
	}
	wallets, err = s.identities.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallets == nil {
		wallets = []*authdomain.WalletIdentity{}
	}

	return &WalletLinkResult{Challenge: challenge, Linked: linked, Wallets: wallets}, nil
}

func NewWalletAccountMergeService(challenges *WalletChallengeService, identities authdomain.WalletIdentityStore) *WalletAccountMergeService {
	if challenges == nil {
		return &WalletAccountMergeService{identities: identities}
	}
	return &WalletAccountMergeService{challenges: *challenges, identities: identities}
}

func (s *WalletAccountMergeService) CreateChallenge(ctx context.Context, userID, address, chain string) (*authdomain.WalletChallenge, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" || s == nil {
		return nil, ErrUnauthorized
	}
	return s.challenges.CreateWithOptions(ctx, address, chain, authdomain.WalletChallengeOptions{
		Purpose: WalletChallengePurposeAccountMerge, RequestedByUserID: userID,
	})
}

func (s *WalletAccountMergeService) VerifyAndMerge(ctx context.Context, targetUserID, challengeID, address, signature string) (*WalletAccountMergeResult, error) {
	targetUserID = strings.TrimSpace(targetUserID)
	if s == nil || s.identities == nil || targetUserID == "" {
		return nil, ErrUnauthorized
	}
	address = normalizeWalletAddress(address)
	if !walletEVMAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}
	challenge, err := s.challenges.Get(ctx, strings.TrimSpace(challengeID))
	if err != nil {
		return nil, err
	}
	if challenge == nil {
		return nil, ErrWalletChallengeNotFound
	}
	if purpose, ok := canonicalWalletChallengePurpose(challenge.Purpose); !ok || purpose != WalletChallengePurposeAccountMerge {
		return nil, ErrWalletChallengePurpose
	}
	if strings.TrimSpace(challenge.RequestedByUserID) != targetUserID {
		return nil, ErrWalletLinkChallengeMismatch
	}
	if normalizeWalletAddress(challenge.Address) != address {
		return nil, ErrInvalidWalletSignature
	}
	recoveredAddress, err := recoverWalletAddress(challenge.Message, signature)
	if err != nil {
		return nil, err
	}
	if normalizeWalletAddress(recoveredAddress) != address {
		return nil, ErrInvalidWalletSignature
	}

	identity, err := s.identities.GetByAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	if identity == nil || strings.TrimSpace(identity.UserID) == "" {
		return nil, ErrWalletMergeSourceNotLinked
	}
	sourceUserID := strings.TrimSpace(identity.UserID)
	if sourceUserID == targetUserID {
		return nil, ErrWalletMergeSameUser
	}

	sourceWallets, err := s.identities.ListByUser(ctx, sourceUserID)
	if err != nil {
		return nil, err
	}
	if len(sourceWallets) == 0 {
		return nil, ErrWalletMergeSourceNotLinked
	}

	challenge, err = s.challenges.MarkUsed(ctx, challenge.ID, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	wallets, err := s.identities.MergeUsers(ctx, sourceUserID, targetUserID)
	if err != nil {
		return nil, err
	}
	if wallets == nil {
		wallets = []*authdomain.WalletIdentity{}
	}

	var mergedWallet *authdomain.WalletIdentity
	for _, wallet := range wallets {
		if wallet != nil && strings.EqualFold(wallet.Address, address) {
			mergedWallet = wallet
			break
		}
	}
	return &WalletAccountMergeResult{Challenge: challenge, MergedWallet: mergedWallet, Wallets: wallets, SourceUserID: sourceUserID, TargetUserID: targetUserID}, nil
}

func NewWalletPrimaryService(identities authdomain.WalletIdentityStore) *WalletPrimaryService {
	return &WalletPrimaryService{identities: identities}
}

func (s *WalletPrimaryService) SetPrimary(ctx context.Context, userID, address string) (*WalletPrimaryResult, error) {
	userID = strings.TrimSpace(userID)
	if s == nil || s.identities == nil || userID == "" {
		return nil, ErrUnauthorized
	}
	address = normalizeWalletAddress(address)
	if !walletEVMAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}
	primary, err := s.identities.SetPrimary(ctx, userID, address)
	if err != nil {
		return nil, err
	}
	wallets, err := s.identities.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallets == nil {
		wallets = []*authdomain.WalletIdentity{}
	}
	return &WalletPrimaryResult{Primary: primary, Wallets: wallets}, nil
}

func NewWalletDetachService(identities authdomain.WalletIdentityStore) *WalletDetachService {
	return &WalletDetachService{identities: identities}
}

func (s *WalletDetachService) CheckEligibility(ctx context.Context, userID, address string) (*WalletDetachCheckResult, error) {
	userID = strings.TrimSpace(userID)
	if s == nil || s.identities == nil || userID == "" {
		return nil, ErrUnauthorized
	}
	address = normalizeWalletAddress(address)
	if !walletEVMAddressRE.MatchString(address) {
		return nil, ErrInvalidWalletAddress
	}
	identity, err := s.identities.GetByAddress(ctx, address)
	if err != nil {
		return nil, err
	}
	if identity == nil || strings.TrimSpace(identity.UserID) != userID {
		return nil, ErrWalletNotOwnedByUser
	}
	wallets, err := s.identities.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallets == nil {
		wallets = []*authdomain.WalletIdentity{}
	}

	result := &WalletDetachCheckResult{WalletAddress: address, Eligible: true, IsPrimary: identity.IsPrimary, OwnedWalletCount: len(wallets), Reasons: []string{}}
	if identity.IsPrimary {
		result.Reasons = append(result.Reasons, authdomain.WalletDetachReasonWalletIsPrimary)
	}
	if len(wallets) <= 1 {
		result.Reasons = append(result.Reasons, authdomain.WalletDetachReasonUserWouldBeEmpty)
	}
	if len(result.Reasons) > 0 {
		result.Eligible = false
	}
	return result, nil
}

func (s *WalletDetachService) Execute(ctx context.Context, userID, address string) (*WalletDetachExecuteResult, error) {
	check, err := s.CheckEligibility(ctx, userID, address)
	if err != nil {
		return nil, err
	}
	if check == nil {
		return nil, ErrUnauthorized
	}
	if !check.Eligible {
		return &WalletDetachExecuteResult{Check: check, Wallets: []*authdomain.WalletIdentity{}}, ErrWalletDetachNotEligible
	}
	detached, wallets, err := s.identities.DetachUser(ctx, userID, address)
	if err != nil {
		return nil, err
	}
	if wallets == nil {
		wallets = []*authdomain.WalletIdentity{}
	}
	return &WalletDetachExecuteResult{Detached: detached, Wallets: wallets, Check: check}, nil
}

func canonicalWalletChallengePurpose(purpose string) (string, bool) {
	purpose = strings.TrimSpace(strings.ToLower(purpose))
	switch purpose {
	case WalletChallengePurposeAuthBootstrap, "auth", "auth_login", "bootstrap":
		return WalletChallengePurposeAuthBootstrap, true
	case WalletChallengePurposeLinkWallet, "link", "wallet-link":
		return WalletChallengePurposeLinkWallet, true
	case WalletChallengePurposeAccountMerge, "merge", "account-merge":
		return WalletChallengePurposeAccountMerge, true
	default:
		return "", false
	}
}

func resolveWalletChallengePurposeForCreate(purpose string) (string, error) {
	purpose = strings.TrimSpace(strings.ToLower(purpose))
	if purpose == "" {
		return WalletChallengePurposeAuthBootstrap, nil
	}
	if canonical, ok := canonicalWalletChallengePurpose(purpose); ok {
		return canonical, nil
	}
	return "", ErrWalletChallengePurpose
}

func normalizeWalletChallengeLoaded(ch *authdomain.WalletChallenge) {
	if ch == nil {
		return
	}
	ch.Address = normalizeWalletAddress(ch.Address)
	ch.Chain = normalizeChain(ch.Chain)
	ch.Purpose = normalizeStoredWalletChallengePurpose(ch.Purpose)
	ch.RequestedByUserID = strings.TrimSpace(ch.RequestedByUserID)
	ch.Message = strings.TrimSpace(ch.Message)
}

func normalizeStoredWalletChallengePurpose(purpose string) string {
	if canonical, ok := canonicalWalletChallengePurpose(purpose); ok {
		return canonical
	}
	return strings.TrimSpace(strings.ToLower(purpose))
}

func randomToken(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("invalid token length")
	}
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	token := base64.RawURLEncoding.EncodeToString(buf)
	if len(token) > length {
		token = token[:length]
	}
	return token, nil
}

func ptrTime(v time.Time) *time.Time { return &v }

func buildChallengeMessage(publicBaseURL string, ch *authdomain.WalletChallenge) string {
	domain := "SCAVO Exchange"
	uri := "http://localhost"
	if publicBaseURL != "" {
		if u, err := url.Parse(publicBaseURL); err == nil {
			if host := strings.TrimSpace(u.Host); host != "" {
				domain = host
			}
			uri = publicBaseURL
		}
	}
	purposeLine := "Purpose: SCAVO Exchange wallet authentication bootstrap."
	if purpose, ok := canonicalWalletChallengePurpose(ch.Purpose); ok {
		switch purpose {
		case WalletChallengePurposeLinkWallet:
			purposeLine = "Purpose: SCAVO Exchange authenticated wallet linking confirmation."
		case WalletChallengePurposeAccountMerge:
			purposeLine = "Purpose: SCAVO Exchange authenticated account merge confirmation."
		}
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
