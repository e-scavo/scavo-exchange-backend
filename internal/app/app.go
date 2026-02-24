package app

import (
	"context"
	"net/http"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/config"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/httpx"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/ws"
	"github.com/e-scavo/scavo-exchange-backend/internal/modules/system"
)

type App struct {
	cfg    config.Config
	log    *logger.Logger
	server *http.Server

	hub        *ws.Hub
	dispatcher *ws.Dispatcher
}

func New(cfg config.Config) *App {
	lg := logger.New(cfg.Env)

	hub := ws.NewHub(lg)
	dispatcher := ws.NewDispatcher()

	// Register WS modules
	system.Register(dispatcher)

	r := httpx.NewRouter(httpx.RouterParams{
		Log:        lg,
		Hub:        hub,
		Dispatcher: dispatcher,
		Config:     cfg,
	})

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: r,
	}

	return &App{
		cfg:        cfg,
		log:        lg,
		server:     srv,
		hub:        hub,
		dispatcher: dispatcher,
	}
}

func (a *App) Start(ctx context.Context) error {
	go a.hub.Run(ctx)

	a.log.Info("http server starting", "addr", a.cfg.HTTPAddr, "env", a.cfg.Env)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.Error("http server error", "err", err)
		}
	}()

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info("http server stopping")
	return a.server.Shutdown(ctx)
}
