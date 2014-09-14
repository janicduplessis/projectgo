package domain

type MessageRepository interface {
	Store(message *Message) error
	FindById(id int64) (*Message, error)
}

type Message struct {
	Id      int64
	Body    string
	Client  *Client
	Channel *Channel
}
