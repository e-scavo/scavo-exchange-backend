package httpx

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/config"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/ws"
	authmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"
)

type RouterParams struct {
	Log        *logger.Logger
	Hub        *ws.Hub
	Dispatcher *ws.Dispatcher
	Config     config.Config

	TokenService *coreauth.TokenService
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

	// ✅ WS sin AccessLog/Timeout (no romper Upgrade)
	r.Get("/ws", ws.NewHandler(ws.HandlerParams{
		Log:        p.Log,
		Hub:        p.Hub,
		Dispatcher: p.Dispatcher,
		TokenSvc:   p.TokenService,
	}))

	// ✅ HTTP normal
	r.Group(func(r chi.Router) {
		r.Use(AccessLog(p.Log))
		r.Use(Timeout(30 * time.Second))

		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
		})

		r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
			WriteJSON(w, http.StatusOK, map[string]any{
				"version": p.Config.Version,
				"commit":  p.Config.Commit,
				"env":     p.Config.Env,
			})
		})

		handlers := authmod.HTTPHandlers{
			Tokens: p.TokenService,
			TTL:    time.Duration(p.Config.JWTTTLHrs) * time.Hour,
		}
		r.Post("/auth/login", handlers.Login)
	})

	return r
}
