package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"

	coreauth "github.com/e-scavo/scavo-exchange-backend/internal/core/auth"
	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
)

type HandlerParams struct {
	Log        *logger.Logger
	Hub        *Hub
	Dispatcher *Dispatcher
	TokenSvc   *coreauth.TokenService
}

type Handler struct {
	log        *logger.Logger
	hub        *Hub
	dispatcher *Dispatcher
	tokenSvc   *coreauth.TokenService
}

func NewHandler(p HandlerParams) http.HandlerFunc {
	h := &Handler{
		log:        p.Log,
		hub:        p.Hub,
		dispatcher: p.Dispatcher,
		tokenSvc:   p.TokenSvc,
	}
	return h.serveWS
}

func (h *Handler) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		CompressionMode:    websocket.CompressionDisabled,
	})
	if err != nil {
		return
	}

	client := NewClient(conn)
	h.hub.Register(client)

	ctx := r.Context()
	writeCtx, cancelWrite := context.WithCancel(ctx)
	defer cancelWrite()

	go client.WriteLoop(writeCtx)

	h.tryAuth(r, client)

	hello := Envelope{
		ID:     uuid.NewString(),
		Type:   MsgTypeEvt,
		Action: "system.hello",
		Data: JSON(map[string]any{
			"ts":   time.Now().UTC().Format(time.RFC3339),
			"auth": client.Session() != nil,
		}),
	}
	client.TrySend(mustMarshal(hello))

	for {
		rctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		typ, data, err := conn.Read(rctx)
		cancel()

		if err != nil {
			client.Close(websocket.StatusNormalClosure, "bye")
			h.hub.Unregister(client)
			return
		}
		if typ != websocket.MessageText {
			continue
		}

		var env Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			res := Envelope{
				ID:     uuid.NewString(),
				Type:   MsgTypeRes,
				Action: "system.error",
				Error:  &ErrPayload{Code: "bad_json", Msg: "invalid json"},
			}
			client.TrySend(mustMarshal(res))
			continue
		}

		res := h.dispatcher.Dispatch(ctx, client, env)
		client.TrySend(mustMarshal(res))
	}
}

func (h *Handler) tryAuth(r *http.Request, c *Client) {
	if h.tokenSvc == nil {
		return
	}

	token := coreauth.ExtractTokenFromRequest(r, true)
	if token == "" {
		return
	}

	claims, err := h.tokenSvc.Parse(token)
	if err != nil || claims == nil || claims.UserID == "" {
		return
	}

	var expiresAt *time.Time
	if claims.ExpiresAt != nil {
		ts := claims.ExpiresAt.Time.UTC()
		expiresAt = &ts
	}

	c.SetSession(Session{
		Claims:    claims,
		UserID:    claims.UserID,
		Email:     claims.Email,
		Subject:   claims.Subject,
		Issuer:    claims.Issuer,
		ExpiresAt: expiresAt,
	})
}

func mustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
