package domain

type ServerRepository interface {
	Get() *Server
}

type Server struct {
	Channels map[int64]*Channel
	Clients  map[int64]*Client
}

func (s *Server) Join(client *Client) {
	s.Clients[client.Id] = client
}

func (s *Server) Leave(client *Client) {
	delete(s.Clients, client.Id)
}

func (s *Server) AddChannel(channel *Channel) {
	s.Channels[channel.Id] = channel
}
