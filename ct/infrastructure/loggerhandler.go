package infrastructure

import (
	"log"
)

type LoggerHandler struct{}

func (handler *LoggerHandler) Log(message string) {
	log.Println(message)
}

func (handler *LoggerHandler) Error(err error) {
	log.Println(err.Error())
}
