package ct

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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
	logCh     chan string
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
	logCh := make(chan string)

	// Db connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if err != nil {
		panic(err)
	}

	auth := NewAuth(db)

	return &Server{
		messages:  messages,
		clients:   clients,
		users:     users,
		addCh:     addCh,
		delCh:     delCh,
		sendAllCh: sendAllCh,
		doneCh:    doneCh,
		logCh:     logCh,
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
	s.Log("d", fmt.Sprintf("Received message from %s: %s", msg.Author, msg.Body))
	s.sendAllCh <- msg
}

// Done shuts down the server
func (s *Server) Done() {
	s.doneCh <- true
}

// Err logs an error
func (s *Server) Err(err error) {
	s.Log("e", err.Error())
}

func (s *Server) Log(tag string, message string) {
	t := time.Now()
	s.logCh <- fmt.Sprintf("[%s][%s] %s", t.Format("2006-01-02 15:04:05"), tag, message)
}

// Listen starts the server
// Handle client connections and message broadcast
func (s *Server) Listen() {
	onConnected := func(conn *websocket.Conn) {
		// Whenever this function exits close the connection
		defer func() {
			err := conn.Close()
			if err != nil {
				s.Err(err)
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

	defer s.db.Close()

	// Server listen loop
	for {
		select {

		// Add a client
		case c := <-s.addCh:
			s.clients[c.id] = c

		// Disconnect a client
		case c := <-s.delCh:
			delete(s.clients, c.id)

		// Send to all clients
		case msg := <-s.sendAllCh:
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		// Log errors
		case message := <-s.logCh:
			log.Println(message)

		// Stop server
		case <-s.doneCh:
			return
		}
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	s.Log("d", "Login request")
	var loginInfo = new(LoginInfo)
	s.readJSON(r, loginInfo)

	user, err := s.auth.Login(loginInfo)
	if err != nil {
		s.Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var result bool
	if user != nil {
		//Authentification successful
		result = true
	} else {
		//Authentification unsuccessful
		result = false
	}

	s.sendJSON(w, &JSON{"Result": result, "User": user})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	s.Log("d", "Register request")
	var regInfo = new(RegisterInfo)
	s.readJSON(r, regInfo)

	res, err := s.auth.Register(regInfo)
	if err != nil {
		s.Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.sendJSON(w, &JSON{"Result": res})
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
