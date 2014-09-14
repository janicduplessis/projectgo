package domain

type ClientRepository interface {
	Store(client *Client) error
	FindById(id int64) (*Client, error)
}

type ClientSender interface {
	Send(message *Message)
}

type Client struct {
	Id          int64
	DisplayName string
	FirstName   string
	LastName    string
	Email       string

	Channel      *Channel
	ClientSender ClientSender
}
