package auth

import (
	"context"
	"errors"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	authapp "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/app"
	authdomain "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/domain"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"
)

var (
	ErrApplicationNotConfigured = authapp.ErrApplicationNotConfigured
	ErrWalletIdentityStore      = authapp.ErrWalletIdentityStore
)

type Application struct {
	inner *authapp.Application
}

// Compile-time contract assertions keep the provider boundary explicit. They
// fail compilation if HTTP handlers can no longer consume Application through
// AuthProvider or if the cross-module providers used by auth drift from the
// minimal domain contracts.
var (
	_ AuthProvider                    = (*Application)(nil)
	_ authdomain.UserProvider         = (*usermod.Service)(nil)
	_ authdomain.UserSettingsProvider = (*usersettingsmod.Service)(nil)
)

func NewApplication(
	tokens *coreauth.TokenService,
	ttl time.Duration,
	users *usermod.Service,
	userSettings *usersettingsmod.Service,
	publicBaseURL string,
	challengeTTL time.Duration,
	challenges WalletChallengeStore,
	walletIdentities WalletIdentityStore,
) *Application {
	var userSettingsProvider authdomain.UserSettingsProvider
	if userSettings != nil {
		userSettingsProvider = userSettings
	}

	return &Application{
		inner: authapp.NewApplication(
			tokens,
			ttl,
			users,
			userSettingsProvider,
			publicBaseURL,
			challengeTTL,
			newWalletChallengeStoreAdapter(challenges, publicBaseURL, challengeTTL),
			walletIdentities,
		),
	}
}

func (h HTTPHandlers) AuthProvider() AuthProvider {
	if h.Provider != nil {
		return h.Provider
	}
	return h.Application()
}

func (h HTTPHandlers) Application() *Application {
	return NewApplication(
		h.Tokens,
		h.TTL,
		h.Users,
		h.UserSettings,
		h.PublicBaseURL,
		h.ChallengeTTL,
		h.Challenges,
		h.WalletIdentities,
	)
}

func (a *Application) Login(ctx context.Context, email, password string) (LoginResponse, error) {
	response, err := a.inner.Login(ctx, email, password)
	if err != nil {
		return LoginResponse{}, normalizeApplicationError(err)
	}

	return LoginResponse{
		AccessToken: response.AccessToken,
		TokenType:   response.TokenType,
		ExpiresIn:   response.ExpiresIn,
		UserID:      response.UserID,
	}, nil
}

func (a *Application) GetMe(ctx context.Context, claims *coreauth.Claims) (MeResponse, error) {
	response, err := a.inner.GetMe(ctx, claims)
	if err != nil {
		return MeResponse{}, normalizeApplicationError(err)
	}

	return MeResponse{
		User:    response.User,
		Profile: mapProfileViewFromApp(response.Profile),
	}, nil
}

func (a *Application) GetSession(ctx context.Context, claims *coreauth.Claims) (SessionResponse, error) {
	response, err := a.inner.GetSession(ctx, claims)
	if err != nil {
		return SessionResponse{}, normalizeApplicationError(err)
	}

	return SessionResponse{
		Session: mapSessionViewFromApp(response.Session),
	}, nil
}

func (a *Application) UpdateProfile(ctx context.Context, claims *coreauth.Claims, input authdomain.ProfileUpdateInput) (MeResponse, error) {
	response, err := a.inner.UpdateProfile(ctx, claims, input)
	if err != nil {
		return MeResponse{}, normalizeApplicationError(err)
	}

	return MeResponse{
		User:    response.User,
		Profile: mapProfileViewFromApp(response.Profile),
	}, nil
}

func (a *Application) GetSettings(ctx context.Context, claims *coreauth.Claims) (MeSettingsResponse, error) {
	response, err := a.inner.GetSettings(ctx, claims)
	if err != nil {
		return MeSettingsResponse{}, normalizeApplicationError(err)
	}

	return MeSettingsResponse{Settings: response.Settings}, nil
}

