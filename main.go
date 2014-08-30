package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/janicduplessis/projectgo/ct"
)

const configFile = "server.json"

func main() {
	log.SetFlags(log.Lshortfile)

	// Default config
	config := ct.ServerConfig{
		DbUser:     "ct",
		DbPassword: "***",
		DbName:     "ct",
		DbUrl:      "localhost",
		DbPort:     "3306",
	}

	//Get server config
	file, err := ioutil.ReadFile(configFile)

	if err != nil {
		data, err := json.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(configFile, data, 0600)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		if err = json.Unmarshal(file, &config); err != nil {
			log.Fatal(err)
		}
	}

	// Chat server
	server := ct.NewServer(&config)
	go server.Listen()

	http.Handle("/", http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux)))
}
