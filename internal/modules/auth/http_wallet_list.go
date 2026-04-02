package auth

import (
	"net/http"
	"sort"
	"strconv"
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
	Wallets        []*WalletReadModel `json:"wallets"`
	Total          int                `json:"total"`
	Limit          int                `json:"limit"`
	Offset         int                `json:"offset"`
	Returned       int                `json:"returned"`
	HasMore        bool               `json:"has_more"`
	NextOffset     *int               `json:"next_offset,omitempty"`
	PreviousOffset *int               `json:"previous_offset,omitempty"`
}

type WalletsQuery struct {
	Status  string
	Primary *bool
	Sort    string
	Order   string
	Limit   int
	Offset  int
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

	limit := strings.TrimSpace(params.Get("limit"))
	if limit != "" {
		value, ok := parsePositiveInt(limit)
		if !ok {
			return WalletsQuery{}, "invalid_limit"
		}
		q.Limit = value
	}

	offset := strings.TrimSpace(params.Get("offset"))
	if offset != "" {
		value, ok := parseNonNegativeInt(offset)
		if !ok {
			return WalletsQuery{}, "invalid_offset"
		}
		q.Offset = value
	}

	return q, ""
}

func parsePositiveInt(raw string) (int, bool) {
	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, false
	}
	return value, true
}

func parseNonNegativeInt(raw string) (int, bool) {
	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0, false
	}
	return value, true
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
		writeJSON(w, http.StatusOK, buildWalletsResponse([]*WalletReadModel{}, 0, query))
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
	window, total := applyWalletsQuery(mapped, query)
	writeJSON(w, http.StatusOK, buildWalletsResponse(window, total, query))
}
