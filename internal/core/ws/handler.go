package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
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

	// ✅ auth opcional: setea session si token válido
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

	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if token == "" {
		authz := r.Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(authz), "bearer ") {
			token = strings.TrimSpace(authz[7:])
		}
	}
	if token == "" {
		return
	}

	claims, err := h.tokenSvc.Parse(token)
	if err != nil || claims == nil || claims.UserID == "" {
		return
	}

	c.SetSession(Session{
		UserID: claims.UserID,
		Email:  claims.Email,
	})
}

func mustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
