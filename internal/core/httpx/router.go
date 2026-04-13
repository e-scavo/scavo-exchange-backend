package httpx

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/config"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/status"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/ws"
	authmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
	usersettingsmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/usersettings"
)

type RouterParams struct {
	Log        *logger.Logger
	Hub        *ws.Hub
	Dispatcher *ws.Dispatcher
	Config     config.Config

	TokenService        *coreauth.TokenService
	Status              *status.Service
	UserService         *usermod.Service
	UserSettingsService *usersettingsmod.Service
	ChallengeStore      authmod.WalletChallengeStore
	WalletIdentityStore authmod.WalletIdentityStore
	ChallengeTTL        time.Duration
	PublicBaseURL       string
}

func NewRouter(p RouterParams) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   p.Config.CORSAllowOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		ExposedHeaders:   []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(RequestID())
	r.Use(Recoverer(p.Log))

	r.Get("/ws", ws.NewHandler(ws.HandlerParams{
		Log:        p.Log,
		Hub:        p.Hub,
		Dispatcher: p.Dispatcher,
		TokenSvc:   p.TokenService,
	}))

	r.Group(func(r chi.Router) {
		r.Use(AccessLog(p.Log))
		r.Use(Timeout(30 * time.Second))

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			if p.Status == nil {
				WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "status": "up"})
				return
			}
			WriteJSON(w, http.StatusOK, p.Status.Health())
		})

		r.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
			if p.Status == nil {
				WriteJSON(w, http.StatusOK, map[string]any{"ok": true, "status": "ready"})
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
			defer cancel()

			code, payload := p.Status.Readiness(ctx)
			WriteJSON(w, code, payload)
		})

		r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
			WriteJSON(w, http.StatusOK, map[string]any{
				"version": p.Config.Version,
				"commit":  p.Config.Commit,
				"env":     p.Config.Env,
			})
		})

		handlers := authmod.HTTPHandlers{
			Tokens:           p.TokenService,
			TTL:              time.Duration(p.Config.JWTTTLHrs) * time.Hour,
			Users:            p.UserService,
			UserSettings:     p.UserSettingsService,
			PublicBaseURL:    p.PublicBaseURL,
			ChallengeTTL:     p.ChallengeTTL,
			Challenges:       p.ChallengeStore,
			WalletIdentities: p.WalletIdentityStore,
		}

		registerAuthRoutes(r, "", p.TokenService, handlers)
		registerAuthRoutes(r, "/api/v1", p.TokenService, handlers)
	})

	return r
}

func registerAuthRoutes(r chi.Router, prefix string, tokens *coreauth.TokenService, handlers authmod.HTTPHandlers) {
	r.Post(routePath(prefix, "/auth/login"), handlers.Login)
	r.Post(routePath(prefix, "/auth/wallet/challenge"), handlers.WalletChallenge)
	r.Post(routePath(prefix, "/auth/wallet/verify"), handlers.WalletVerify)

	requireAuth := RequireAuth(tokens, false)

	r.With(requireAuth).Post(routePath(prefix, "/auth/wallets/link/challenge"), handlers.WalletLinkChallenge)
	r.With(requireAuth).Post(routePath(prefix, "/auth/wallets/link/verify"), handlers.WalletLinkVerify)
	r.With(requireAuth).Post(routePath(prefix, "/auth/account/merge/wallet/challenge"), handlers.WalletAccountMergeChallenge)
	r.With(requireAuth).Post(routePath(prefix, "/auth/account/merge/wallet/verify"), handlers.WalletAccountMergeVerify)

	r.With(requireAuth).Get(routePath(prefix, "/auth/bootstrap"), handlers.Bootstrap)
	r.With(requireAuth).Get(routePath(prefix, "/auth/me"), handlers.Me)
	r.With(requireAuth).Patch(routePath(prefix, "/auth/me"), handlers.UpdateMe)
	r.With(requireAuth).Get(routePath(prefix, "/auth/me/settings"), handlers.MeSettings)
	r.With(requireAuth).Patch(routePath(prefix, "/auth/me/settings"), handlers.UpdateMeSettings)
	r.With(requireAuth).Get(routePath(prefix, "/auth/session"), handlers.Session)
	r.With(requireAuth).Get(routePath(prefix, "/auth/wallets"), handlers.Wallets)
	r.With(requireAuth).Post(routePath(prefix, "/auth/wallets/detach/check"), handlers.WalletDetachCheck)
	r.With(requireAuth).Post(routePath(prefix, "/auth/wallets/detach"), handlers.WalletDetach)
	r.With(requireAuth).Post(routePath(prefix, "/auth/wallets/primary"), handlers.WalletSetPrimary)
}

func routePath(prefix, suffix string) string {
	if prefix == "" {
		return suffix
	}
	return prefix + suffix
}
