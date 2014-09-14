package interfaces

import (
	"github.com/janicduplessis/projectgo/ct/domain"
)

type DbMessageRepo DbRepo

func NewDbMessageRepo(dbHandlers map[string]DbHandler) *DbMessageRepo {
	dbHandler := dbHandlers["DbMessageRepo"]

	return &DbMessageRepo{
		dbHandlers: dbHandlers,
		dbHandler:  dbHandler,
	}
}

func (repo *DbMessageRepo) FindById(id int64) (*domain.Message, error) {
	//TODO: NYI
	return nil, nil
}

func (repo *DbMessageRepo) Store(msg *domain.Message) error {
	//TODO: NYI
	return nil
}
