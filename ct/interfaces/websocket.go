package interfaces

import (
	"code.google.com/p/go.net/context"
	"net/http"

	"github.com/janicduplessis/projectgo/ct/domain"
)

type Websocket interface {
	AddHandler(command string, fn func(context.Context, WebsocketClient, WebsocketCommand))
	AddClient(ctx context.Context, w http.ResponseWriter, r *http.Request, client *domain.Client)
	RemoveClient(ctx context.Context, client *domain.Client)
	Log(msg string)
}

type WebsocketClient interface {
	ReadJson(cmd WebsocketCommand, obj interface{}) error
	SendJson(cmd WebsocketCommand, obj interface{}) error
	Error(cmd WebsocketCommand, err error)
}

type WebsocketCommand interface {
	SetType(name string)
}
