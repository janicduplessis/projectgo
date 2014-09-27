package domain

type ServerRepository interface {
	Get() *Server
}

// Server
type Server struct {
	Channels map[int64]*Channel
	Clients  map[int64]*Client
}

// Join adds a client to the server
func (s *Server) Join(client *Client) error {
	s.Clients[client.Id] = client
	return nil
}

// Leave removes a client from the server
func (s *Server) Leave(client *Client) error {
	delete(s.Clients, client.Id)
	return nil
}

// AddChannel creates a channel in the server
func (s *Server) AddChannel(channel *Channel) {
	s.Channels[channel.Id] = channel
}
