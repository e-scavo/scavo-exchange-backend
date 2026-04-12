package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"
)

var (
	ErrApplicationNotConfigured = errors.New("application not configured")
	ErrWalletIdentityStore      = errors.New("wallet identity store error")
)

type Application struct {
	Tokens           *coreauth.TokenService
	TTL              time.Duration
	Users            *usermod.Service
	UserSettings     *usersettingsmod.Service
	PublicBaseURL    string
	ChallengeTTL     time.Duration
	Challenges       WalletChallengeStore
	WalletIdentities WalletIdentityStore
}

func NewApplication(tokens *coreauth.TokenService, ttl time.Duration, users *usermod.Service, userSettings *usersettingsmod.Service, publicBaseURL string, challengeTTL time.Duration, challenges WalletChallengeStore, walletIdentities WalletIdentityStore) *Application {
	return &Application{
		Tokens:           tokens,
		TTL:              ttl,
		Users:            users,
		UserSettings:     userSettings,
		PublicBaseURL:    publicBaseURL,
		ChallengeTTL:     challengeTTL,
		Challenges:       challenges,
		WalletIdentities: walletIdentities,
	}
}

func (h HTTPHandlers) Application() *Application {
	return NewApplication(h.Tokens, h.TTL, h.Users, h.UserSettings, h.PublicBaseURL, h.ChallengeTTL, h.Challenges, h.WalletIdentities)
}

func (a *Application) challengeTTL() time.Duration {
	if a == nil || a.ChallengeTTL <= 0 {
		return 5 * time.Minute
	}
	return a.ChallengeTTL
}

func (a *Application) challengeService() WalletChallengeService {
	return *NewWalletChallengeService(a.Challenges, a.PublicBaseURL, a.challengeTTL())
}

func (a *Application) Login(ctx context.Context, email, password string) (LoginResponse, error) {
	if a == nil || a.Tokens == nil {
		return LoginResponse{}, ErrApplicationNotConfigured
	}

	svc := NewService(a.Tokens, a.Users, a.TTL)
	result, err := svc.LoginDev(ctx, email, password)
	if err != nil {
		return LoginResponse{}, err
	}

	userID := ""
	if result.User != nil {
		userID = result.User.ID
	}

	return LoginResponse{
		AccessToken: result.AccessToken,
		TokenType:   result.TokenType,
		ExpiresIn:   result.ExpiresIn,
		UserID:      userID,
	}, nil
}

func (a *Application) GetMe(ctx context.Context, claims *coreauth.Claims) (MeResponse, error) {
	if claims == nil {
		return MeResponse{}, ErrUnauthorized
	}
	if a == nil {
		return MeResponse{}, ErrApplicationNotConfigured
	}

	profile, err := buildProfileView(ctx, claims, a.Users, a.WalletIdentities)
	if err != nil {
		return MeResponse{}, err
	}

	return MeResponse{
		User:    profile.User,
		Profile: profile,
	}, nil
}

func (a *Application) GetSession(ctx context.Context, claims *coreauth.Claims) (SessionResponse, error) {
	if claims == nil {
		return SessionResponse{}, ErrUnauthorized
	}
	if a == nil {
		return SessionResponse{}, ErrApplicationNotConfigured
	}

	svc := NewService(a.Tokens, a.Users, a.TTL)
	session, err := svc.ResolveSessionClaims(ctx, claims)
	if err != nil {
		return SessionResponse{}, err
	}

	return SessionResponse{Session: session}, nil
}

func (a *Application) GetBootstrap(ctx context.Context, claims *coreauth.Claims) (BootstrapResponse, error) {
	if claims == nil {
		return BootstrapResponse{}, ErrUnauthorized
	}
	if a == nil || a.UserSettings == nil {
		return BootstrapResponse{}, ErrApplicationNotConfigured
	}

	svc := NewService(a.Tokens, a.Users, a.TTL)
	user, err := svc.ResolveCurrentUserClaims(ctx, claims)
	if err != nil {
		return BootstrapResponse{}, err
	}

	session := buildSessionViewWithUser(claims, user)
	profile, err := buildProfileViewWithUser(ctx, claims, user, a.WalletIdentities)
	if err != nil {
		return BootstrapResponse{}, err
	}

	settings, err := a.UserSettings.GetOrDefault(ctx, session.UserID)
	if err != nil {
		return BootstrapResponse{}, err
	}

	wallets, err := listWalletReadModelsForUser(ctx, session.UserID, a.WalletIdentities)
	if err != nil {
		return BootstrapResponse{}, fmt.Errorf("%w: %v", ErrWalletIdentityStore, err)
	}

	return BootstrapResponse{
		Session:  session,
		User:     user,
		Profile:  profile,
		Settings: usersettingsmod.ToView(settings),
		Wallets:  buildBootstrapWalletsView(wallets),
	}, nil
}

func (a *Application) ListWallets(ctx context.Context, userID string, query WalletsQuery) (WalletsResponse, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return WalletsResponse{}, ErrUnauthorized
	}
	if a == nil {
		return WalletsResponse{}, ErrApplicationNotConfigured
	}
	if a.WalletIdentities == nil {
		return buildWalletsResponse([]*WalletReadModel{}, 0, query), nil
	}

	wallets, err := a.WalletIdentities.ListByUser(ctx, userID)
	if err != nil {
		return WalletsResponse{}, fmt.Errorf("%w: %v", ErrWalletIdentityStore, err)
	}
	if wallets == nil {
		wallets = []*WalletIdentity{}
	}

	mapped := mapWalletIdentitiesToReadModels(wallets)
	window, total := applyWalletsQuery(mapped, query)
	return buildWalletsResponse(window, total, query), nil
}

