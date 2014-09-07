package domain

import (
	"errors"
)

type ChannelRepository interface {
	Store(channel *Channel)
	FindById(id int64) *Channel
}

// Chat channel
type Channel struct {
	Id       int64
	Name     string
	Public   bool
	Capacity int
	Clients  []*Client
}

// Add adds a client to the channel
func (c *Channel) Add(client *Client) error {
	if !c.hasAccess(client) {
		return errors.New("Client cannot join this channel")
	}

	if c.Capacity != 0 && len(c.Clients) >= c.Capacity {
		return errors.New("Channel is full")
	}

	c.Clients = append(c.Clients, client)
	return nil
}

func (c *Channel) hasAccess(client *Client) bool {
	if c.Public {
		return true
	}
	//TODO: Channel permissions
	return false
}
