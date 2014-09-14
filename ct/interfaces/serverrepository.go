package interfaces

import (
	"github.com/janicduplessis/projectgo/ct/domain"
)

var server *domain.Server

type SingletonServerRepo struct{}

func NewSingletonServerRepo() *SingletonServerRepo {
	return new(SingletonServerRepo)
}

func (repo *SingletonServerRepo) Get() *domain.Server {
	if server == nil {
		channels := make(map[int64]*domain.Channel)
		clients := make(map[int64]*domain.Client)

		server = &domain.Server{
			Channels: channels,
			Clients:  clients,
		}
	}

	return server
}
