package domain

import (
	"sync"
)

type ServerRepository interface {
	Get() *Server
}

// Server
// Goroutine safe
type Server struct {
	channels map[int64]*Channel
	clients  map[int64]*Client
	// We will simply use a RWMutex to handle read/writes to the maps
	channelsMutex sync.RWMutex
	clientsMutex  sync.RWMutex
}

func NewServer(channels map[int64]*Channel, clients map[int64]*Client) *Server {
	return &Server{
		channels: channels,
		clients:  clients,
	}
}

// AddClient adds a client to the server
func (s *Server) AddClient(client *Client) error {
	s.clientsMutex.Lock()
	s.clients[client.Id] = client
	s.clientsMutex.Unlock()
	return nil
}

// RemoveClient removes a client from the server
func (s *Server) RemoveClient(clientId int64) error {
	s.clientsMutex.Lock()
	delete(s.clients, clientId)
	s.clientsMutex.Unlock()
	return nil
}

// GetClient returns the client with clientId
func (s *Server) GetClient(clientId int64) *Client {
	s.clientsMutex.RLock()
	client := s.clients[clientId]
	s.clientsMutex.RUnlock()
	return client
}

// GetClients returns a copy of the clients map
func (s *Server) GetClients() map[int64]*Client {
	s.clientsMutex.RLock()
	clients := make(map[int64]*Client)
	for id, client := range s.clients {
		clients[id] = client
	}
	s.clientsMutex.RUnlock()
	return clients
}

// AddChannel creates a channel in the server
func (s *Server) AddChannel(channel *Channel) {
	s.channelsMutex.Lock()
	s.channels[channel.Id] = channel
	s.channelsMutex.Unlock()
}

// RemoveChannel deletes a channel in the server
func (s *Server) RemoveChannel(channelId int64) {
	s.channelsMutex.Lock()
	delete(s.channels, channelId)
	s.channelsMutex.Unlock()
}

// GetChannel returns the channel with channelId
func (s *Server) GetChannel(channelId int64) *Channel {
	s.channelsMutex.RLock()
	channel := s.channels[channelId]
	s.channelsMutex.RUnlock()
	return channel
}

// GetChannels returns a copy of the channels map
func (s *Server) GetChannels() map[int64]*Channel {
	s.channelsMutex.RLock()
	channels := make(map[int64]*Channel)
	for id, channel := range s.channels {
		channels[id] = channel
	}
	s.channelsMutex.RUnlock()
	return channels
}
