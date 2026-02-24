package ws

type Session struct {
	UserID string
	Email  string
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
