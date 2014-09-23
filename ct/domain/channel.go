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
	Messages []*Message
}

func NewChannel(name string) *Channel {
	clients := make([]*Client, 0)
	messages := make([]*Message, 0)
	return &Channel{
		Name:     name,
		Public:   true,
		Capacity: 0,
		Clients:  clients,
		Messages: messages,
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

	// Check if the client is already in the channel
	for _, curClient := range c.Clients {
		if curClient.Id == client.Id {
			client.Channel = c
			return nil
		}
	}

	c.Clients = append(c.Clients, client)
	client.Channel = c
	return nil
}

func (c *Channel) Leave(client *Client) error {
	for i, curClient := range c.Clients {
		if curClient.Id == client.Id {
			copy(c.Clients[i:], c.Clients[i+1:])
			c.Clients[len(c.Clients)-1] = nil
			c.Clients = c.Clients[:len(c.Clients)-1]
			return nil
		}
	}
	return errors.New("Client not found")
}

func (c *Channel) Send(message *Message) {
	c.Messages = append(c.Messages, message)
	for _, client := range c.Clients {
		client.ClientSender.Send(message)
	}
	// If we have more than 100 messages delete the 50 oldest from memory
	if len(c.Messages) > 100 {
		// Nil pointers to prevent leaks
		for i := 0; i < 50; i++ {
			c.Messages[i] = nil
		}
		c.Messages = c.Messages[50:]
	}
}

func (c *Channel) HasAccess(client *Client) bool {
	if c.Public {
		return true
	}
	//TODO: Channel permissions
	return false
}
