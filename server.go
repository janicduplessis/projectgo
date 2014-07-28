package main

import (
	"net/http"
)

type Server struct {
	//Rooms   []Room
	Files   map[string]*File
	Clients map[string]*Client
}

func (s *Server) makeHander(fn func(http.ResponseWriter, *http.Request, *Client)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("id")
		client := s.Clients[id]
		fn(w, r, client)
	}
}

func (s *Server) sendMessage(http.ResponseWriter, *http.Request, *Client) {

}

func (s *Server) Start() {
	s.Files = make(map[string]*File)
	s.Clients = make(map[string]*Client)
	http.HandleFunc("/send", s.makeHander(s.sendMessage))
}

func (s *Server) AddFile(name string, data []byte) {
	file := &File{Name: name, Size: len(data)}
	s.Files[name] = file

}
