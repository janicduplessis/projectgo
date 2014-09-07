package domain

type MessageRepository interface {
	Store(message *Message)
	FindById(id int64) *Message
}

type Message struct {
	Id      int64
	Author  string
	Body    string
	Channel *Channel
}
