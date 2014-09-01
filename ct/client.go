package ct

import (
	"io"

	"code.google.com/p/go.net/websocket"
)

const chanBufferSize = 128

var curId int = 0

type Client struct {
	id     int
	conn   *websocket.Conn
	server *Server
	msgCh  chan []*Message
	doneCh chan bool

	User *User
}

func NewClient(conn *websocket.Conn, s *Server, user *User) *Client {
	curId++
	msgCh := make(chan []*Message, chanBufferSize)
	doneCh := make(chan bool)

	return &Client{
		id:     curId,
		conn:   conn,
		server: s,
		msgCh:  msgCh,
		doneCh: doneCh,
		User:   user,
	}
}

func (c *Client) Conn() *websocket.Conn {
	return c.conn
}

func (c *Client) Send(msg *Message) {
	select {
	case c.msgCh <- []*Message{msg}:
	default:

	}
}

func (c *Client) SendArray(msg []*Message) {
	select {
	case c.msgCh <- msg:
	default:

	}
}

func (c *Client) Listen() {
	go c.listenSend()
	//c.SendArray(c.server.Messages())
	c.listenReceive()
}

func (c *Client) Done() {
	c.doneCh <- true
}

func (c *Client) listenSend() {
	for {
		select {
		case msg := <-c.msgCh:
			websocket.JSON.Send(c.conn, msg)
		}
	}
}

func (c *Client) listenReceive() {
	for {
		select {
		case <-c.doneCh:
			return
		default:
			var msg Message
			err := websocket.JSON.Receive(c.conn, &msg)
			msg.Author = c.User.String()
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				c.server.Send(&msg)
			}
		}
	}
}
