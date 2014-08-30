package ct

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/sessions"
)

const (
	urlLogin    = "/login"
	urlRegister = "/register"
	urlSend     = "/send"
	urlLogout   = "/logout"

	urlProfileModel = "/models/getProfileModel"

	dbUser     = "ct"
	dbPassword = "wyty640"
	dbName     = "ct"

	sessionKey  = "IAsOAlsdkawpkodpwaoADas"
	sessionName = "ct-session"
)

// Server manages login, registration and messaging
type Server struct {
	messages  []*Message
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	logCh     chan string
	db        *sql.DB
	auth      *Auth
	store     *sessions.CookieStore
}

// JSON object
type JSON map[string]interface{}

// NewServer creates an instance of server
func NewServer() *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	logCh := make(chan string)

	store := sessions.NewCookieStore([]byte(sessionKey))
	gob.Register(&User{})

	// Db connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	if err != nil {
		panic(err)
	}

	auth := NewAuth(db)

	return &Server{
		messages:  messages,
		clients:   clients,
		addCh:     addCh,
		delCh:     delCh,
		sendAllCh: sendAllCh,
		doneCh:    doneCh,
		logCh:     logCh,
		db:        db,
		auth:      auth,
		store:     store,
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

func (s *Server) Messages() []*Message {
	return s.messages
}

// Listen starts the server
// Handle client connections and message broadcast
func (s *Server) Listen() {

	// Registered handlers
	http.HandleFunc(urlSend, s.authenticate(s.handleInitChat))
	http.HandleFunc(urlLogout, s.authenticate(s.handleLogout))
	http.HandleFunc(urlProfileModel, s.authenticate(s.handleGetProfileModel))

	// Public handlers
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

func (s *Server) handleInitChat(w http.ResponseWriter, r *http.Request, user *User) {
	onConnected := func(conn *websocket.Conn) {
		// Whenever this function exits close the connection
		defer func() {
			err := conn.Close()
			if err != nil {
				s.Err(err)
			}
		}()

		// Add then client to the server
		client := NewClient(conn, s, user)
		user.Client = client
		s.Add(client)
		// Listen untill the client disconnects
		client.Listen()
	}

	onConnectedHander := websocket.Handler(onConnected)
	onConnectedHander.ServeHTTP(w, r)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request, user *User) {
	s.endSession(w, r)
	if user.Client != nil {
		s.Del(user.Client)
	}
	s.sendJSON(w, &JSON{"Result": true})
}

func (s *Server) handleGetProfileModel(w http.ResponseWriter, r *http.Request, user *User) {
	s.sendJSON(w, &JSON{"Model": user})
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
		err = s.startSession(w, r, user)
		if err != nil {
			s.Err(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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

	user, err := s.auth.Register(regInfo)
	if err != nil {
		s.Err(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user != nil {
		err = s.startSession(w, r, user)
		if err != nil {
			s.Err(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	s.sendJSON(w, &JSON{"Result": user != nil, "User": user})
}

// Send a message to every client
func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		c.Send(msg)
	}
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

// authenticate makes sure the user is logged, if not it returns an AUTH_NEEDED_ERROR to the client
func (s *Server) authenticate(fn func(http.ResponseWriter, *http.Request, *User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.store.Get(r, sessionName)
		if err != nil {
			s.authNeededError(w)
			s.Err(err)
			return
		}
		if session.Values["User"] == nil {
			s.authNeededError(w)
			return
		}
		user := session.Values["User"].(*User)

		fn(w, r, user)
	}

}

func (s *Server) startSession(w http.ResponseWriter, r *http.Request, user *User) error {
	session, err := s.store.Get(r, sessionName)
	session.Values["User"] = user
	session.Save(r, w)

	return err
}

func (s *Server) endSession(w http.ResponseWriter, r *http.Request) error {
	session, err := s.store.Get(r, sessionName)
	user := session.Values["User"].(*User)
	if user != nil && user.Client != nil {
		user.Client.Done()
		s.Del(user.Client)
	}
	session.Values["User"] = nil
	session.Save(r, w)

	return err
}

// authNeededError returns an error to an unauthentificated client when the page requires authentification
func (s *Server) authNeededError(w http.ResponseWriter) {
	s.sendJSON(w, &JSON{"Response": "AUTH_NEEDED_ERROR"})
}
