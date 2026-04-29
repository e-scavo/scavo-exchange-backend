package auth

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"

	coreerrs "github.com/e-scavo/scavo-exchange-backend/internal/core/errs"
	authapp "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/app"
	authmappers "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth/mappers"
)

type WalletReadModel = authapp.WalletReadModel
type WalletsResponse = authapp.WalletsResponse
type WalletsQuery = authapp.WalletsQuery

func mapWalletIdentityToReadModel(wallet *WalletIdentity) *WalletReadModel {
	return authmappers.WalletIdentityToReadModel(wallet)
}

func mapWalletIdentitiesToReadModels(wallets []*WalletIdentity) []*WalletReadModel {
	return authmappers.WalletIdentitiesToActionableReadModels(wallets)
}

func enrichWalletReadModelsActionability(wallets []*WalletReadModel) []*WalletReadModel {
	return authmappers.EnrichWalletReadModelsActionability(wallets)
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
		q.SortProvided = true
	}

	order := strings.TrimSpace(strings.ToLower(params.Get("order")))
	if order != "" {
		switch order {
		case "asc", "desc":
			q.Order = order
			q.OrderProvided = true
		default:
			return WalletsQuery{}, "invalid_order"
		}
	}

	limit := strings.TrimSpace(params.Get("limit"))
	if limit != "" {
		value, ok := parsePositiveInt(limit)
		if !ok {
			return WalletsQuery{}, "invalid_limit"
		}
		q.Limit = value
		q.LimitProvided = true
	}

	offset := strings.TrimSpace(params.Get("offset"))
	if offset != "" {
		value, ok := parseNonNegativeInt(offset)
		if !ok {
			return WalletsQuery{}, "invalid_offset"
		}
		q.Offset = value
		q.OffsetProvided = true
	}

	if errCode := validateWalletsQueryContract(&q); errCode != "" {
		return WalletsQuery{}, errCode
	}

	return q, ""
}

func validateWalletsQueryContract(q *WalletsQuery) string {
	if q == nil {
		return ""
	}

	if q.OrderProvided && !q.SortProvided {
		return "invalid_order_requires_sort"
	}

	if q.SortProvided && !q.OrderProvided {
		q.Order = "asc"
	}

	return ""
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

func (h HTTPHandlers) Wallets(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.requireClaims(w, r)
	if !ok {
		return
	}

	query, queryErr := parseWalletsQuery(r)
	if queryErr != "" {
		writeAppErrorJSON(w, coreerrs.AppErrorFromLegacyAuthKey(queryErr, nil))
		return
	}

	response, err := h.AuthProvider().ListWallets(r.Context(), claims.UserID, query)
	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			writeAppErrorJSON(w, coreerrs.AuthUnauthorized())
		case errors.Is(err, ErrApplicationNotConfigured):
			writeAppErrorJSON(w, coreerrs.WalletIdentityError(nil))
		case errors.Is(err, ErrWalletIdentityStore):
			writeAppErrorJSON(w, coreerrs.WalletIdentityError(nil))
		default:
			writeAppErrorJSON(w, coreerrs.WalletIdentityError(nil))
		}
		return
	}

	writeJSON(w, http.StatusOK, response)
}
