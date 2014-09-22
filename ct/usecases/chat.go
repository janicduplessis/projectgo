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

func (ci *ChatInteractor) JoinServer(clientId int64) (*domain.Client, error) {
	server := ci.ServerRepository.Get()
	client, err := ci.ClientRepository.FindById(clientId)
	if err != nil {
		return nil, err
	}

	return client, server.Join(client)
}

func (ci *ChatInteractor) SendMessage(clientId int64, body string) error {
	server := ci.ServerRepository.Get()
	client := server.Clients[clientId]

	message := &domain.Message{
		Body:    body,
		Client:  client,
		Channel: client.Channel,
	}

	if client.Channel == nil || !client.Channel.HasAccess(client) {
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
	if channel == nil {
		return ErrInvalidChannel
	}

	if client.Channel != nil {
		client.Channel.Leave(client)
	}

	return channel.Join(client)
}

func (ci *ChatInteractor) CreateChannel(clientId int64, name string) error {
	server := ci.ServerRepository.Get()

	channel := domain.NewChannel(name)
	err := ci.ChannelRepository.Store(channel)
	if err != nil {
		return err
	}

	server.Channels[channel.Id] = channel

	for _, client := range server.Clients {
		client.ClientSender.ChannelCreated(channel)
	}

	return nil
}

func (ci *ChatInteractor) Channels(clientId int64) (map[int64]*domain.Channel, error) {
	server := ci.ServerRepository.Get()

	return server.Channels, nil
}
