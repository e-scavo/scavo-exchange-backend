package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	UserID      string `json:"user_id"`
}

type HTTPHandlers struct {
	Tokens *coreauth.TokenService
	TTL    time.Duration
	Users  *usermod.Service
}

func (h HTTPHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "bad_request"})
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Email == "" || req.Password != "dev" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid_credentials"})
		return
	}

	userID := "u_" + strings.ReplaceAll(req.Email, "@", "_")
	if h.Users != nil {
		u, err := h.Users.ResolveOrCreateDevUser(r.Context(), req.Email)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "user_persistence_error"})
			return
		}
		userID = u.ID
	}

	token, err := h.Tokens.Mint(userID, req.Email)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "token_error"})
		return
	}

	exp := h.TTL
	if exp <= 0 {
		exp = 24 * time.Hour
	}

	writeJSON(w, http.StatusOK, LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(exp.Seconds()),
		UserID:      userID,
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
