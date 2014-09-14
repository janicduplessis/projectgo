package domain

type ServerRepository interface {
	Get() *Server
}

type Server struct {
	Channels map[int64]*Channel
	Clients  map[int64]*Client
}

func (s *Server) Join(client *Client) error {
	s.Clients[client.Id] = client
	return nil
}

func (s *Server) Leave(client *Client) error {
	delete(s.Clients, client.Id)
	return nil
}

func (s *Server) AddChannel(channel *Channel) {
	s.Channels[channel.Id] = channel
}
