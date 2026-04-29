package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	authdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/domain"
	authmappers "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/mappers"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	usermappers "github.com/e-scavo/scavo-exchange-backend/internal/modules/user/mappers"
	usersettingsmappers "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings/mappers"
)

var (
	ErrApplicationNotConfigured = errors.New("application not configured")
	ErrWalletIdentityStore      = errors.New("wallet identity store error")
)

type Application struct {
	Tokens           *coreauth.TokenService
	TTL              time.Duration
	Users            authdomain.UserProvider
	UserSettings     authdomain.UserSettingsProvider
	PublicBaseURL    string
	ChallengeTTL     time.Duration
	Challenges       authdomain.WalletChallengeStore
	WalletIdentities authdomain.WalletIdentityStore
}

func NewApplication(tokens *coreauth.TokenService, ttl time.Duration, users authdomain.UserProvider, userSettings authdomain.UserSettingsProvider, publicBaseURL string, challengeTTL time.Duration, challenges authdomain.WalletChallengeStore, walletIdentities authdomain.WalletIdentityStore) *Application {
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
	return authmappers.NewAuthLoginReadModel(
		result.AccessToken,
		result.TokenType,
		result.ExpiresIn,
		userID,
	), nil
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

func (a *Application) UpdateProfile(ctx context.Context, claims *coreauth.Claims, input authdomain.ProfileUpdateInput) (MeResponse, error) {
	if claims == nil {
		return MeResponse{}, ErrUnauthorized
	}
	if a == nil || a.Users == nil {
		return MeResponse{}, ErrApplicationNotConfigured
	}

	updatedUser, err := a.Users.UpdateDisplayName(ctx, claims.UserID, input.DisplayName)
	if err != nil {
		return MeResponse{}, err
	}

	profile, err := buildProfileViewWithUser(ctx, claims, updatedUser, a.WalletIdentities)
	if err != nil {
		return MeResponse{}, err
	}

	return MeResponse{
		User:    profile.User,
		Profile: profile,
	}, nil
}

func (a *Application) GetSettings(ctx context.Context, claims *coreauth.Claims) (MeSettingsResponse, error) {
	if claims == nil {
		return MeSettingsResponse{}, ErrUnauthorized
	}
	if a == nil || a.UserSettings == nil {
		return MeSettingsResponse{}, ErrApplicationNotConfigured
	}

	settings, err := a.UserSettings.GetOrDefault(ctx, claims.UserID)
	if err != nil {
		return MeSettingsResponse{}, err
	}

	return MeSettingsResponse{
		Settings: usersettingsmappers.UserSettingsToReadModel(settings),
	}, nil
}

func (a *Application) UpdateSettings(ctx context.Context, claims *coreauth.Claims, input authdomain.SettingsUpdateInput) (MeSettingsResponse, error) {
	if claims == nil {
		return MeSettingsResponse{}, ErrUnauthorized
	}
	if a == nil || a.UserSettings == nil {
		return MeSettingsResponse{}, ErrApplicationNotConfigured
	}

	settings, err := a.UserSettings.UpdatePreferences(ctx, claims.UserID, input.Preferences)
	if err != nil {
		return MeSettingsResponse{}, err
	}

	return MeSettingsResponse{
		Settings: usersettingsmappers.UserSettingsToReadModel(settings),
	}, nil
}

func (a *Application) CreateWalletChallenge(ctx context.Context, address, chain string) (WalletChallengeResponse, error) {
	if a == nil {
		return WalletChallengeResponse{}, ErrApplicationNotConfigured
	}

	challenge, err := NewWalletChallengeService(a.Challenges, a.PublicBaseURL, a.challengeTTL()).Create(ctx, address, chain)
	if err != nil {
		return WalletChallengeResponse{}, err
	}

	return WalletChallengeResponse{
		Challenge: authmappers.WalletChallengeToReadModel(challenge),
	}, nil
}

func (a *Application) VerifyWallet(ctx context.Context, challengeID, address, signature string) (WalletVerifyResponse, error) {
	if a == nil || a.Tokens == nil || a.WalletIdentities == nil {
		return WalletVerifyResponse{}, ErrUnauthorized
	}

	address = normalizeWalletAddress(address)
	if !walletEVMAddressRE.MatchString(address) {
		return WalletVerifyResponse{}, ErrInvalidWalletAddress
	}

	challengeSvc := NewWalletChallengeService(a.Challenges, a.PublicBaseURL, a.challengeTTL())
	challenge, err := challengeSvc.Get(ctx, strings.TrimSpace(challengeID))
	if err != nil {
		return WalletVerifyResponse{}, err
	}
	if challenge == nil {
		return WalletVerifyResponse{}, ErrWalletChallengeNotFound
	}
	if purpose, ok := canonicalWalletChallengePurpose(challenge.Purpose); !ok || purpose != WalletChallengePurposeAuthBootstrap {
		return WalletVerifyResponse{}, ErrWalletChallengePurpose
	}
	if normalizeWalletAddress(challenge.Address) != address {
		return WalletVerifyResponse{}, ErrInvalidWalletSignature
	}

	recoveredAddress, err := recoverWalletAddress(challenge.Message, signature)
	if err != nil {
		return WalletVerifyResponse{}, err
	}
	if normalizeWalletAddress(recoveredAddress) != address {
		return WalletVerifyResponse{}, ErrInvalidWalletSignature
	}

	usedAt := time.Now().UTC()
	challenge, err = challengeSvc.MarkUsed(ctx, challenge.ID, usedAt)
	if err != nil {
		return WalletVerifyResponse{}, err
	}

	identity, err := a.WalletIdentities.GetOrCreate(ctx, address)
	if err != nil {
		return WalletVerifyResponse{}, err
	}

	user, identity, err := a.resolveWalletBootstrapUser(ctx, identity, address)
	if err != nil {
		return WalletVerifyResponse{}, err
	}

	loginSvc := NewService(a.Tokens, a.Users, a.TTL)
	result, err := loginSvc.LoginWalletForUser(ctx, user, identity.ID, address, challenge.Chain)
	if err != nil {
		return WalletVerifyResponse{}, err
	}

	userID := ""
	if result.User != nil {
		userID = result.User.ID
	}

	return WalletVerifyResponse{
		AccessToken:   result.AccessToken,
		TokenType:     result.TokenType,
		ExpiresIn:     result.ExpiresIn,
		UserID:        userID,
		WalletID:      result.WalletID,
		WalletAddress: result.WalletAddress,
		Chain:         result.Chain,
		AuthMethod:    result.AuthMethod,
		User:          usermappers.UserToReadModel(result.User),
		Challenge:     authmappers.WalletChallengeToReadModel(challenge),
	}, nil
}

func (a *Application) resolveWalletBootstrapUser(ctx context.Context, identity *authdomain.WalletIdentity, address string) (*usermod.User, *authdomain.WalletIdentity, error) {
	if identity == nil {
		return nil, nil, ErrUnauthorized
	}

	if strings.TrimSpace(identity.UserID) != "" {
		if a.Users == nil {
			return walletUser(address), identity, nil
		}

		user, err := a.Users.GetByID(ctx, identity.UserID, walletUserEmail(address))
		if err == nil {
			return user, identity, nil
		}
		if !errors.Is(err, usermod.ErrUserNotFound) {
			return nil, nil, err
		}
	}

	if a.Users == nil {
		linked := walletUser(address)
		identity.UserID = linked.ID
		now := time.Now().UTC()
		identity.LinkedAt = &now
		identity.IsPrimary = true
		return linked, identity, nil
	}

	user, err := a.Users.ResolveOrCreateWalletUser(ctx, address)
	if err != nil {
		return nil, nil, err
	}

	identity, err = a.WalletIdentities.AttachUser(ctx, identity.ID, user.ID, true)
	if err != nil {
		return nil, nil, err
	}

	return user, identity, nil
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
		User:     usermappers.UserToReadModel(user),
		Profile:  profile,
		Settings: usersettingsmappers.UserSettingsToReadModel(settings),
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
		wallets = []*authdomain.WalletIdentity{}
	}

	mapped := authmappers.WalletIdentitiesToActionableReadModels(wallets)
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
	return WalletLinkChallengeResponse{Challenge: authmappers.WalletChallengeToReadModel(challenge)}, nil
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
		LinkedWallet: authmappers.WalletIdentityToReadModel(result.Linked),
		Wallets:      authmappers.WalletIdentitiesToActionableReadModels(result.Wallets),
		Challenge:    authmappers.WalletChallengeToReadModel(result.Challenge),
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
	return WalletAccountMergeChallengeResponse{Challenge: authmappers.WalletChallengeToReadModel(challenge)}, nil
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
		MergedWallet: authmappers.WalletIdentityToReadModel(result.MergedWallet),
		Wallets:      authmappers.WalletIdentitiesToActionableReadModels(result.Wallets),
		Challenge:    authmappers.WalletChallengeToReadModel(result.Challenge),
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
		PrimaryWallet: authmappers.WalletIdentityToReadModel(result.Primary),
		Wallets:       authmappers.WalletIdentitiesToActionableReadModels(result.Wallets),
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
		response.DetachedWallet = authmappers.WalletIdentityToReadModel(result.Detached)
		response.Wallets = authmappers.WalletIdentitiesToActionableReadModels(result.Wallets)
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