func (a *Application) UpdateSettings(ctx context.Context, claims *coreauth.Claims, input authdomain.SettingsUpdateInput) (MeSettingsResponse, error) {
	response, err := a.inner.UpdateSettings(ctx, claims, input)
	if err != nil {
		return MeSettingsResponse{}, normalizeApplicationError(err)
	}

	return MeSettingsResponse{Settings: response.Settings}, nil
}

func (a *Application) CreateWalletChallenge(ctx context.Context, address, chain string) (WalletChallengeResponse, error) {
	response, err := a.inner.CreateWalletChallenge(ctx, address, chain)
	if err != nil {
		return WalletChallengeResponse{}, normalizeApplicationError(err)
	}

	return WalletChallengeResponse{Challenge: response.Challenge}, nil
}

func (a *Application) VerifyWallet(ctx context.Context, challengeID, address, signature string) (WalletVerifyResponse, error) {
	response, err := a.inner.VerifyWallet(ctx, challengeID, address, signature)
	if err != nil {
		return WalletVerifyResponse{}, normalizeApplicationError(err)
	}

	return WalletVerifyResponse{
		AccessToken:   response.AccessToken,
		TokenType:     response.TokenType,
		ExpiresIn:     response.ExpiresIn,
		UserID:        response.UserID,
		WalletID:      response.WalletID,
		WalletAddress: response.WalletAddress,
		Chain:         response.Chain,
		AuthMethod:    response.AuthMethod,
		User:          response.User,
		Challenge:     response.Challenge,
	}, nil
}

func (a *Application) GetBootstrap(ctx context.Context, claims *coreauth.Claims) (BootstrapResponse, error) {
	response, err := a.inner.GetBootstrap(ctx, claims)
	if err != nil {
		return BootstrapResponse{}, normalizeApplicationError(err)
	}

	return BootstrapResponse{
		Session:  mapSessionViewFromApp(response.Session),
		User:     response.User,
		Profile:  mapProfileViewFromApp(response.Profile),
		Settings: response.Settings,
		Wallets: BootstrapWalletsView{
			Items: mapWalletReadModelsFromApp(response.Wallets.Items),
			Total: response.Wallets.Total,
		},
	}, nil
}

func (a *Application) ListWallets(ctx context.Context, userID string, query WalletsQuery) (WalletsResponse, error) {
	response, err := a.inner.ListWallets(ctx, userID, authapp.WalletsQuery(query))
	if err != nil {
		return WalletsResponse{}, normalizeApplicationError(err)
	}

	return WalletsResponse{
		Items:          mapWalletReadModelsFromApp(response.Items),
		Wallets:        mapWalletReadModelsFromApp(response.Wallets),
		Total:          response.Total,
		Limit:          response.Limit,
		Offset:         response.Offset,
		Returned:       response.Returned,
		HasMore:        response.HasMore,
		NextOffset:     response.NextOffset,
		PreviousOffset: response.PreviousOffset,
	}, nil
}

func (a *Application) CreateWalletLinkChallenge(ctx context.Context, userID, address, chain string) (WalletLinkChallengeResponse, error) {
	response, err := a.inner.CreateWalletLinkChallenge(ctx, userID, address, chain)
	if err != nil {
		return WalletLinkChallengeResponse{}, normalizeApplicationError(err)
	}

	return WalletLinkChallengeResponse{
		Challenge: response.Challenge,
	}, nil
}

func (a *Application) VerifyWalletLink(ctx context.Context, userID, challengeID, address, signature string) (WalletLinkVerifyResponse, error) {
	response, err := a.inner.VerifyWalletLink(ctx, userID, challengeID, address, signature)
	if err != nil {
		return WalletLinkVerifyResponse{}, normalizeApplicationError(err)
	}

	return WalletLinkVerifyResponse{
		LinkedWallet: response.LinkedWallet,
		Wallets:      mapWalletReadModelsFromApp(response.Wallets),
		Challenge:    response.Challenge,
	}, nil
}

