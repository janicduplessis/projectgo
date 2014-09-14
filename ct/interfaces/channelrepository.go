package interfaces

import (
	"github.com/janicduplessis/projectgo/ct/domain"
)

type DbChannelRepo DbRepo

func NewDbChannelRepo(dbHandlers map[string]DbHandler) *DbChannelRepo {
	dbHandler := dbHandlers["DbChannelRepo"]

	return &DbChannelRepo{
		dbHandlers: dbHandlers,
		dbHandler:  dbHandler,
	}
}

func (repo *DbChannelRepo) Channels() ([]*domain.Channel, error) {

	channels := make([]*domain.Channel, 0)

	// Check if the username is available
	rows, err := repo.dbHandler.Query(`SELECT *
				 	    		 	   FROM channel`)

	if err != nil && err != ErrNoRows {
		return channels, err
	}

	var (
		id   int64
		name string
	)

	for rows.Next() {
		if err := rows.Scan(&id, &name); err != nil {
			return channels, err
		}

		channel := domain.NewChannel(name)
		channel.Id = id
		channels = append(channels, channel)
	}

	return channels, nil
}

func (repo *DbChannelRepo) Store(channel *domain.Channel) error {
	return nil
}

func (repo *DbChannelRepo) FindById(id int64) (*domain.Channel, error) {
	return nil, nil
}
