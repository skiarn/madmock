package ws

import "golang.org/x/net/websocket"

// Client user.
type Client struct {
	// ClientId
	ClientID string
	// ws is the web socket for this client.
	Ws *websocket.Conn
	// send is a channel on which messages are sent.
	Send chan *Message

	handler *Handler
}

func (c *Client) write() {
	for msg := range c.Send {
		if _, ew := c.Ws.Write(msg.Data); ew != nil {
			break
		}
	}
	c.Ws.Close()
}
