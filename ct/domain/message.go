package domain

import (
	"time"
)

type MessageRepository interface {
	Store(message *Message) error
	FindById(id int64) (*Message, error)
	FindByChannelId(channelId int64, max int) ([]*Message, error)
	FindByClientId(clientId int64, max int) ([]*Message, error)
}

// Message
type Message struct {
	Id        int64
	Body      string
	ClientId  int64
	ChannelId int64
	Author    string
	Time      time.Time
}
