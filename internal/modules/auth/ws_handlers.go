package auth

import (
	"context"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/ws"
)

type WSHandlers struct {
	Service *Service
}

func RegisterWS(d *ws.Dispatcher, svc *Service) {
	h := WSHandlers{Service: svc}

	d.Handle("auth.whoami", ws.RequireAuth(h.whoami))
	d.Handle("auth.session", ws.RequireAuth(h.session))
}

func (h WSHandlers) whoami(ctx context.Context, c *ws.Client, env ws.Envelope) ws.Envelope {
	s := c.Session()
	if s == nil || s.UserID == "" {
		return ws.Envelope{
			ID:     env.ID,
			Type:   ws.MsgTypeRes,
			Action: env.Action,
			Error:  &ws.ErrPayload{Code: "unauthorized", Msg: "missing session"},
		}
	}

	data := map[string]any{
		"authenticated": true,
		"user_id":       s.UserID,
		"email":         s.Email,
	}
	if s.WalletAddress != "" {
		data["wallet_address"] = s.WalletAddress
	}
	if s.AuthMethod != "" {
		data["auth_method"] = s.AuthMethod
	}
	if s.Chain != "" {
		data["chain"] = s.Chain
	}

	if s.Subject != "" {
		data["subject"] = s.Subject
	}
	if s.Issuer != "" {
		data["issuer"] = s.Issuer
	}
	if s.ExpiresAt != nil {
		data["expires_at"] = s.ExpiresAt.UTC()
	}

	return ws.Envelope{
		ID:     env.ID,
		Type:   ws.MsgTypeRes,
		Action: env.Action,
		Data:   ws.JSON(data),
	}
}

func (h WSHandlers) session(ctx context.Context, c *ws.Client, env ws.Envelope) ws.Envelope {
	s := c.Session()
	if s == nil || s.Claims == nil || s.UserID == "" {
		return ws.Envelope{
			ID:     env.ID,
			Type:   ws.MsgTypeRes,
			Action: env.Action,
			Error:  &ws.ErrPayload{Code: "unauthorized", Msg: "missing session"},
		}
	}

	if h.Service == nil {
		return ws.Envelope{
			ID:     env.ID,
			Type:   ws.MsgTypeRes,
			Action: env.Action,
			Error:  &ws.ErrPayload{Code: "auth_service_error", Msg: "auth service not configured"},
		}
	}

	view, err := h.Service.ResolveSessionClaims(ctx, s.Claims)
	if err != nil {
		return ws.Envelope{
			ID:     env.ID,
			Type:   ws.MsgTypeRes,
			Action: env.Action,
			Error:  &ws.ErrPayload{Code: "auth_service_error", Msg: err.Error()},
		}
	}

	return ws.Envelope{
		ID:     env.ID,
		Type:   ws.MsgTypeRes,
		Action: env.Action,
		Data:   ws.JSON(map[string]any{"session": view}),
	}
}
