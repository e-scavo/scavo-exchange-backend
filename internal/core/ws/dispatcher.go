package ws

import (
	"context"
	"fmt"
	"sync"
)

type HandlerFunc func(ctx context.Context, c *Client, env Envelope) Envelope

type Dispatcher struct {
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[string]HandlerFunc, 32),
	}
}

func (d *Dispatcher) Handle(action string, fn HandlerFunc) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[action] = fn
}

func (d *Dispatcher) Dispatch(ctx context.Context, c *Client, env Envelope) Envelope {
	if env.Type != MsgTypeReq {
		return Envelope{
			ID:     env.ID,
			Type:   MsgTypeRes,
			Action: env.Action,
			Error:  &ErrPayload{Code: "bad_request", Msg: "type must be req"},
		}
	}

	d.mu.RLock()
	fn, ok := d.handlers[env.Action]
	d.mu.RUnlock()

	if !ok {
		return Envelope{
			ID:     env.ID,
			Type:   MsgTypeRes,
			Action: env.Action,
			Error:  &ErrPayload{Code: "unknown_action", Msg: fmt.Sprintf("unknown action: %s", env.Action)},
		}
	}

	res := fn(ctx, c, env)

	if res.Type == "" {
		res.Type = MsgTypeRes
	}
	if res.ID == "" {
		res.ID = env.ID
	}
	if res.Action == "" {
		res.Action = env.Action
	}

	return res
}
