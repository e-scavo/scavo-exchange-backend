package system

import (
	"context"
	"time"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/ws"
)

func Register(d *ws.Dispatcher) {
	d.Handle("system.ping", ping)
}

func ping(ctx context.Context, c *ws.Client, env ws.Envelope) ws.Envelope {
	return ws.Envelope{
		ID:     env.ID,
		Type:   ws.MsgTypeRes,
		Action: "system.ping",
		Data: ws.JSON(map[string]any{
			"pong": true,
			"ts":   time.Now().UTC().Format(time.RFC3339),
		}),
	}
}
