package ws

import (
	"log"

	"golang.org/x/net/websocket"
)

// Handler handles client communication.
type Handler struct {
	// clients holds all current clients.
	Clients map[*Client]bool

	// join is a channel for clients wishing to join.
	join chan *Client
	// leave is a channel for clients wishing to leave.
	leave chan *Client

	// Send message to clients.
	send chan *Message

	//Logger logs information.
	Logger *log.Logger
}

const (
	messageBufferSize = 256
)

// WSMockInfoServer echoes mock responses back to gui.
func (ch *Handler) WSMockInfoServer(ws *websocket.Conn) {

	client := &Client{
		Ws:      ws,
		Send:    make(chan *Message, messageBufferSize),
		handler: ch,
	}
	ch.join <- client
	defer func() { ch.leave <- client }()
	client.write()
}

// NewHandler creates a new client handler.
func NewHandler(log *log.Logger) *Handler {
	return &Handler{
		join:    make(chan *Client),
		leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
		send:    make(chan *Message),
		Logger:  log}
}

// Send message to client.
func (ch *Handler) Send(msg *Message) {
	ch.send <- msg
}

// Run should be used as a goroutine to handle clients.
func (ch *Handler) Run() {
	for {
		select {
		case client := <-ch.join:
			// joining
			ch.Clients[client] = true
			ch.Logger.Println("New client joined!")
		case client := <-ch.leave:
			// leaving
			delete(ch.Clients, client)
			close(client.Send)
			ch.Logger.Println("Client left!")
		case msg := <-ch.send:
			// forward message to all clients
			for client := range ch.Clients {
				select {
				case client.Send <- msg:
				default:
					// failed to send
					delete(ch.Clients, client)
					close(client.Send)
					ch.Logger.Println("-- failed to send, clean up client")
				}
			}
		}
	}
}
