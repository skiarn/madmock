package ws

//Message represents a single message.
type Message struct {
	Data   []byte
	Sender *Client
}
