package config

import (
	"encoding/json"
	"flag"
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

	UseS3       bool
	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
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

	S3AccessKey string
	S3SecretKey string
	S3Bucket    string
}

func init() {
	envConfig := flag.Bool("useenv", false, "Use environnement variables config")
	useS3 := flag.Bool("uses3", false, "Use amazon s3 for file storage")
	flag.Parse()

	UseS3 = *useS3

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

	if *envConfig {
		val := os.Getenv("SITE_URL")
		if len(val) > 0 {
			config.SiteUrl = val
		}
		val = os.Getenv("PORT")
		if len(val) > 0 {
			config.SitePort = val
		}
		val = os.Getenv("DB_USER")
		if len(val) > 0 {
			config.DbUser = val
		}
		val = os.Getenv("DB_PASSWORD")
		if len(val) > 0 {
			config.DbPassword = val
		}
		val = os.Getenv("DB_NAME")
		if len(val) > 0 {
			config.DbName = val
		}
		val = os.Getenv("DB_URL")
		if len(val) > 0 {
			config.DbUrl = val
		}
		val = os.Getenv("DB_PORT")
		if len(val) > 0 {
			config.DbPort = val
		}
		val = os.Getenv("OAUTH2_CLIENT_ID")
		if len(val) > 0 {
			config.OAuth2ClientId = val
		}
		val = os.Getenv("OAUTH2_CLIENT_SECRET")
		if len(val) > 0 {
			config.OAuth2ClientSecret = val
		}
		val = os.Getenv("S3_ACCESS_KEY")
		if len(val) > 0 {
			config.S3AccessKey = val
		}
		val = os.Getenv("S3_SECRET_KEY")
		if len(val) > 0 {
			config.S3SecretKey = val
		}
		val = os.Getenv("S3_BUCKET")
		if len(val) > 0 {
			config.S3Bucket = val
		}
	} else {
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
		}

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

	S3AccessKey = config.S3AccessKey
	S3SecretKey = config.S3SecretKey
	S3Bucket = config.S3Bucket

	log.Println("---------------------")
	log.Println("-     Config        -")
	log.Println("---------------------")
	log.Println(fmt.Sprintf("%s: %s", "SiteUrl", SiteUrl))
	log.Println(fmt.Sprintf("%s: %s", "SitePort", SitePort))
	log.Println(fmt.Sprintf("%s: %s", "DbUser", DbUser))
	log.Println(fmt.Sprintf("%s: %s", "DbPassword", DbPassword))
	log.Println(fmt.Sprintf("%s: %s", "DbName", DbName))
	log.Println(fmt.Sprintf("%s: %s", "DbUrl", DbUrl))
	log.Println(fmt.Sprintf("%s: %s", "DbPort", DbPort))
	log.Println("-------- OAuth ------")
	log.Println(fmt.Sprintf("%s: %s", "OAuth2ClientId", OAuth2ClientId))
	log.Println(fmt.Sprintf("%s: %s", "OAuth2ClientSecret", OAuth2ClientSecret))
	if UseS3 {
		log.Println("--------- S3 --------")
		log.Println(fmt.Sprintf("%s: %s", "S3AccessKey", S3AccessKey))
		log.Println(fmt.Sprintf("%s: %s", "S3SecretKey", S3SecretKey))
		log.Println(fmt.Sprintf("%s: %s", "S3Bucket", S3Bucket))
	}
	log.Println("---------------------")

	_, err := os.Stat("upload")
	if os.IsNotExist(err) {
		os.Mkdir("upload", 0777)
	}
}