func (a *Application) CreateWalletAccountMergeChallenge(ctx context.Context, userID, address, chain string) (WalletAccountMergeChallengeResponse, error) {
	response, err := a.inner.CreateWalletAccountMergeChallenge(ctx, userID, address, chain)
	if err != nil {
		return WalletAccountMergeChallengeResponse{}, normalizeApplicationError(err)
	}

	return WalletAccountMergeChallengeResponse{
		Challenge: response.Challenge,
	}, nil
}

func (a *Application) VerifyWalletAccountMerge(ctx context.Context, userID, challengeID, address, signature string) (WalletAccountMergeVerifyResponse, error) {
	response, err := a.inner.VerifyWalletAccountMerge(ctx, userID, challengeID, address, signature)
	if err != nil {
		return WalletAccountMergeVerifyResponse{}, normalizeApplicationError(err)
	}

	return WalletAccountMergeVerifyResponse{
		MergedWallet: response.MergedWallet,
		Wallets:      mapWalletReadModelsFromApp(response.Wallets),
		Challenge:    response.Challenge,
		SourceUserID: response.SourceUserID,
		TargetUserID: response.TargetUserID,
	}, nil
}

func (a *Application) SetPrimaryWallet(ctx context.Context, userID, address string) (WalletPrimarySetResponse, error) {
	response, err := a.inner.SetPrimaryWallet(ctx, userID, address)
	if err != nil {
		return WalletPrimarySetResponse{}, normalizeApplicationError(err)
	}

	return WalletPrimarySetResponse{
		PrimaryWallet: response.PrimaryWallet,
		Wallets:       mapWalletReadModelsFromApp(response.Wallets),
	}, nil
}

func (a *Application) CheckWalletDetach(ctx context.Context, userID, address string) (WalletDetachCheckResponse, error) {
	response, err := a.inner.CheckWalletDetach(ctx, userID, address)
	if err != nil {
		return WalletDetachCheckResponse{}, normalizeApplicationError(err)
	}

	return WalletDetachCheckResponse{
		WalletAddress:    response.WalletAddress,
		Eligible:         response.Eligible,
		IsPrimary:        response.IsPrimary,
		OwnedWalletCount: response.OwnedWalletCount,
		Reasons:          response.Reasons,
	}, nil
}

func (a *Application) ExecuteWalletDetach(ctx context.Context, userID, address string) (WalletDetachExecuteResponse, error) {
	response, err := a.inner.ExecuteWalletDetach(ctx, userID, address)

	out := WalletDetachExecuteResponse{
		DetachedWallet: response.DetachedWallet,
		Wallets:        mapWalletReadModelsFromApp(response.Wallets),
		Check:          mapWalletDetachCheckFromApp(response.Check),
	}

	if err != nil {
		return out, normalizeApplicationError(err)
	}

	return out, nil
}

type walletChallengeStoreAdapter struct {
	store         WalletChallengeStore
	publicBaseURL string
	ttl           time.Duration
}

func newWalletChallengeStoreAdapter(
	store WalletChallengeStore,
	publicBaseURL string,
	ttl time.Duration,
) authdomain.WalletChallengeStore {
	if store == nil {
		return nil
	}
	return &walletChallengeStoreAdapter{
		store:         store,
		publicBaseURL: publicBaseURL,
		ttl:           ttl,
	}
}

func (a *walletChallengeStoreAdapter) Create(ctx context.Context, address, chain string, ttl time.Duration) (*authdomain.WalletChallenge, error) {
	svc := NewWalletChallengeService(a.store, a.publicBaseURL, ttl)
	challenge, err := svc.Create(ctx, address, chain)
	if err != nil {
		return nil, err
	}
	return mapWalletChallengeToDomain(challenge), nil
}

func (a *walletChallengeStoreAdapter) CreateWithOptions(ctx context.Context, address, chain string, opts authdomain.WalletChallengeOptions) (*authdomain.WalletChallenge, error) {
	svc := NewWalletChallengeService(a.store, a.publicBaseURL, a.ttl)
	challenge, err := svc.CreateWithOptions(ctx, address, chain, WalletChallengeOptions{
		Purpose:           opts.Purpose,
		RequestedByUserID: opts.RequestedByUserID,
	})
	if err != nil {
		return nil, err
	}
	return mapWalletChallengeToDomain(challenge), nil
}

