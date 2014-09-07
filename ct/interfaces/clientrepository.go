package interfaces

import (
//"github.com/janicduplessis/projectgo/ct/domain"
//"github.com/janicduplessis/projectgo/ct/usecases"
)

type DbClientRepo DbRepo

func NewDbClientRepo(dbHandlers map[string]DbHandler) *DbClientRepo {
	dbHandler := dbHandlers["DbClientRepo"]

	return &DbClientRepo{
		dbHandlers: dbHandlers,
		dbHandler:  dbHandler,
	}
}
