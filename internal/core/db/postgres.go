package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/config"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
)

var ErrNotConfigured = errors.New("postgres not configured")

type Client struct {
	pool *pgxpool.Pool
	log  *logger.Logger
}

func New(ctx context.Context, cfg config.Config, log *logger.Logger) (*Client, error) {
	if strings.TrimSpace(cfg.PostgresURL) == "" {
		return &Client{pool: nil, log: log}, nil
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.PostgresURL)
	if err != nil {
		return nil, err
	}

	poolCfg.MaxConnIdleTime = 5 * time.Minute
	poolCfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		pool: pool,
		log:  log,
	}, nil
}

func (c *Client) Enabled() bool {
	return c != nil && c.pool != nil
}

func (c *Client) Pool() *pgxpool.Pool {
	if c == nil {
		return nil
	}
	return c.pool
}

func (c *Client) Ping(ctx context.Context) error {
	if !c.Enabled() {
		return ErrNotConfigured
	}
	return c.pool.Ping(ctx)
}

func (c *Client) Close() {
	if c != nil && c.pool != nil {
		c.pool.Close()
	}
}
