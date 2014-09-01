package ct

import (
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"code.google.com/p/go.net/websocket"
	"github.com/gorilla/sessions"
)

const (
	urlLogin      = "login"
	urlRegister   = "register"
	urlSend       = "send"
	urlLogout     = "logout"
	urlJoinChan   = "chan/join"
	urlLeaveChan  = "chan/leave"
	urlCreateChan = "chan/create"
	urlListChan   = "chan/list"

	urlProfileModel = "/models/getProfileModel"

	sessionKey  = "IAsOAlsdkawpkodpwaoADas"
	sessionName = "ct-session"
)

// Server manages login, registration and messaging
type Server struct {
	Channels map[int64]*Channel
	Db       *sql.DB

	clients         map[int]*Client
	addCh           chan *Client
	delCh           chan *Client
	sendCh          chan *Message
	joinChannelCh   chan ChannelUser
	leaveChannelCh  chan ChannelUser
	createChannelCh chan *ChannelInfo
	channelsMutex   sync.RWMutex
	doneCh          chan bool
	logCh           chan string
	auth            *Auth
	store           *sessions.CookieStore
	config          *ServerConfig
}

type ServerConfig struct {
	SiteRoot   string
	SitePort   string
	DbUser     string
	DbPassword string
	DbName     string
	DbUrl      string
	DbPort     string
}

// JSON object
type JSON map[string]interface{}

// NewServer creates an instance of server
func NewServer(config *ServerConfig) *Server {
	channels := make(map[int64]*Channel)
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendCh := make(chan *Message)
	doneCh := make(chan bool)
	logCh := make(chan string)
	joinChannelCh := make(chan ChannelUser)
	leaveChannelCh := make(chan ChannelUser)
	createChannelCh := make(chan *ChannelInfo)

	store := sessions.NewCookieStore([]byte(sessionKey))
	gob.Register(&User{})

	// Db connection
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DbUser, config.DbPassword, config.DbUrl, config.DbPort, config.DbName))
	if err != nil {
		panic(err)
	}

	auth := NewAuth(db)

	return &Server{
		Channels: channels,
		Db:       db,

		clients:         clients,
		addCh:           addCh,
		delCh:           delCh,
		sendCh:          sendCh,
		joinChannelCh:   joinChannelCh,
		leaveChannelCh:  leaveChannelCh,
		createChannelCh: createChannelCh,
		doneCh:          doneCh,
		logCh:           logCh,
		auth:            auth,
		store:           store,
		config:          config,
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

func (s *Server) Send(msg *Message) {
	s.sendCh <- msg
}

func (s *Server) GetChannel(channelId int64) *Channel {
	//We'll user a mutex for reading channels values, writes will be handled in main loop
	s.channelsMutex.RLock()
	channel := s.Channels[channelId]
	s.channelsMutex.RUnlock()
	if channel == nil {
		s.Log("e", fmt.Sprintf("Cannot find channel %s", channelId))
	}
	return channel
}

func (s *Server) GetAllChannels() []*Channel {
	s.channelsMutex.RLock()
	channels := make([]*Channel, 0, len(s.Channels))
	for _, channel := range s.Channels {
		channels = append(channels, channel)
	}
	s.channelsMutex.RUnlock()
	return channels
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

	// Registered handlers
	http.HandleFunc(s.makeServerUrl(urlSend), s.authenticate(s.handleInitChat))
	http.HandleFunc(s.makeServerUrl(urlLogout), s.authenticate(s.handleLogout))
	http.HandleFunc(s.makeServerUrl(urlProfileModel), s.authenticate(s.handleGetProfileModel))
	http.HandleFunc(s.makeServerUrl(urlJoinChan), s.authenticate(s.handleJoinChan))
	http.HandleFunc(s.makeServerUrl(urlLeaveChan), s.authenticate(s.handleLeaveChan))
	http.HandleFunc(s.makeServerUrl(urlCreateChan), s.authenticate(s.handleCreateChan))
	http.HandleFunc(s.makeServerUrl(urlListChan), s.authenticate(s.handleListChan))

	// Public handlers
	http.HandleFunc(s.makeServerUrl(urlLogin), s.handleLogin)
	http.HandleFunc(s.makeServerUrl(urlRegister), s.handleRegister)

	defer s.Db.Close()

	// Server listen loop
	for {
		select {

		// Add a client
		case c := <-s.addCh:
			s.clients[c.id] = c

		// Disconnect a client
		case c := <-s.delCh:
			delete(s.clients, c.id)

		case msg := <-s.sendCh:
			s.send(msg)

		// Log errors
		case message := <-s.logCh:
			log.Println(message)

		case chanUser := <-s.joinChannelCh:
			s.joinChannel(chanUser)

		case chanUser := <-s.leaveChannelCh:
			s.leaveChannel(chanUser)

		case chanInfo := <-s.createChannelCh:
			s.createChannel(chanInfo)

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

func (s *Server) handleJoinChan(w http.ResponseWriter, r *http.Request, user *User) {
	channel := ChannelUser{}
	s.readJSON(r, &channel)
	channel.User = user
	s.joinChannelCh <- channel
	s.sendJSON(w, &JSON{"Result": true})
}

func (s *Server) handleLeaveChan(w http.ResponseWriter, r *http.Request, user *User) {
	channel := ChannelUser{}
	s.readJSON(r, &channel)
	channel.User = user
	s.leaveChannelCh <- channel
	s.sendJSON(w, &JSON{"Result": true})
}

func (s *Server) handleCreateChan(w http.ResponseWriter, r *http.Request, user *User) {
	channel := &ChannelInfo{}
	s.readJSON(r, channel)
	s.Log("d", channel.Name)
	s.createChannelCh <- channel
	s.sendJSON(w, &JSON{"Result": true})
}

func (s *Server) handleListChan(w http.ResponseWriter, r *http.Request, user *User) {
	channels := s.GetAllChannels()
	s.sendJSON(w, channels)
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

func (s *Server) send(msg *Message) {
	channel := s.Channels[msg.ChannelId]
	if channel == nil {
		return
	}
	msg.Channel = channel
	channel.AddMessage(msg)
	for _, user := range channel.Users {
		user.Client.Send(msg)
	}
}

func (s *Server) joinChannel(cu ChannelUser) {
	channel := s.GetChannel(cu.ChannelId)
	channel.Join(cu.User)
}

func (s *Server) leaveChannel(cu ChannelUser) {
	channel := s.GetChannel(cu.ChannelId)
	channel.Leave(cu.User)
}

func (s *Server) createChannel(info *ChannelInfo) {
	channel, err := NewChannel(info.Name, s)
	if err != nil {
		s.Err(err)
		return
	}
	if channel == nil {
		s.Log("e", "Could not create channel")
		return
	}
	s.Channels[channel.Id] = channel
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

func (s *Server) makeServerUrl(url string) string {
	return s.config.SiteRoot + url
}
