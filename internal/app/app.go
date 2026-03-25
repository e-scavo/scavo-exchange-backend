package app

import (
	"context"
	"net/http"
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/cache"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/config"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/db"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/httpx"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/status"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/ws"
	authmod "github.com/e-scavo/scavo-exchange-backend/internal/modules/auth"
	"github.com/e-scavo/scavo-exchange-backend/internal/modules/system"
	usermod "github.com/e-scavo/scavo-exchange-backend/internal/modules/user"
)

type App struct {
	cfg    config.Config
	log    *logger.Logger
	server *http.Server

	hub        *ws.Hub
	dispatcher *ws.Dispatcher
	tokens     *coreauth.TokenService

	dbClient    *db.Client
	cacheClient *cache.Client
	statusSvc   *status.Service

	userService *usermod.Service
	authService *authmod.Service
}

func New(cfg config.Config) *App {
	lg := logger.New(cfg.Env)

	hub := ws.NewHub(lg)
	dispatcher := ws.NewDispatcher()

	system.Register(dispatcher)

	ttl := time.Duration(cfg.JWTTTLHrs) * time.Hour
	tokens, err := coreauth.NewTokenService(cfg.JWTSecret, cfg.JWTIssuer, ttl)
	if err != nil {
		lg.Error("jwt config invalid", "err", err)
		tokens, _ = coreauth.NewTokenService("dev_dev_dev_dev_dev_dev_dev_dev", "scavo-exchange-backend", 24*time.Hour)
	}

	initCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbClient, dbErr := db.New(initCtx, cfg, lg)
	if dbErr != nil {
		lg.Error("postgres init failed", "err", dbErr)
	}

	cacheClient, cacheErr := cache.New(cfg, lg)
	if cacheErr != nil {
		lg.Error("redis init failed", "err", cacheErr)
	}

	var userService *usermod.Service
	if dbClient != nil && dbClient.Enabled() && dbClient.Pool() != nil {
		userRepo := usermod.NewPostgresRepository(dbClient.Pool(), lg)
		userService = usermod.NewService(userRepo)
	} else {
		userService = usermod.NewService(nil)
	}

	authService := authmod.NewService(tokens, userService, ttl)
	authmod.RegisterWS(dispatcher, authService)

	statusSvc := status.New(
		"scavo-exchange-backend",
		cfg.Env,
		cfg.Version,
		cfg.Commit,
		status.FuncChecker{
			NameValue:     "postgres",
			RequiredValue: cfg.ReadinessRequirePostgres,
			Fn: func(ctx context.Context) error {
				if dbErr != nil {
					return dbErr
				}
				if dbClient == nil {
					return db.ErrNotConfigured
				}
				return dbClient.Ping(ctx)
			},
		},
		status.FuncChecker{
			NameValue:     "redis",
			RequiredValue: cfg.ReadinessRequireRedis,
			Fn: func(ctx context.Context) error {
				if cacheErr != nil {
					return cacheErr
				}
				if cacheClient == nil {
					return cache.ErrNotConfigured
				}
				return cacheClient.Ping(ctx)
			},
		},
	)

	r := httpx.NewRouter(httpx.RouterParams{
		Log:          lg,
		Hub:          hub,
		Dispatcher:   dispatcher,
		Config:       cfg,
		TokenService: tokens,
		Status:       statusSvc,
		UserService:  userService,
	})

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: r,
	}

	return &App{
		cfg:         cfg,
		log:         lg,
		server:      srv,
		hub:         hub,
		dispatcher:  dispatcher,
		tokens:      tokens,
		dbClient:    dbClient,
		cacheClient: cacheClient,
		statusSvc:   statusSvc,
		userService: userService,
		authService: authService,
	}
}

func (a *App) Start(ctx context.Context) error {
	go a.hub.Run(ctx)

	a.log.Info("http server starting",
		"addr", a.cfg.HTTPAddr,
		"env", a.cfg.Env,
		"postgres_enabled", a.dbClient != nil && a.dbClient.Enabled(),
		"redis_enabled", a.cacheClient != nil && a.cacheClient.Enabled(),
	)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error("http server error", "err", err)
		}
	}()

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info("http server stopping")

	if a.dbClient != nil {
		a.dbClient.Close()
	}

	if a.cacheClient != nil {
		if err := a.cacheClient.Close(); err != nil {
			a.log.Error("redis close error", "err", err)
		}
	}

	return a.server.Shutdown(ctx)
}
