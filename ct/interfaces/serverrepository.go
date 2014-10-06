package interfaces

import (
	"github.com/janicduplessis/projectgo/ct/domain"
)

var server *domain.Server

type SingletonServerRepo struct {
	dbHandlers map[string]DbHandler
}

func NewSingletonServerRepo(dbHandlers map[string]DbHandler) *SingletonServerRepo {
	return &SingletonServerRepo{
		dbHandlers: dbHandlers,
	}
}

func (repo *SingletonServerRepo) Get() *domain.Server {
	if server == nil {
		chanRepo := NewDbChannelRepo(repo.dbHandlers)
		messageRepo := NewDbMessageRepo(repo.dbHandlers)
		// Get channels from db
		channels, err := chanRepo.GetChannels()
		if err != nil {
			panic(err)
		}

		// Get last 50 messages in each channel
		for _, channel := range channels {
			messages, err := messageRepo.FindByChannelId(channel.Id, 50)
			if err != nil {
				panic(err)
			}
			channel.Messages = messages
		}

		clients := make(map[int64]*domain.Client)

		server = domain.NewServer(channels, clients)
	}

	return server
}
