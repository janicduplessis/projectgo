package interfaces

import (
	"github.com/janicduplessis/projectgo/ct/domain"
)

type DbClientRepo DbRepo

func NewDbClientRepo(dbHandlers map[string]DbHandler) *DbClientRepo {
	dbHandler := dbHandlers["DbClientRepo"]

	return &DbClientRepo{
		dbHandlers: dbHandlers,
		dbHandler:  dbHandler,
	}
}

func (repo *DbClientRepo) FindById(id int64) (*domain.Client, error) {
	var (
		clientId    int64
		displayName string
		firstName   string
		lastName    string
		email       string
	)
	err := repo.dbHandler.QueryRow(`SELECT *
				 	    		    FROM client
				 	     		    WHERE ClientId = ?`, id).Scan(&clientId, &displayName, &firstName, &lastName, &email)

	if err != nil {
		if err == ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	client := &domain.Client{
		Id:          clientId,
		DisplayName: displayName,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
	}

	return client, nil
}

func (repo *DbClientRepo) Store(msg *domain.Client) error {
	//TODO: NYI
	return nil
}
