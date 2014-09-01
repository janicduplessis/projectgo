package main

import (
	"encoding/json"
	"fmt"
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
		SiteRoot:   "/",
		SitePort:   "8080",
		DbUser:     "ct",
		DbPassword: "***",
		DbName:     "ct",
		DbUrl:      "localhost",
		DbPort:     "3306",
	}

	//Get server config
	file, err := ioutil.ReadFile(configFile)

	if err != nil {
		// No config found, we will create the default one and tell the user to set it up
		data, err := json.Marshal(config)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(configFile, data, 0600)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("No config found, created default config file. Please edit 'server.json' and try again.")
		// Exit the program
		return
	} else {
		if err = json.Unmarshal(file, &config); err != nil {
			log.Fatal(err)
		}
	}

	// Chat server
	server := ct.NewServer(&config)
	go server.Listen()

	http.Handle("/", http.FileServer(http.Dir("public")))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.SitePort), context.ClearHandler(http.DefaultServeMux)))
}
