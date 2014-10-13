package infrastructure

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"

	"code.google.com/p/go.net/context"
	"code.google.com/p/go.net/websocket"

	"github.com/janicduplessis/projectgo/ct/domain"
	"github.com/janicduplessis/projectgo/ct/interfaces"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

// WebsocketHandler handleS websocket connections
type WebsocketHandler struct {
	Logger   usecases.Logger
	Handlers map[string]func(context.Context, interfaces.WebsocketClient, interfaces.WebsocketCommand)
	Clients  map[int64]*WebsocketClientHandler

	mutex sync.RWMutex
}

// WebsocketClientHandler handles connections for a client
type WebsocketClientHandler struct {
	WebsocketHandler *WebsocketHandler
	Context          context.Context

	id       int64
	chsMutex sync.Mutex

	// 1 channel per connection
	doneChs    []chan bool
	commandChs []chan *WebsocketCommandHandler
}

// WebsocketClientError error model
type WebsocketClientError struct {
	Message string
}

// WebsocketCommandHandler command data
type WebsocketCommandHandler struct {
	Command string
	Data    interface{}
	client  *WebsocketClientHandler
}

// WebsocketHandler impl

// NewWebsocketHandler ctor
func NewWebsocketHandler(logger usecases.Logger) *WebsocketHandler {
	handlers := make(map[string]func(context.Context, interfaces.WebsocketClient, interfaces.WebsocketCommand))
	clients := make(map[int64]*WebsocketClientHandler)

	return &WebsocketHandler{
		Handlers: handlers,
		Clients:  clients,
		Logger:   logger,
	}
}

// AddHandler adds an handler for a command
func (handler *WebsocketHandler) AddHandler(commandName string, fn func(context.Context, interfaces.WebsocketClient, interfaces.WebsocketCommand)) {
	handler.Handlers[commandName] = fn
}

// AddClient connects a client
func (handler *WebsocketHandler) AddClient(ctx context.Context, w http.ResponseWriter, r *http.Request, client *domain.Client) {
	onConnected := func(conn *websocket.Conn) {
		// Whenever this function exits close the connection
		defer func() {
			err := conn.Close()
			if err != nil {
				handler.Logger.Error(err)
			}
		}()

		handler.mutex.RLock()
		clientHandler := handler.Clients[client.Id]
		handler.mutex.RUnlock()
		// If we dont have an handler yet create it
		if clientHandler == nil {
			clientHandler = &WebsocketClientHandler{
				WebsocketHandler: handler,
				Context:          ctx,
				id:               client.Id,
			}
			command := &WebsocketCommandHandler{
				client: clientHandler,
			}
			sender := &interfaces.SenderHandler{
				Handler: clientHandler,
				Command: command,
			}
			client.ClientSender = sender
			handler.mutex.Lock()
			handler.Clients[client.Id] = clientHandler
			handler.mutex.Unlock()
		}

		// Create channels for this connection
		commandCh := make(chan *WebsocketCommandHandler)
		doneCh := make(chan bool)
		clientHandler.chsMutex.Lock()
		clientHandler.commandChs = append(clientHandler.commandChs, commandCh)
		clientHandler.doneChs = append(clientHandler.doneChs, doneCh)
		clientHandler.chsMutex.Unlock()

		// Listen using the new connection
		clientHandler.listen(conn, commandCh, doneCh)
	}

	onConnectedHander := websocket.Handler(onConnected)
	onConnectedHander.ServeHTTP(w, r)
}

func (handler *WebsocketHandler) Log(message string) {
	handler.Logger.Log(message)
}

func (handler *WebsocketHandler) executeCommand(cmd *WebsocketCommandHandler) {
	fn := handler.Handlers[cmd.Command]
	if fn == nil {
		handler.Logger.Log("No handler for requested command")
		return
	}

	fn(cmd.client.Context, cmd.client, cmd)
}

// WebsocketClientHandler impl

// ReadJson parses data from a command
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

// SendJson sends a command to the client
func (client *WebsocketClientHandler) SendJson(cmd interfaces.WebsocketCommand, obj interface{}) error {
	cmdData, ok := cmd.(*WebsocketCommandHandler)
	if !ok {
		return errors.New("Invalid command")
	}
	cmdData.Data = obj
	for _, channel := range client.commandChs {
		channel <- cmdData
	}
	return nil
}

// Error sends an error to the client
func (client *WebsocketClientHandler) Error(cmd interfaces.WebsocketCommand, err error) {
	cmd.SetType("Error")
	client.WebsocketHandler.Logger.Error(err)
	client.SendJson(cmd, &WebsocketClientError{Message: "INTERNAL_SERVER_ERROR"})
}

func (client *WebsocketClientHandler) listen(conn *websocket.Conn, commandCh chan *WebsocketCommandHandler, doneCh chan bool) {
	go client.listenSend(conn, commandCh, doneCh)
	client.listenReceive(conn, commandCh, doneCh)
	client.done(commandCh, doneCh)
}

func (client *WebsocketClientHandler) done(commandCh chan *WebsocketCommandHandler, doneCh chan bool) {
	doneCh <- true
	client.chsMutex.Lock()
	for i, channel := range client.doneChs {
		if doneCh == channel {
			client.doneChs[i], client.doneChs = client.doneChs[len(client.doneChs)-1], client.doneChs[:len(client.doneChs)-1]
		}
	}
	for i, channel := range client.commandChs {
		if commandCh == channel {
			client.commandChs[i], client.commandChs = client.commandChs[len(client.commandChs)-1], client.commandChs[:len(client.commandChs)-1]
		}
	}
	if len(client.doneChs) == 0 {
		client.WebsocketHandler.mutex.Lock()
		delete(client.WebsocketHandler.Clients, client.id)
		client.WebsocketHandler.mutex.Unlock()
	}
	client.chsMutex.Unlock()
}

func (client *WebsocketClientHandler) listenSend(conn *websocket.Conn, commandCh chan *WebsocketCommandHandler, doneCh chan bool) {
	for {
		select {
		case <-doneCh:
			return
		case command := <-commandCh:
			client.WebsocketHandler.Logger.Log(fmt.Sprintf("Sending command %s for connection %s", command.Command, conn.RemoteAddr().String()))
			err := websocket.JSON.Send(conn, command)
			if err != nil {
				client.Error(command, err)
			}
		}
	}
}

func (client *WebsocketClientHandler) listenReceive(conn *websocket.Conn, commandCh chan *WebsocketCommandHandler, doneCh chan bool) {
	for {
		select {
		case <-doneCh:
			return
		default:
			command := new(WebsocketCommandHandler)
			err := websocket.JSON.Receive(conn, command)
			if err == io.EOF {
				client.done(commandCh, doneCh)
			} else if err != nil {
				client.Error(command, err)
			} else {
				client.WebsocketHandler.Logger.Log(fmt.Sprintf("Received command %s for connection %s", command.Command, conn.RemoteAddr().String()))
				command.client = client
				go client.WebsocketHandler.executeCommand(command)
			}
		}
	}
}

// WebsocketCommandHandler impl

// SetType sets the command name
func (cmd *WebsocketCommandHandler) SetType(name string) {
	cmd.Command = name
}
