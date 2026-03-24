package cache

import (
	"context"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/config"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
)

var ErrNotConfigured = errors.New("redis not configured")

type Client struct {
	rdb *redis.Client
	log *logger.Logger
}

func New(cfg config.Config, log *logger.Logger) (*Client, error) {
	if strings.TrimSpace(cfg.RedisAddr) == "" {
		return &Client{rdb: nil, log: log}, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &Client{
		rdb: rdb,
		log: log,
	}, nil
}

func (c *Client) Enabled() bool {
	return c != nil && c.rdb != nil
}

func (c *Client) Redis() *redis.Client {
	if c == nil {
		return nil
	}
	return c.rdb
}

func (c *Client) Ping(ctx context.Context) error {
	if !c.Enabled() {
		return ErrNotConfigured
	}
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) Close() error {
	if c != nil && c.rdb != nil {
		return c.rdb.Close()
	}
	return nil
}
