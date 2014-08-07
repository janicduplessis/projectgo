package main

import (
	"log"
	"net/http"

	"github.com/janicduplessis/projectgo/ct"
)

func main() {
	log.SetFlags(log.Lshortfile)

	// Chat server
	server := ct.NewServer()
	go server.Listen()

	http.Handle("/", http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