func (a *Application) CreateWalletLinkChallenge(ctx context.Context, userID, address, chain string) (WalletLinkChallengeResponse, error) {
	if a == nil {
		return WalletLinkChallengeResponse{}, ErrApplicationNotConfigured
	}
	linkSvc := NewWalletLinkingService(ptrWalletChallengeService(a.challengeService()), a.WalletIdentities)
	challenge, err := linkSvc.CreateChallenge(ctx, userID, address, chain)
	if err != nil {
		return WalletLinkChallengeResponse{}, err
	}
	return WalletLinkChallengeResponse{Challenge: challenge}, nil
}

func (a *Application) VerifyWalletLink(ctx context.Context, userID, challengeID, address, signature string) (WalletLinkVerifyResponse, error) {
	if a == nil {
		return WalletLinkVerifyResponse{}, ErrApplicationNotConfigured
	}
	linkSvc := NewWalletLinkingService(ptrWalletChallengeService(a.challengeService()), a.WalletIdentities)
	result, err := linkSvc.VerifyAndLink(ctx, userID, challengeID, address, signature)
	if err != nil {
		return WalletLinkVerifyResponse{}, err
	}
	return WalletLinkVerifyResponse{
		LinkedWallet: result.Linked,
		Wallets:      result.Wallets,
		Challenge:    result.Challenge,
	}, nil
}

func (a *Application) CreateWalletAccountMergeChallenge(ctx context.Context, userID, address, chain string) (WalletAccountMergeChallengeResponse, error) {
	if a == nil {
		return WalletAccountMergeChallengeResponse{}, ErrApplicationNotConfigured
	}
	mergeSvc := NewWalletAccountMergeService(ptrWalletChallengeService(a.challengeService()), a.WalletIdentities)
	challenge, err := mergeSvc.CreateChallenge(ctx, userID, address, chain)
	if err != nil {
		return WalletAccountMergeChallengeResponse{}, err
	}
	return WalletAccountMergeChallengeResponse{Challenge: challenge}, nil
}

func (a *Application) VerifyWalletAccountMerge(ctx context.Context, userID, challengeID, address, signature string) (WalletAccountMergeVerifyResponse, error) {
	if a == nil {
		return WalletAccountMergeVerifyResponse{}, ErrApplicationNotConfigured
	}
	mergeSvc := NewWalletAccountMergeService(ptrWalletChallengeService(a.challengeService()), a.WalletIdentities)
	result, err := mergeSvc.VerifyAndMerge(ctx, userID, challengeID, address, signature)
	if err != nil {
		return WalletAccountMergeVerifyResponse{}, err
	}
	return WalletAccountMergeVerifyResponse{
		MergedWallet: result.MergedWallet,
		Wallets:      result.Wallets,
		Challenge:    result.Challenge,
		SourceUserID: result.SourceUserID,
		TargetUserID: result.TargetUserID,
	}, nil
}

func (a *Application) SetPrimaryWallet(ctx context.Context, userID, address string) (WalletPrimarySetResponse, error) {
	if a == nil {
		return WalletPrimarySetResponse{}, ErrApplicationNotConfigured
	}
	svc := NewWalletPrimaryService(a.WalletIdentities)
	result, err := svc.SetPrimary(ctx, userID, address)
	if err != nil {
		return WalletPrimarySetResponse{}, err
	}
	return WalletPrimarySetResponse{
		PrimaryWallet: result.Primary,
		Wallets:       result.Wallets,
	}, nil
}

func (a *Application) CheckWalletDetach(ctx context.Context, userID, address string) (WalletDetachCheckResponse, error) {
	if a == nil {
		return WalletDetachCheckResponse{}, ErrApplicationNotConfigured
	}
	svc := NewWalletDetachService(a.WalletIdentities)
	result, err := svc.CheckEligibility(ctx, userID, address)
	if err != nil {
		return WalletDetachCheckResponse{}, err
	}
	return WalletDetachCheckResponse{
		WalletAddress:    result.WalletAddress,
		Eligible:         result.Eligible,
		IsPrimary:        result.IsPrimary,
		OwnedWalletCount: result.OwnedWalletCount,
		Reasons:          result.Reasons,
	}, nil
}

func (a *Application) ExecuteWalletDetach(ctx context.Context, userID, address string) (WalletDetachExecuteResponse, error) {
	if a == nil {
		return WalletDetachExecuteResponse{}, ErrApplicationNotConfigured
	}
	svc := NewWalletDetachService(a.WalletIdentities)
	result, err := svc.Execute(ctx, userID, address)
	response := WalletDetachExecuteResponse{}
	if result != nil {
		response.DetachedWallet = result.Detached
		response.Wallets = result.Wallets
		if result.Check != nil {
			response.Check = &WalletDetachCheckResponse{
				WalletAddress:    result.Check.WalletAddress,
				Eligible:         result.Check.Eligible,
				IsPrimary:        result.Check.IsPrimary,
				OwnedWalletCount: result.Check.OwnedWalletCount,
				Reasons:          result.Check.Reasons,
			}
		}
	}
	return response, err
}

func ptrWalletChallengeService(s WalletChallengeService) *WalletChallengeService {
	return &s
}
