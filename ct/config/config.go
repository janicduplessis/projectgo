package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const configFile = "server.json"

var (
	SiteUrl    string
	SitePort   string
	DbUser     string
	DbPassword string
	DbName     string
	DbUrl      string
	DbPort     string

	OAuth2ClientId     string
	OAuth2ClientSecret string
)

type serverConfig struct {
	SiteUrl    string
	SitePort   string
	DbUser     string
	DbPassword string
	DbName     string
	DbUrl      string
	DbPort     string

	OAuth2ClientId     string
	OAuth2ClientSecret string
}

func init() {
	// Default config
	config := serverConfig{
		SiteUrl:    "localhost:8080",
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

	SiteUrl = config.SiteUrl
	SitePort = config.SitePort
	DbUser = config.DbUser
	DbPassword = config.DbPassword
	DbName = config.DbName
	DbUrl = config.DbUrl
	DbPort = config.DbPort

	OAuth2ClientId = config.OAuth2ClientId
	OAuth2ClientSecret = config.OAuth2ClientSecret

	log.Println("---------------------")
	log.Println("- Config            -")
	log.Println("---------------------")
	log.Println(fmt.Sprintf("%s: %s", "SiteUrl", SiteUrl))
	log.Println(fmt.Sprintf("%s: %s", "SitePort", SitePort))
	log.Println(fmt.Sprintf("%s: %s", "DbUser", DbUser))
	log.Println(fmt.Sprintf("%s: %s", "DbPassword", DbPassword))
	log.Println(fmt.Sprintf("%s: %s", "DbName", DbName))
	log.Println(fmt.Sprintf("%s: %s", "DbUrl", DbUrl))
	log.Println(fmt.Sprintf("%s: %s", "DbPort", DbPort))
	log.Println(fmt.Sprintf("%s: %s", "OAuth2ClientId", OAuth2ClientId))
	log.Println(fmt.Sprintf("%s: %s", "OAuth2ClientSecret", OAuth2ClientSecret))
	log.Println("---------------------")

	_, err = os.Stat("upload")
	if os.IsNotExist(err) {
		os.Mkdir("upload", 0777)
	}
}
