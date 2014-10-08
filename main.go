package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"

	//"github.com/janicduplessis/projectgo/ct/domain"
	"github.com/janicduplessis/projectgo/ct/infrastructure"
	"github.com/janicduplessis/projectgo/ct/interfaces"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

const configFile = "server.json"

type ServerConfig struct {
	SiteRoot   string
	SitePort   string
	DbUser     string
	DbPassword string
	DbName     string
	DbUrl      string
	DbPort     string
}

func main() {
	log.SetFlags(log.Lshortfile)

	// Default config
	config := ServerConfig{
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

	//TODO move this elsewhere
	_, err = os.Stat("upload")
	if os.IsNotExist(err) {
		os.Mkdir("upload", 0777)
	}

	//Console logger
	logger := new(infrastructure.LoggerHandler)

	// Crypto handler
	crypto := new(infrastructure.CryptoHandler)

	// Base webservice handler
	webservice := infrastructure.NewWebserviceHandler(logger)

	// Base websocket handler
	websocket := infrastructure.NewWebsocketHandler(logger)

	imageUtils := new(infrastructure.ImageUtilsHandler)

	// Database
	dbConfig := infrastructure.MySqlDbConfig{
		User:     config.DbUser,
		Password: config.DbPassword,
		Name:     config.DbName,
		Url:      config.DbUrl,
		Port:     config.DbPort,
	}
	dbHandler := infrastructure.NewMySqlHandler(dbConfig)

	// Database handlers for each repo
	handlers := make(map[string]interfaces.DbHandler)
	handlers["DbInitializerRepo"] = dbHandler
	handlers["DbClientRepo"] = dbHandler
	handlers["DbUserRepo"] = dbHandler
	handlers["DbMessageRepo"] = dbHandler
	handlers["DbChannelRepo"] = dbHandler

	// Initialize the database
	dbInit := interfaces.NewDbInitializerRepo(handlers)
	dbInit.Init()

	//Repos
	clientRepo := interfaces.NewDbClientRepo(handlers)

	// Interactors
	authInteractor := usecases.NewAuthentificationInteractor(interfaces.NewDbUserRepo(handlers), crypto, logger)

	chatInteractor := new(usecases.ChatInteractor)
	chatInteractor.ServerRepository = interfaces.NewSingletonServerRepo(handlers)
	chatInteractor.ChannelRepository = interfaces.NewDbChannelRepo(handlers)
	chatInteractor.MessageRepository = interfaces.NewDbMessageRepo(handlers)
	chatInteractor.ClientRepository = clientRepo
	chatInteractor.Logger = logger

	homeInteractor := usecases.NewHomeInteractor(clientRepo, logger)

	// Webservices
	interfaces.NewAuthentificationWebservice(webservice, authInteractor, chatInteractor)
	interfaces.NewChatWebservice(webservice, websocket, chatInteractor)
	interfaces.NewHomeWebservice(webservice, homeInteractor, imageUtils)

	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.SitePort), context.ClearHandler(http.DefaultServeMux)))
}
