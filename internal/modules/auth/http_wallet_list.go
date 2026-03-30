package auth

import (
	"net/http"
	"sort"
	"strings"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
)

type WalletReadModel struct {
	ID         string     `json:"id"`
	Address    string     `json:"address"`
	UserID     string     `json:"user_id,omitempty"`
	LinkedAt   *time.Time `json:"linked_at,omitempty"`
	DetachedAt *time.Time `json:"detached_at,omitempty"`
	IsPrimary  bool       `json:"is_primary"`
	Status     string     `json:"status"`
}

type WalletsResponse struct {
	Wallets []*WalletReadModel `json:"wallets"`
}

type WalletsQuery struct {
	Status  string
	Primary *bool
	Sort    string
	Order   string
}

func mapWalletIdentityToReadModel(wallet *WalletIdentity) *WalletReadModel {
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

func mapWalletIdentitiesToReadModels(wallets []*WalletIdentity) []*WalletReadModel {
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

	return out
}

func parseWalletsQuery(r *http.Request) (WalletsQuery, string) {
	q := WalletsQuery{}
	params := r.URL.Query()

	status := strings.TrimSpace(strings.ToLower(params.Get("status")))
	if status != "" {
		switch status {
		case "active", "detached":
			q.Status = status
		default:
			return WalletsQuery{}, "invalid_status"
		}
	}

	primary := strings.TrimSpace(strings.ToLower(params.Get("primary")))
	if primary != "" {
		switch primary {
		case "true":
			v := true
			q.Primary = &v
		case "false":
			v := false
			q.Primary = &v
		default:
			return WalletsQuery{}, "invalid_primary"
		}
	}

	sortBy := strings.TrimSpace(strings.ToLower(params.Get("sort")))
	if sortBy != "" {
		if sortBy != "linked_at" {
			return WalletsQuery{}, "invalid_sort"
		}
		q.Sort = sortBy
	}

	order := strings.TrimSpace(strings.ToLower(params.Get("order")))
	if order != "" {
		switch order {
		case "asc", "desc":
			q.Order = order
		default:
			return WalletsQuery{}, "invalid_order"
		}
		if q.Sort == "" {
			return WalletsQuery{}, "invalid_sort"
		}
	}

	return q, ""
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

func applyWalletsQuery(wallets []*WalletReadModel, q WalletsQuery) []*WalletReadModel {
	filtered := filterWalletReadModels(wallets, q)
	return sortWalletReadModels(filtered, q)
}

func (h HTTPHandlers) Wallets(w http.ResponseWriter, r *http.Request) {
	claims, ok := coreauth.ClaimsFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
		return
	}

	query, queryErr := parseWalletsQuery(r)
	if queryErr != "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": queryErr})
		return
	}

	if h.WalletIdentities == nil {
		writeJSON(w, http.StatusOK, WalletsResponse{Wallets: []*WalletReadModel{}})
		return
	}

	wallets, err := h.WalletIdentities.ListByUser(r.Context(), claims.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "wallet_identity_error"})
		return
	}
	if wallets == nil {
		wallets = []*WalletIdentity{}
	}

	mapped := mapWalletIdentitiesToReadModels(wallets)
	writeJSON(w, http.StatusOK, WalletsResponse{Wallets: applyWalletsQuery(mapped, query)})
}
