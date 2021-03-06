package usecases

import (
	"fmt"
	"html/template"
	"regexp"
	"time"

	"github.com/janicduplessis/projectgo/ct/domain"
)

type ChatInteractor struct {
	ServerRepository  domain.ServerRepository
	ChannelRepository domain.ChannelRepository
	MessageRepository domain.MessageRepository
	ClientRepository  domain.ClientRepository
	Logger            Logger

	urlRegex *regexp.Regexp
}

func NewChatInteractor(serverRepo domain.ServerRepository, channelRepo domain.ChannelRepository,
	msgRepo domain.MessageRepository, clientRepo domain.ClientRepository, logger Logger) *ChatInteractor {

	urlRegex := regexp.MustCompile(`(https?:\/\/([-\w\.]+)+(:\d+)?(\/([\w\/_\.]*(\?\S+)?)?)?)`)

	return &ChatInteractor{
		ServerRepository:  serverRepo,
		ChannelRepository: channelRepo,
		MessageRepository: msgRepo,
		ClientRepository:  clientRepo,
		Logger:            logger,

		urlRegex: urlRegex,
	}
}

func (ci *ChatInteractor) JoinServer(clientId int64) (*domain.Client, error) {
	server := ci.ServerRepository.Get()

	// If the client is already connected on the server
	client := server.GetClient(clientId)
	if client != nil {
		return client, nil
	}

	// If not get its info
	client, err := ci.ClientRepository.FindById(clientId)
	if err != nil {
		return nil, err
	}

	return client, server.AddClient(client)
}

func (ci *ChatInteractor) SendMessage(clientId int64, body string) error {
	server := ci.ServerRepository.Get()
	client := server.GetClient(clientId)

	if client.Channel == nil {
		return ErrNoChannel
	}

	safeBody := template.HTMLEscapeString(body)

	// Create hyperlinks
	safeBody = ci.urlRegex.ReplaceAllString(safeBody, `<a href="$1" target="_blank">$1</a>`)

	message := &domain.Message{
		Body:      safeBody,
		ClientId:  client.Id,
		Author:    client.DisplayName,
		Time:      time.Now(),
		ChannelId: client.Channel.Id,
	}

	if client.Channel == nil || !client.Channel.HasAccess(client) {
		return ErrAccessDenied
	}

	client.Channel.Send(message)

	return ci.MessageRepository.Store(message)
}

func (ci *ChatInteractor) JoinChannel(clientId int64, channelId int64) error {
	server := ci.ServerRepository.Get()
	client := server.GetClient(clientId)
	channel := server.GetChannel(channelId)
	if channel == nil {
		return ErrInvalidChannelId
	}

	if client.Channel != nil {
		ci.Logger.Log(fmt.Sprintf("Client %s is in channel %s", client.DisplayName, client.Channel.Name))
		// If the client is already in the channel we have nothing to do
		if client.Channel.Id == channelId {
			ci.Logger.Log("Client already in the channel")
			return nil
		}
		// If the client is in a channel leave it
		client.Channel.Leave(client)
	}

	// Try to join...
	if err := channel.Join(client); err != nil {
		return err
	}

	// Alert other clients
	for _, c := range server.GetClients() {
		c.ClientSender.ChannelJoined(channel, client)
	}

	return nil
}

func (ci *ChatInteractor) CreateChannel(clientId int64, name string) error {
	server := ci.ServerRepository.Get()

	channel := domain.NewChannel(name)
	err := ci.ChannelRepository.Store(channel)
	if err != nil {
		return err
	}

	server.AddChannel(channel)

	for _, client := range server.GetClients() {
		client.ClientSender.ChannelCreated(channel)
	}

	return nil
}

func (ci *ChatInteractor) Channels(clientId int64) (map[int64]*domain.Channel, int64, error) {
	server := ci.ServerRepository.Get()
	client := server.GetClient(clientId)
	var channelId int64 = -1
	if client == nil {
		return nil, channelId, ErrInvalidClientId
	}
	if client.Channel != nil {
		channelId = client.Channel.Id
	}

	return server.GetChannels(), channelId, nil
}

func (ci *ChatInteractor) Disconnect(clientId int64) error {
	server := ci.ServerRepository.Get()
	client := server.GetClient(clientId)
	if client == nil {
		return ErrInvalidChannelId
	}
	server.RemoveClient(clientId)
	if client != nil && client.Channel != nil {
		client.Channel.Leave(client)
	}

	return nil
}

func (ci *ChatInteractor) GetMessages(channelId int64) ([]*domain.Message, error) {
	server := ci.ServerRepository.Get()
	channel := server.GetChannel(channelId)
	if channel == nil {
		return nil, ErrInvalidChannelId
	}

	return channel.Messages, nil
}
