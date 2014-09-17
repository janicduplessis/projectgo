package infrastructure

import (
	"errors"
	"io"
	"net/http"
	"reflect"

	"code.google.com/p/go.net/context"
	"code.google.com/p/go.net/websocket"

	"github.com/janicduplessis/projectgo/ct/domain"
	"github.com/janicduplessis/projectgo/ct/interfaces"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

type WebsocketHandler struct {
	Logger   usecases.Logger
	Handlers map[string]func(context.Context, interfaces.WebsocketClient, interfaces.WebsocketCommand)
	Clients  map[int64]*WebsocketClientHandler
}

type WebsocketClientHandler struct {
	WebsocketHandler *WebsocketHandler
	Context          context.Context

	conn *websocket.Conn

	doneCh    chan bool
	commandCh chan *WebsocketCommandHandler
}

type WebsocketClientError struct {
	Message string
}

type WebsocketCommandHandler struct {
	Command string
	Data    interface{}
	client  *WebsocketClientHandler
}

func NewWebsocketHandler(logger usecases.Logger) *WebsocketHandler {
	handlers := make(map[string]func(context.Context, interfaces.WebsocketClient, interfaces.WebsocketCommand))
	clients := make(map[int64]*WebsocketClientHandler)

	return &WebsocketHandler{
		Handlers: handlers,
		Clients:  clients,
		Logger:   logger,
	}
}

func (handler *WebsocketHandler) AddHandler(commandName string, fn func(context.Context, interfaces.WebsocketClient, interfaces.WebsocketCommand)) {
	handler.Handlers[commandName] = fn
}

func (handler *WebsocketHandler) AddClient(ctx context.Context, w http.ResponseWriter, r *http.Request, client *domain.Client) {
	onConnected := func(conn *websocket.Conn) {
		// Whenever this function exits close the connection
		defer func() {
			err := conn.Close()
			if err != nil {
				handler.Logger.Error(err)
			}
		}()

		// Add then client to the server
		doneCh := make(chan bool)
		commandCh := make(chan *WebsocketCommandHandler)
		clientHandler := &WebsocketClientHandler{
			WebsocketHandler: handler,
			Context:          ctx,
			conn:             conn,
			doneCh:           doneCh,
			commandCh:        commandCh,
		}
		command := &WebsocketCommandHandler{
			client: clientHandler,
		}
		sender := &interfaces.SenderHandler{
			Handler: clientHandler,
			Command: command,
		}
		client.ClientSender = sender
		handler.Clients[client.Id] = clientHandler

		// Listen untill the client disconnects
		clientHandler.listen()
	}

	onConnectedHander := websocket.Handler(onConnected)
	onConnectedHander.ServeHTTP(w, r)
}

func (handler *WebsocketHandler) RemoveClient(ctx context.Context, client *domain.Client) {
	//TODO: NYI
}

func (handler *WebsocketHandler) Log(msg string) {
	handler.Logger.Log(msg)
}

func (handler *WebsocketHandler) executeCommand(cmd *WebsocketCommandHandler) {
	fn := handler.Handlers[cmd.Command]
	if fn == nil {
		handler.Logger.Log("No handler for requested command")
		return
	}

	fn(cmd.client.Context, cmd.client, cmd)
}

func (client *WebsocketClientHandler) ReadJson(cmd interfaces.WebsocketCommand, obj interface{}) error {
	cmdData, ok := cmd.(*WebsocketCommandHandler)
	if !ok {
		return errors.New("Invalid command")
	}

	// Need to use reflect to copy values into obj
	dest := reflect.ValueOf(obj).Elem()
	src, ok := cmdData.Data.(map[string]interface{})
	if !ok {
		return errors.New("Invalid command")
	}
	for srcKey, srcVal := range src {
		destF := dest.FieldByName(srcKey)
		srcF := reflect.ValueOf(srcVal)
		if destF.IsValid() && srcF.Type().ConvertibleTo(destF.Type()) {
			destF.Set(srcF.Convert(destF.Type()))
		}
	}

	return nil
}

func (client *WebsocketClientHandler) SendJson(cmd interfaces.WebsocketCommand, obj interface{}) error {
	cmdData, ok := cmd.(*WebsocketCommandHandler)
	if !ok {
		return errors.New("Invalid command")
	}
	cmdData.Data = obj
	client.commandCh <- cmdData
	return nil
}

func (client *WebsocketClientHandler) Error(cmd interfaces.WebsocketCommand, err error) {
	cmd.SetType("ERROR")
	client.WebsocketHandler.Logger.Error(err)
	client.SendJson(cmd, &WebsocketClientError{Message: "INTERNAL_SERVER_ERROR"})
}

func (client *WebsocketClientHandler) listen() {
	go client.listenSend()
	client.listenReceive()
}

func (client *WebsocketClientHandler) Done() {
	client.doneCh <- true
}

func (client *WebsocketClientHandler) listenSend() {
	for {
		select {
		case command := <-client.commandCh:
			err := websocket.JSON.Send(client.conn, command)
			if err != nil {
				client.Error(command, err)
			}
		}
	}
}

func (client *WebsocketClientHandler) listenReceive() {
	for {
		select {
		case <-client.doneCh:
			return
		default:
			command := new(WebsocketCommandHandler)
			err := websocket.JSON.Receive(client.conn, command)
			if err == io.EOF {
				client.Done()
			} else if err != nil {
				client.Error(command, err)
			} else {
				command.client = client
				go client.WebsocketHandler.executeCommand(command)
			}
		}
	}
}

func (cmd *WebsocketCommandHandler) SetType(name string) {
	cmd.Command = name
}