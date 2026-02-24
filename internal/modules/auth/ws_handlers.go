package auth

import (
	"context"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/ws"
)

func RegisterWS(d *ws.Dispatcher) {
	d.Handle("auth.whoami", whoami)
}

func whoami(ctx context.Context, c *ws.Client, env ws.Envelope) ws.Envelope {
	s := c.Session()
	if s == nil || s.UserID == "" {
		return ws.Envelope{
			ID:     env.ID,
			Type:   ws.MsgTypeRes,
			Action: env.Action,
			Error:  &ws.ErrPayload{Code: "unauthorized", Msg: "missing session"},
		}
	}

	return ws.Envelope{
		ID:     env.ID,
		Type:   ws.MsgTypeRes,
		Action: env.Action,
		Data: ws.JSON(map[string]any{
			"user_id": s.UserID,
			"email":   s.Email,
		}),
	}
}
