package interfaces

import (
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
	SendMessage(userId int64, body string) error
}

type SendMessageRequest struct {
	Message string
}

type JoinChannelRequest struct {
	ChannelId int64
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
		ChatInteractor: ci,
	}

	ws.AddHandler(urlChat, true, wsHandler.JoinServer)

	wsocket.AddHandler("SendMessage", wsHandler.SendMessage)
	wsocket.AddHandler("JoinChannel", wsHandler.JoinChannel)

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
	request := new(SendMessageRequest)
	err := client.ReadJson(cmd, request)
	if err != nil {
		client.Error(cmd, err)
		return
	}
	user := ctx.Value(KeyUser).(*usecases.User)
	err = handler.ChatInteractor.SendMessage(user.Id, request.Message)
	if err != nil {
		client.Error(cmd, err)
		return
	}
}

func (handler *ChatWebserviceHandler) JoinChannel(ctx context.Context, client WebsocketClient, cmd WebsocketCommand) {
	request := new(JoinChannelRequest)
	err := client.ReadJson(cmd, request)
	if err != nil {
		client.Error(cmd, err)
		return
	}
	user := ctx.Value(KeyUser).(*usecases.User)
	err = handler.ChatInteractor.JoinChannel(user.Id, request.ChannelId)
	if err != nil {
		client.Error(cmd, err)
		return
	}
}

func (sender *SenderHandler) Send(msg *domain.Message) {
	sender.Command.SetType("SendMessage")
	if err := sender.Handler.SendJson(sender.Command, msg); err != nil {
		sender.Handler.Error(sender.Command, err)
	}
}
