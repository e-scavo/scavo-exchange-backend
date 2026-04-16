package app

import (
	"context"
	"sort"
	"strings"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	rootauth "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"
)

type authenticatedContextView struct {
	UserID           string
	Email            string
	WalletID         string
	WalletAddress    string
	AuthMethod       string
	Chain            string
	HasWalletSession bool
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func normalizeWalletAddress(address string) string {
	return strings.ToLower(strings.TrimSpace(address))
}

func normalizeChain(chain string) string {
	chain = strings.TrimSpace(strings.ToLower(chain))
	if chain == "" {
		return "scavium"
	}
	return chain
}

func buildAuthenticatedContextView(claims *coreauth.Claims) authenticatedContextView {
	if claims == nil {
		return authenticatedContextView{AuthMethod: "password_dev"}
	}

	view := authenticatedContextView{
		UserID:        strings.TrimSpace(claims.UserID),
		Email:         normalizeEmail(claims.Email),
		WalletID:      strings.TrimSpace(claims.WalletID),
		WalletAddress: normalizeWalletAddress(claims.WalletAddress),
		AuthMethod:    strings.TrimSpace(claims.AuthMethod),
		Chain:         normalizeChain(claims.Chain),
	}

	if view.AuthMethod == "" {
		view.AuthMethod = "password_dev"
	}
	if view.WalletAddress == "" {
		view.WalletID = ""
		view.Chain = ""
	}
	view.HasWalletSession = view.WalletAddress != ""

	return view
}

func buildSessionViewWithUser(claims *coreauth.Claims, user *usermod.User) *SessionView {
	ctxView := buildAuthenticatedContextView(claims)

	var expiresAt *time.Time
	if claims != nil && claims.ExpiresAt != nil {
		ts := claims.ExpiresAt.Time.UTC()
		expiresAt = &ts
	}

	view := &SessionView{
		Authenticated: true,
		TokenType:     "Bearer",
		UserID:        ctxView.UserID,
		Email:         ctxView.Email,
		WalletID:      ctxView.WalletID,
		WalletAddress: ctxView.WalletAddress,
		AuthMethod:    ctxView.AuthMethod,
		Chain:         ctxView.Chain,
		ExpiresAt:     expiresAt,
		User:          user,
	}
	if claims != nil {
		view.Subject = strings.TrimSpace(claims.Subject)
		view.Issuer = strings.TrimSpace(claims.Issuer)
	}
	if view.Subject == "" {
		view.Subject = view.UserID
	}

	return view
}

func buildProfileView(ctx context.Context, claims *coreauth.Claims, users *usermod.Service, walletStore rootauth.WalletIdentityStore) (*ProfileView, error) {
	svc := NewService(nil, users, 24*time.Hour)
	user, err := svc.ResolveCurrentUserClaims(ctx, claims)
	if err != nil {
		return nil, err
	}
	return buildProfileViewWithUser(ctx, claims, user, walletStore)
}

func buildProfileViewWithUser(ctx context.Context, claims *coreauth.Claims, user *usermod.User, walletStore rootauth.WalletIdentityStore) (*ProfileView, error) {
	ctxView := buildAuthenticatedContextView(claims)

	view := &ProfileView{
		User:             user,
		UserID:           ctxView.UserID,
		AuthMethod:       ctxView.AuthMethod,
		WalletID:         ctxView.WalletID,
		WalletAddress:    ctxView.WalletAddress,
		Chain:            ctxView.Chain,
		Wallets:          []*ProfileWalletView{},
		HasWalletSession: ctxView.HasWalletSession,
	}
	if walletStore == nil || view.UserID == "" {
		return view, nil
	}

	wallets, err := walletStore.ListByUser(ctx, view.UserID)
	if err != nil {
		return nil, err
	}

	for _, wallet := range wallets {
		mapped := mapWalletIdentityToProfileWallet(wallet)
		if mapped == nil {
			continue
		}
		view.Wallets = append(view.Wallets, mapped)
		view.WalletCount++
		switch mapped.Status {
		case "active":
			view.ActiveWalletCount++
		case "detached":
			view.DetachedWalletCount++
		}
		if mapped.IsPrimary && view.PrimaryWallet == nil {
			copy := *mapped
			view.PrimaryWallet = &copy
		}
	}

	return view, nil
}

func mapWalletIdentityToProfileWallet(wallet *rootauth.WalletIdentity) *ProfileWalletView {
	mapped := mapWalletIdentityToReadModel(wallet)
	if mapped == nil {
		return nil
	}

	return &ProfileWalletView{
		ID:         mapped.ID,
		Address:    mapped.Address,
		IsPrimary:  mapped.IsPrimary,
		Status:     mapped.Status,
		LinkedAt:   mapped.LinkedAt,
		DetachedAt: mapped.DetachedAt,
	}
}

func listWalletReadModelsForUser(ctx context.Context, userID string, store rootauth.WalletIdentityStore) ([]*WalletReadModel, error) {
	if store == nil || userID == "" {
		return []*WalletReadModel{}, nil
	}

	wallets, err := store.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallets == nil {
		wallets = []*rootauth.WalletIdentity{}
	}

	return mapWalletIdentitiesToReadModels(wallets), nil
}

func buildBootstrapWalletsView(wallets []*WalletReadModel) BootstrapWalletsView {
	if wallets == nil {
		wallets = []*WalletReadModel{}
	}

	return BootstrapWalletsView{
		Items: wallets,
		Total: len(wallets),
	}
}

func mapWalletIdentityToReadModel(wallet *rootauth.WalletIdentity) *WalletReadModel {
	if wallet == nil {
		return nil
	}

	status := "unlinked"
	switch {
	case wallet.UserID != "":
		status = "active"
	case wallet.DetachedAt != nil:
		status = "detached"
	}

	return &WalletReadModel{
		ID:         wallet.ID,
		Address:    wallet.Address,
		UserID:     wallet.UserID,
		LinkedAt:   wallet.LinkedAt,
		DetachedAt: wallet.DetachedAt,
		IsPrimary:  wallet.IsPrimary,
		Status:     status,
	}
}

func mapWalletIdentitiesToReadModels(wallets []*rootauth.WalletIdentity) []*WalletReadModel {
	if len(wallets) == 0 {
		return []*WalletReadModel{}
	}

	out := make([]*WalletReadModel, 0, len(wallets))
	for _, wallet := range wallets {
		mapped := mapWalletIdentityToReadModel(wallet)
		if mapped != nil {
			out = append(out, mapped)
		}
	}

	if out == nil {
		return []*WalletReadModel{}
	}

	return enrichWalletReadModelsActionability(out)
}

func enrichWalletReadModelsActionability(wallets []*WalletReadModel) []*WalletReadModel {
	if len(wallets) == 0 {
		return []*WalletReadModel{}
	}

	activeOwnedCount := 0
	for _, wallet := range wallets {
		if wallet != nil && wallet.Status == "active" {
			activeOwnedCount++
		}
	}

	for _, wallet := range wallets {
		if wallet == nil {
			continue
		}

		wallet.CanSetPrimary = wallet.Status == "active" && !wallet.IsPrimary
		wallet.CanDetach = false
		wallet.DetachBlockReasons = []string{}

		if wallet.Status != "active" {
			continue
		}

		if wallet.IsPrimary {
			wallet.DetachBlockReasons = append(wallet.DetachBlockReasons, rootauth.WalletDetachReasonWalletIsPrimary)
		}
		if activeOwnedCount <= 1 {
			wallet.DetachBlockReasons = append(wallet.DetachBlockReasons, rootauth.WalletDetachReasonUserWouldBeEmpty)
		}
		if len(wallet.DetachBlockReasons) == 0 {
			wallet.CanDetach = true
		}
	}

	return wallets
}

func filterWalletReadModels(wallets []*WalletReadModel, q WalletsQuery) []*WalletReadModel {
	if len(wallets) == 0 {
		return []*WalletReadModel{}
	}

	out := make([]*WalletReadModel, 0, len(wallets))
	for _, wallet := range wallets {
		if wallet == nil {
			continue
		}
		if q.Status != "" && wallet.Status != q.Status {
			continue
		}
		if q.Primary != nil && wallet.IsPrimary != *q.Primary {
			continue
		}
		out = append(out, wallet)
	}

	if out == nil {
		return []*WalletReadModel{}
	}

	return out
}

func sortWalletReadModels(wallets []*WalletReadModel, q WalletsQuery) []*WalletReadModel {
	if len(wallets) <= 1 || q.Sort == "" {
		if wallets == nil {
			return []*WalletReadModel{}
		}
		return wallets
	}

	out := make([]*WalletReadModel, 0, len(wallets))
	out = append(out, wallets...)

	desc := q.Order == "desc"
	sort.SliceStable(out, func(i, j int) bool {
		left := out[i]
		right := out[j]

		switch {
		case left == nil && right == nil:
			return false
		case left == nil:
			return false
		case right == nil:
			return true
		}

		switch {
		case left.LinkedAt == nil && right.LinkedAt == nil:
			return left.Address < right.Address
		case left.LinkedAt == nil:
			return false
		case right.LinkedAt == nil:
			return true
		case left.LinkedAt.Equal(*right.LinkedAt):
			return left.Address < right.Address
		case desc:
			return left.LinkedAt.After(*right.LinkedAt)
		default:
			return left.LinkedAt.Before(*right.LinkedAt)
		}
	})

	return out
}

func paginateWalletReadModels(wallets []*WalletReadModel, q WalletsQuery) []*WalletReadModel {
	if len(wallets) == 0 {
		return []*WalletReadModel{}
	}

	if q.Offset >= len(wallets) {
		return []*WalletReadModel{}
	}

	start := q.Offset
	end := len(wallets)
	if q.Limit > 0 && start+q.Limit < end {
		end = start + q.Limit
	}

	out := make([]*WalletReadModel, 0, end-start)
	out = append(out, wallets[start:end]...)
	return out
}

func applyWalletsQuery(wallets []*WalletReadModel, q WalletsQuery) ([]*WalletReadModel, int) {
	filtered := filterWalletReadModels(wallets, q)
	sorted := sortWalletReadModels(filtered, q)
	total := len(sorted)
	return paginateWalletReadModels(sorted, q), total
}

func buildWalletsResponse(window []*WalletReadModel, total int, q WalletsQuery) WalletsResponse {
	if window == nil {
		window = []*WalletReadModel{}
	}

	returned := len(window)
	hasMore := false
	var nextOffset *int
	var previousOffset *int

	if q.Limit > 0 {
		hasMore = q.Offset+returned < total
		if hasMore {
			v := q.Offset + returned
			nextOffset = &v
		}
		if q.Offset > 0 {
			v := q.Offset - q.Limit
			if v < 0 {
				v = 0
			}
			previousOffset = &v
		}
	}

	return WalletsResponse{
		Items:          window,
		Wallets:        window,
		Total:          total,
		Limit:          q.Limit,
		Offset:         q.Offset,
		Returned:       returned,
		HasMore:        hasMore,
		NextOffset:     nextOffset,
		PreviousOffset: previousOffset,
	}
}
