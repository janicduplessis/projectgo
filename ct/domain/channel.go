package domain

import (
	"errors"
)

type ChannelRepository interface {
	Channels() ([]*Channel, error)
	Store(channel *Channel) error
	FindById(id int64) (*Channel, error)
}

// Chat channel
type Channel struct {
	Id       int64
	Name     string
	Public   bool
	Capacity int
	Clients  []*Client
}

func NewChannel(name string) *Channel {
	clients := make([]*Client, 0)
	return &Channel{
		Name:     name,
		Public:   true,
		Capacity: 0,
		Clients:  clients,
	}
}

// Add adds a client to the channel
func (c *Channel) Join(client *Client) error {
	if !c.HasAccess(client) {
		return errors.New("Client cannot join this channel")
	}

	if c.Capacity != 0 && len(c.Clients) >= c.Capacity {
		return errors.New("Channel is full")
	}

	c.Clients = append(c.Clients, client)
	client.Channel = c
	return nil
}

func (c *Channel) Send(message *Message) {
	for _, client := range c.Clients {
		client.ClientSender.Send(message)
	}
}

func (c *Channel) HasAccess(client *Client) bool {
	if c.Public {
		return true
	}
	//TODO: Channel permissions
	return false
}
