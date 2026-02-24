package ws

import (
	"context"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte

	mu     sync.Mutex
	closed bool
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 128),
	}
}

func (c *Client) Close(status websocket.StatusCode, reason string) {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return
	}
	c.closed = true
	c.mu.Unlock()

	_ = c.conn.Close(status, reason)
	close(c.send)
}

func (c *Client) TrySend(msg []byte) {
	// Non-blocking send: si está lleno, se descarta para no trabar el hub.
	select {
	case c.send <- msg:
	default:
	}
}

func (c *Client) WriteLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			wctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			_ = c.conn.Write(wctx, websocket.MessageText, msg)
			cancel()
		}
	}
}
