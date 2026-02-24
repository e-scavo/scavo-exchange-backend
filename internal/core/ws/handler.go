package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"nhooyr.io/websocket"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
)

type HandlerParams struct {
	Log        *logger.Logger
	Hub        *Hub
	Dispatcher *Dispatcher
}

type Handler struct {
	log        *logger.Logger
	hub        *Hub
	dispatcher *Dispatcher
}

func NewHandler(p HandlerParams) http.HandlerFunc {
	h := &Handler{
		log:        p.Log,
		hub:        p.Hub,
		dispatcher: p.Dispatcher,
	}
	return h.serveWS
}

func (h *Handler) serveWS(w http.ResponseWriter, r *http.Request) {
	// TODO: auth handshake con token en query/header.

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		// En prod: validar Origin.
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

	// Mensaje hello (evt)
	hello := Envelope{
		ID:     uuid.NewString(),
		Type:   MsgTypeEvt,
		Action: "system.hello",
		Data:   JSON(map[string]any{"ts": time.Now().UTC().Format(time.RFC3339)}),
	}
	client.TrySend(mustMarshal(hello))

	// Read loop
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

		// ✅ Dispatcher
		res := h.dispatcher.Dispatch(ctx, client, env)
		client.TrySend(mustMarshal(res))
	}
}

func mustMarshal(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
