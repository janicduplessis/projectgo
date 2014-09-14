package usecases

import (
	"github.com/janicduplessis/projectgo/ct/domain"
)

type ChatInteractor struct {
	ServerRepository  domain.ServerRepository
	ChannelRepository domain.ChannelRepository
	MessageRepository domain.MessageRepository
	ClientRepository  domain.ClientRepository
	Logger            Logger
}

func (ci *ChatInteractor) JoinServer(clientId int64) error {
	server := ci.ServerRepository.Get()
	client, err := ci.ClientRepository.FindById(clientId)
	if err != nil {
		return err
	}
	return server.Join(client)
}

func (ci *ChatInteractor) SendMessage(clientId int64, body string) error {
	server := ci.ServerRepository.Get()
	client := server.Clients[clientId]

	message := &domain.Message{
		Body:    body,
		Client:  client,
		Channel: client.Channel,
	}

	if !client.Channel.HasAccess(client) {
		return ErrAccessDenied
	}

	client.Channel.Send(message)

	go ci.MessageRepository.Store(message)

	return nil
}

func (ci *ChatInteractor) JoinChannel(clientId int64, channelId int64) error {
	server := ci.ServerRepository.Get()
	client := server.Clients[clientId]
	channel := server.Channels[channelId]

	return channel.Join(client)
}
