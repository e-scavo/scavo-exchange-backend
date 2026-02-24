package ws

import (
	"context"

	"github.com/e-scavo/scavo-exchange-backend/internal/core/logger"
)

type Hub struct {
	log *logger.Logger

	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte

	clients map[*Client]struct{}
}

func NewHub(log *logger.Logger) *Hub {
	return &Hub{
		log:        log,
		register:   make(chan *Client, 32),
		unregister: make(chan *Client, 32),
		broadcast:  make(chan []byte, 256),
		clients:    make(map[*Client]struct{}),
	}
}

func (h *Hub) Run(ctx context.Context) {
	h.log.Info("ws hub started")
	for {
		select {
		case <-ctx.Done():
			h.log.Info("ws hub stopped")
			return
		case c := <-h.register:
			h.clients[c] = struct{}{}
			h.log.Info("ws client registered", "clients", len(h.clients))
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				h.log.Info("ws client unregistered", "clients", len(h.clients))
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				c.TrySend(msg)
			}
		}
	}
}

func (h *Hub) Register(c *Client)   { h.register <- c }
func (h *Hub) Unregister(c *Client) { h.unregister <- c }
func (h *Hub) Broadcast(msg []byte) { h.broadcast <- msg }
