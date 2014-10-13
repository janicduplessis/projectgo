package domain

type ClientRepository interface {
	Store(client *Client) error
	FindById(id int64) (*Client, error)
}

// ClientSender is used to send information to the client
type ClientSender interface {
	Send(message *Message)
	ChannelCreated(channel *Channel)
	ChannelJoined(channel *Channel, client *Client)
}

// Client is a normal user of the chat service
type Client struct {
	Id          int64
	DisplayName string
	FirstName   string
	LastName    string
	Email       string

	Channel      *Channel
	ClientSender ClientSender
}
