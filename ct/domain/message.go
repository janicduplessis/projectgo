package domain

type MessageRepository interface {
	Store(message *Message) error
	FindById(id int64) (*Message, error)
}

// Message
type Message struct {
	Id       int64
	Body     string
	ClientId int64
	Author   string
	Channel  *Channel
}
