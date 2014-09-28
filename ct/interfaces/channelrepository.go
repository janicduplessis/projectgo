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

func (repo *DbChannelRepo) GetChannels() (map[int64]*domain.Channel, error) {

	channels := make(map[int64]*domain.Channel)

	rows, err := repo.dbHandler.Query(`SELECT *
				 	    		 	   FROM channel`)

	if err != nil && err != ErrNoRows {
		return nil, err
	}

	defer rows.Close()

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
		channels[channel.Id] = channel
	}

	return channels, nil
}

func (repo *DbChannelRepo) Store(channel *domain.Channel) error {
	res, err := repo.dbHandler.Execute(`INSERT INTO channel (Name)
						   			    VALUES (?)`,
		channel.Name)

	if err != nil {
		return err
	}

	chanId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	channel.Id = chanId

	return nil
}

func (repo *DbChannelRepo) FindById(id int64) (*domain.Channel, error) {
	var (
		channelId int64
		name      string
	)
	err := repo.dbHandler.QueryRow(`SELECT *
				 	    		    FROM channel
				 	     		    WHERE ChannelId = ?`, id).Scan(&channelId, &name)

	if err != nil {
		return nil, err
	}

	channel := &domain.Channel{
		Id:   channelId,
		Name: name,
	}

	return channel, nil
}
