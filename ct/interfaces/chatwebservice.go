package interfaces

import (
	"fmt"
	"net/http"

	"code.google.com/p/go.net/context"
	"github.com/janicduplessis/projectgo/ct/domain"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

const (
	urlChat string = "/ct"
)

type ChatInteractor interface {
	JoinServer(userId int64) (*domain.Client, error)
	JoinChannel(userId int64, channelId int64) error
	Channels(clientId int64) (map[int64]*domain.Channel, int64, error)
	CreateChannel(clientId int64, name string) error
	SendMessage(userId int64, body string) error
	GetMessages(channelId int64) ([]*domain.Message, error)
	Disconnect(userId int64) error
}

// ChatWebserviceHandler handles chat requests
type ChatWebserviceHandler struct {
	Webservice     Webservice
	Websocket      Websocket
	ChatInteractor ChatInteractor
}

// SenderHandler handles sends to a client
type SenderHandler struct {
	Handler WebsocketClient
	Command WebsocketCommand
}

// NewChatWebservice ctor
func NewChatWebservice(ws Webservice, wsocket Websocket, ci ChatInteractor) *ChatWebserviceHandler {
	wsHandler := &ChatWebserviceHandler{
		Webservice:     ws,
		Websocket:      wsocket,
		ChatInteractor: ci,
	}
	ws.AddHandler(urlChat, true, wsHandler.JoinServer)

	wsocket.AddHandler("SendMessage", wsHandler.SendMessage)
	wsocket.AddHandler("JoinChannel", wsHandler.JoinChannel)
	wsocket.AddHandler("CreateChannel", wsHandler.CreateChannel)
	wsocket.AddHandler("Channels", wsHandler.Channels)
	wsocket.AddHandler("Ping", wsHandler.Ping)

	return wsHandler
}

// JoinServer handles a join server request
func (handler *ChatWebserviceHandler) JoinServer(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value(KeyUser).(*usecases.User)

	client, err := handler.ChatInteractor.JoinServer(user.Id)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	handler.Websocket.AddClient(ctx, w, r, client)
}

// SendMessage handles a send message request
func (handler *ChatWebserviceHandler) SendMessage(ctx context.Context, client WebsocketClient, cmd WebsocketCommand) {
	request := SendMessageRequest{}
	err := client.ReadJson(cmd, &request)
	if err != nil {
		client.Error(cmd, err)
		return
	}

	user := ctx.Value(KeyUser).(*usecases.User)

	handler.Webservice.Log(fmt.Sprintf("Sending message by id %d with content %s", user.Id, request.Message))

	err = handler.ChatInteractor.SendMessage(user.Id, request.Message)
	if err != nil {
		client.Error(cmd, err)
		return
	}
}

// JoinChannel handles a join channel request
func (handler *ChatWebserviceHandler) JoinChannel(ctx context.Context, client WebsocketClient, cmd WebsocketCommand) {
	handler.Webservice.Log("Join channel request")
	request := JoinChannelRequest{}
	err := client.ReadJson(cmd, &request)
	if err != nil {
		client.Error(cmd, err)
		return
	}
	user := ctx.Value(KeyUser).(*usecases.User)

	handler.Webservice.Log(fmt.Sprintf("Joining channel with id %d for user %d", request.ChannelId, user.Id))

	err = handler.ChatInteractor.JoinChannel(user.Id, request.ChannelId)
	if err != nil {
		client.Error(cmd, err)
		return
	}

	messages, err := handler.ChatInteractor.GetMessages(request.ChannelId)
	if err != nil {
		client.Error(cmd, err)
		return
	}

	messagesModel := make([]MessageModel, len(messages))
	for i, curMessage := range messages {
		messagesModel[i] = MessageModel{
			Author:    curMessage.Author,
			Body:      curMessage.Body,
			UnixTime:  curMessage.Time.Unix(),
			ChannelId: curMessage.ChannelId,
			ClientId:  curMessage.ClientId,
		}
	}

	response := JoinChannelResponse{
		Messages: messagesModel,
		Result:   true,
	}

	err = client.SendJson(cmd, response)
	if err != nil {
		client.Error(cmd, err)
		return
	}
}

// CreateChannel handles a create channel request
func (handler *ChatWebserviceHandler) CreateChannel(ctx context.Context, client WebsocketClient, cmd WebsocketCommand) {
	handler.Webservice.Log("Create channel request")
	request := CreateChannelRequest{}
	err := client.ReadJson(cmd, &request)
	if err != nil {
		client.Error(cmd, err)
		return
	}
	user := ctx.Value(KeyUser).(*usecases.User)
	handler.Webservice.Log(fmt.Sprintf("Creating channel with name %s for user %d", request.Name, user.Id))
	err = handler.ChatInteractor.CreateChannel(user.Id, request.Name)
	if err != nil {
		client.Error(cmd, err)
		return
	}
}

// Channels handles a channels request, returning the list of all channels
func (handler *ChatWebserviceHandler) Channels(ctx context.Context, client WebsocketClient, cmd WebsocketCommand) {
	handler.Webservice.Log("Channels request")
	user := ctx.Value(KeyUser).(*usecases.User)
	channels, curChannel, err := handler.ChatInteractor.Channels(user.Id)
	if err != nil {
		client.Error(cmd, err)
		return
	}
	// Create the response model
	channelsArr := make([]*ChannelModel, len(channels))
	index := 0
	for _, curChan := range channels {
		channelsArr[index] = createChannelModel(curChan)
		index++
	}

	response := ChannelsResponse{
		List:    channelsArr,
		Current: curChannel,
	}

	err = client.SendJson(cmd, response)
	if err != nil {
		client.Error(cmd, err)
		return
	}
}

// Ping handles client heartbeats
func (handler *ChatWebserviceHandler) Ping(ctx context.Context, client WebsocketClient, cmd WebsocketCommand) {
	client.SendJson(cmd, &PingModel{Response: "Pong"})
}

// SenderHandler implementation

// Send handles sending a message to a client
func (sender *SenderHandler) Send(msg *domain.Message) {
	sender.Command.SetType("SendMessage")
	response := &MessageModel{
		Author:    msg.Author,
		Body:      msg.Body,
		UnixTime:  msg.Time.Unix(),
		ChannelId: msg.ChannelId,
		ClientId:  msg.ClientId,
	}
	if err := sender.Handler.SendJson(sender.Command, response); err != nil {
		sender.Handler.Error(sender.Command, err)
	}
}

// ChannelCreated handles warning the client about a new channel
func (sender *SenderHandler) ChannelCreated(channel *domain.Channel) {
	sender.Command.SetType("CreateChannel")
	response := createChannelModel(channel)
	if err := sender.Handler.SendJson(sender.Command, response); err != nil {
		sender.Handler.Error(sender.Command, err)
	}
}

// ChannelJoined handles warning the client about someone joining a channel
func (sender *SenderHandler) ChannelJoined(channel *domain.Channel, client *domain.Client) {
	sender.Command.SetType("ChannelJoined")
	response := &ChannelJoinedResponse{
		ChannelId: channel.Id,
		Client: ClientModel{
			Id:   client.Id,
			Name: client.DisplayName,
		},
	}
	if err := sender.Handler.SendJson(sender.Command, response); err != nil {
		sender.Handler.Error(sender.Command, err)
	}
}

// Utils
func createChannelModel(channel *domain.Channel) *ChannelModel {
	clients := make([]ClientModel, len(channel.Clients))
	for i, curClient := range channel.Clients {
		clients[i] = ClientModel{
			Id:   curClient.Id,
			Name: curClient.DisplayName,
		}
	}
	return &ChannelModel{
		Id:      channel.Id,
		Name:    channel.Name,
		Clients: clients,
	}
}
