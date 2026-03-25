package ws

import (
	"time"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
)

type Session struct {
	Claims        *coreauth.Claims
	UserID        string
	Email         string
	WalletID      string
	WalletAddress string
	AuthMethod    string
	Chain         string
	Subject       string
	Issuer        string
	ExpiresAt     *time.Time
}

func (c *Client) SetSession(s Session) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.session = &s
}

func (c *Client) Session() *Session {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.session
}