func (a *walletChallengeStoreAdapter) Get(ctx context.Context, challengeID string) (*authdomain.WalletChallenge, error) {
	svc := NewWalletChallengeService(a.store, a.publicBaseURL, a.ttl)
	challenge, err := svc.Get(ctx, challengeID)
	if err != nil {
		return nil, err
	}
	return mapWalletChallengeToDomain(challenge), nil
}

func (a *walletChallengeStoreAdapter) MarkUsed(ctx context.Context, challengeID string, usedAt time.Time) error {
	svc := NewWalletChallengeService(a.store, a.publicBaseURL, a.ttl)
	_, err := svc.MarkUsed(ctx, challengeID, usedAt)
	return err
}

func normalizeApplicationError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, authapp.ErrInvalidCredentials):
		return ErrInvalidCredentials
	case errors.Is(err, authapp.ErrUnauthorized):
		return ErrUnauthorized
	case errors.Is(err, authapp.ErrInvalidWalletAddress):
		return ErrInvalidWalletAddress
	case errors.Is(err, authapp.ErrChallengeStore):
		return ErrChallengeStore
	case errors.Is(err, authapp.ErrChallengeExpired):
		return ErrChallengeExpired
	case errors.Is(err, authapp.ErrChallengeUsed):
		return ErrChallengeUsed
	case errors.Is(err, authapp.ErrWalletChallengeNotFound):
		return ErrWalletChallengeNotFound
	case errors.Is(err, authapp.ErrWalletIdentityNotFound):
		return ErrWalletIdentityNotFound
	case errors.Is(err, authapp.ErrWalletIdentityAlreadyLinked):
		return ErrWalletIdentityAlreadyLinked
	case errors.Is(err, authapp.ErrWalletAlreadyLinkedToUser):
		return ErrWalletAlreadyLinkedToUser
	case errors.Is(err, authapp.ErrWalletLinkChallengeMismatch):
		return ErrWalletLinkChallengeMismatch
	case errors.Is(err, authapp.ErrWalletChallengePurpose):
		return ErrWalletChallengePurpose
	case errors.Is(err, authapp.ErrWalletMergeSourceNotLinked):
		return ErrWalletMergeSourceNotLinked
	case errors.Is(err, authapp.ErrWalletMergeSameUser):
		return ErrWalletMergeSameUser
	case errors.Is(err, authapp.ErrWalletNotOwnedByUser):
		return ErrWalletNotOwnedByUser
	case errors.Is(err, authapp.ErrWalletDetachNotEligible):
		return ErrWalletDetachNotEligible
	case errors.Is(err, authapp.ErrInvalidWalletSignature):
		return ErrInvalidWalletSignature
	case errors.Is(err, authapp.ErrApplicationNotConfigured):
		return ErrApplicationNotConfigured
	case errors.Is(err, authapp.ErrWalletIdentityStore):
		return ErrWalletIdentityStore
	default:
		return err
	}
}

func mapSessionViewFromApp(v *authapp.SessionView) *SessionView {
	if v == nil {
		return nil
	}
	return &SessionView{
		Authenticated: v.Authenticated,
		TokenType:     v.TokenType,
		UserID:        v.UserID,
		Email:         v.Email,
		WalletID:      v.WalletID,
		WalletAddress: v.WalletAddress,
		AuthMethod:    v.AuthMethod,
		Chain:         v.Chain,
		Subject:       v.Subject,
		Issuer:        v.Issuer,
		ExpiresAt:     v.ExpiresAt,
		User:          v.User,
	}
}

