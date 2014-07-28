package main

import (
	"code.google.com/p/go.net/websocket"
)

type Client struct {
	id     int
	name   string
	conn   *websocket.Conn
	server *Server
	msgCh  chan *Message
	doneCh chan bool
}

func NewClient() *Client {

}
