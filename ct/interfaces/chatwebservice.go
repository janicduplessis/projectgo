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
	Channels(clientId int64) (map[int64]*domain.Channel, error)
	CreateChannel(clientId int64, name string) error
	SendMessage(userId int64, body string) error
}

type SendMessageRequest struct {
	Message string
}

type JoinChannelRequest struct {
	ChannelId int64
}

type CreateChannelRequest struct {
	Name string
}

type ChannelsResponse struct {
	List []ChannelModel
}

type ClientModel struct {
	Id   int64
	Name string
}

type ChannelModel struct {
	Id      int64
	Name    string
	Clients []ClientModel
}

type SendMessageResponse struct {
	Body      string
	Author    string
	ChannelId int64
	ClientId  int64
}

type ChatWebserviceHandler struct {
	Webservice     Webservice
	Websocket      Websocket
	ChatInteractor ChatInteractor
}

type SenderHandler struct {
	Handler WebsocketClient
	Command WebsocketCommand
}

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

	return wsHandler
}

func (handler *ChatWebserviceHandler) JoinServer(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value(KeyUser).(*usecases.User)

	client, err := handler.ChatInteractor.JoinServer(user.Id)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	handler.Websocket.AddClient(ctx, w, r, client)
}

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
}

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

func (handler *ChatWebserviceHandler) Channels(ctx context.Context, client WebsocketClient, cmd WebsocketCommand) {
	handler.Webservice.Log("Channels request")
	user := ctx.Value(KeyUser).(*usecases.User)
	channels, err := handler.ChatInteractor.Channels(user.Id)
	if err != nil {
		client.Error(cmd, err)
		return
	}
	// Create the response model
	channelsArr := make([]ChannelModel, len(channels))
	index := 0
	for _, curChan := range channels {
		clients := make([]ClientModel, len(curChan.Clients))
		for j, curClient := range curChan.Clients {
			clients[j] = ClientModel{
				Id:   curClient.Id,
				Name: curClient.DisplayName,
			}
		}
		channelsArr[index] = ChannelModel{
			Id:      curChan.Id,
			Name:    curChan.Name,
			Clients: clients,
		}
		index++
	}

	response := ChannelsResponse{
		List: channelsArr,
	}

	err = client.SendJson(cmd, response)
	if err != nil {
		client.Error(cmd, err)
		return
	}
}

func (sender *SenderHandler) Send(msg *domain.Message) {
	sender.Command.SetType("SendMessage")
	response := &SendMessageResponse{
		Body:      msg.Body,
		Author:    msg.Client.DisplayName,
		ChannelId: msg.Channel.Id,
		ClientId:  msg.Client.Id,
	}
	if err := sender.Handler.SendJson(sender.Command, response); err != nil {
		sender.Handler.Error(sender.Command, err)
	}
}