func mapProfileViewFromApp(v *authapp.ProfileView) *ProfileView {
	if v == nil {
		return nil
	}

	wallets := make([]*ProfileWalletView, 0, len(v.Wallets))
	for _, wallet := range v.Wallets {
		wallets = append(wallets, mapProfileWalletViewFromApp(wallet))
	}

	return &ProfileView{
		User:                v.User,
		UserID:              v.UserID,
		AuthMethod:          v.AuthMethod,
		WalletID:            v.WalletID,
		WalletAddress:       v.WalletAddress,
		Chain:               v.Chain,
		PrimaryWallet:       mapProfileWalletViewFromApp(v.PrimaryWallet),
		Wallets:             wallets,
		WalletCount:         v.WalletCount,
		ActiveWalletCount:   v.ActiveWalletCount,
		DetachedWalletCount: v.DetachedWalletCount,
		HasWalletSession:    v.HasWalletSession,
	}
}

func mapProfileWalletViewFromApp(v *authapp.ProfileWalletView) *ProfileWalletView {
	if v == nil {
		return nil
	}
	return &ProfileWalletView{
		ID:         v.ID,
		Address:    v.Address,
		IsPrimary:  v.IsPrimary,
		Status:     v.Status,
		LinkedAt:   v.LinkedAt,
		DetachedAt: v.DetachedAt,
	}
}

func mapWalletReadModelsFromApp(items []*authapp.WalletReadModel) []*WalletReadModel {
	out := make([]*WalletReadModel, 0, len(items))
	for _, item := range items {
		if item == nil {
			out = append(out, nil)
			continue
		}
		out = append(out, &WalletReadModel{
			ID:                 item.ID,
			Address:            item.Address,
			UserID:             item.UserID,
			LinkedAt:           item.LinkedAt,
			DetachedAt:         item.DetachedAt,
			IsPrimary:          item.IsPrimary,
			Status:             item.Status,
			CanSetPrimary:      item.CanSetPrimary,
			CanDetach:          item.CanDetach,
			DetachBlockReasons: item.DetachBlockReasons,
		})
	}
	return out
}

func mapWalletDetachCheckFromApp(v *authapp.WalletDetachCheckResponse) *WalletDetachCheckResponse {
	if v == nil {
		return nil
	}
	return &WalletDetachCheckResponse{
		WalletAddress:    v.WalletAddress,
		Eligible:         v.Eligible,
		IsPrimary:        v.IsPrimary,
		OwnedWalletCount: v.OwnedWalletCount,
		Reasons:          v.Reasons,
	}
}

func mapWalletChallengeToDomain(ch *WalletChallenge) *authdomain.WalletChallenge {
	if ch == nil {
		return nil
	}
	return &authdomain.WalletChallenge{
		ID:                ch.ID,
		Address:           ch.Address,
		Chain:             ch.Chain,
		Nonce:             ch.Nonce,
		Message:           ch.Message,
		Purpose:           ch.Purpose,
		RequestedByUserID: ch.RequestedByUserID,
		IssuedAt:          ch.IssuedAt,
		ExpiresAt:         ch.ExpiresAt,
		UsedAt:            ch.UsedAt,
	}
}

func mapWalletChallengeFromDomain(ch *authdomain.WalletChallenge) *WalletChallenge {
	if ch == nil {
		return nil
	}
	return &WalletChallenge{
		ID:                ch.ID,
		Address:           ch.Address,
		Chain:             ch.Chain,
		Nonce:             ch.Nonce,
		Message:           ch.Message,
		Purpose:           ch.Purpose,
		RequestedByUserID: ch.RequestedByUserID,
		IssuedAt:          ch.IssuedAt,
		ExpiresAt:         ch.ExpiresAt,
		UsedAt:            ch.UsedAt,
	}
}

func mapWalletIdentityFromDomain(w *authdomain.WalletIdentity) *WalletIdentity {
	if w == nil {
		return nil
	}
	return &WalletIdentity{
		ID:         w.ID,
		Address:    w.Address,
		UserID:     w.UserID,
		LinkedAt:   w.LinkedAt,
		DetachedAt: w.DetachedAt,
		IsPrimary:  w.IsPrimary,
	}
}

func mapWalletIdentitiesFromDomain(items []*authdomain.WalletIdentity) []*WalletIdentity {
	out := make([]*WalletIdentity, 0, len(items))
	for _, item := range items {
		out = append(out, mapWalletIdentityFromDomain(item))
	}
	return out
}
