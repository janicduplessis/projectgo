package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/context"

	"github.com/janicduplessis/projectgo/ct/config"
	"github.com/janicduplessis/projectgo/ct/infrastructure"
	"github.com/janicduplessis/projectgo/ct/interfaces"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

func main() {
	log.SetFlags(log.Lshortfile)

	//Console logger
	logger := new(infrastructure.LoggerHandler)

	// Crypto handler
	crypto := new(infrastructure.CryptoHandler)

	// Base webservice handler
	webservice := infrastructure.NewWebserviceHandler(logger)

	// Base websocket handler
	websocket := infrastructure.NewWebsocketHandler(logger)

	imageUtils := new(infrastructure.ImageUtilsHandler)

	var fileStore interfaces.FileStore
	if config.UseS3 {
		fileStoreHandler := new(infrastructure.S3FileStorageHandler)
		fileStoreHandler.Init()
		fileStore = fileStoreHandler
	} else {
		fileStore = new(infrastructure.LocalFileStoreHandler)
	}

	oauth2 := new(infrastructure.OAuth2Handler)
	oauth2.Init()
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

	chatInteractor := usecases.NewChatInteractor(interfaces.NewSingletonServerRepo(handlers), interfaces.NewDbChannelRepo(handlers),
		interfaces.NewDbMessageRepo(handlers), clientRepo, logger)

	homeInteractor := usecases.NewHomeInteractor(clientRepo, logger)

	// Webservices
	interfaces.NewAuthentificationWebservice(webservice, oauth2, authInteractor, chatInteractor)
	interfaces.NewChatWebservice(webservice, websocket, chatInteractor)
	interfaces.NewHomeWebservice(webservice, homeInteractor, imageUtils, fileStore)

	http.Handle("/", http.FileServer(http.Dir("web")))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.SitePort), context.ClearHandler(http.DefaultServeMux)))
}
