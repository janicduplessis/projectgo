package ct

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
)

const (
	urlLogin    = "/login"
	urlRegister = "/register"
	urlSend     = "/send"

	dbUser     = "ct"
	dbPassword = "wyty640"
	dbName     = "ct"
)

// Server manages login, registration and messaging
type Server struct {
	messages  []*Message
	clients   map[int]*Client
	users     map[int]*User
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
	db        *sql.DB
	auth      *Auth
}

// JSON object
type JSON map[string]interface{}

// NewServer creates an instance of server
func NewServer() *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	users := make(map[int]*User)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	// Db connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	auth := NewAuth(db)

	return &Server{
		messages:  messages,
		clients:   clients,
		users:     users,
		addCh:     addCh,
		delCh:     delCh,
		sendAllCh: sendAllCh,
		doneCh:    doneCh,
		errCh:     errCh,
		db:        db,
		auth:      auth,
	}
}

// Add adds a client
func (s *Server) Add(c *Client) {
	s.addCh <- c
}

// Del disconnect a client
func (s *Server) Del(c *Client) {
	s.delCh <- c
}

// SendAll sends a message to all clients
func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

// Done shuts down the server
func (s *Server) Done() {
	s.doneCh <- true
}

// Err logs an error
func (s *Server) Err(err error) {
	s.errCh <- err
}

// Listen starts the server
// Handle client connections and message broadcast
func (s *Server) Listen() {
	onConnected := func(conn *websocket.Conn) {
		// Whenever this function exits close the connection
		defer func() {
			err := conn.Close()
			if err != nil {
				s.errCh <- err
			}
		}()

		// Add then client to the server
		client := NewClient(conn, s)
		s.Add(client)
		// Listen untill the client disconnects
		client.Listen()
	}

	http.Handle(urlSend, websocket.Handler(onConnected))
	http.HandleFunc(urlLogin, s.handleLogin)
	http.HandleFunc(urlRegister, s.handleRegister)

	// Server listen loop
	for {
		select {

		// Add a client
		case c := <-s.addCh:
			log.Println("Adding new client")
			s.clients[c.id] = c

		// Disconnect a client
		case c := <-s.delCh:
			log.Println("Delete client")
			delete(s.clients, c.id)

		// Send to all clients
		case msg := <-s.sendAllCh:
			log.Println("Message : " + msg.Body)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		// Log errors
		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		// Stop server
		case <-s.doneCh:
			return
		}
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	var regInfo = new(RegisterInfo)
	s.readJSON(r, regInfo)

	res, err := s.auth.Register(regInfo)
	if err != nil {
		s.Err(err)
	}

	s.sendJSON(w, &JSON{"result": res})
}

func (s *Server) readJSON(r *http.Request, obj interface{}) {
	// Read the request's body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.Err(err)
	}
	err = json.Unmarshal(body, obj)
	if err != nil {
		s.Err(err)
	}
}

func (s *Server) sendJSON(w http.ResponseWriter, obj interface{}) {
	//Convert object to json
	bytes, err := json.Marshal(obj)
	if err != nil {
		s.Err(err)
	}
	//Set content type
	w.Header().Set("Content-Type", "application/json")
	//Write response body
	w.Write(bytes)
}

// Send a message to every client
func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		c.Send(msg)
	}
}
