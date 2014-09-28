package interfaces

import (
	"time"

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

func (repo *DbMessageRepo) FindByChannelId(channelId int64, max int) ([]*domain.Message, error) {
	messages := make([]*domain.Message, 0)

	rows, err := repo.dbHandler.Query(`SELECT m.MessageId, c.DisplayName, m.Body, m.Time, m.ClientId
				 	    		 	   FROM message m
				 	    		 	   JOIN client c ON m.ClientId = c.ClientId
				 	    		 	   WHERE ChannelId=?
				 	    		 	   LIMIT ?`, channelId, max)

	if err != nil && err != ErrNoRows {
		return nil, err
	}

	defer rows.Close()

	var (
		id          int64
		author      string
		body        string
		dateTimeStr string
		clientId    int64
	)

	for rows.Next() {
		if err := rows.Scan(&id, &author, &body, &dateTimeStr, &clientId); err != nil {
			return nil, err
		}

		dateTime, err := time.Parse(dbDateFormat, dateTimeStr)
		if err != nil {
			return nil, err
		}

		message := &domain.Message{
			Id:        id,
			Author:    author,
			Body:      body,
			Time:      dateTime,
			ClientId:  clientId,
			ChannelId: channelId,
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (repo *DbMessageRepo) FindByClientId(clientId int64, max int) ([]*domain.Message, error) {
	messages := make([]*domain.Message, 0)

	rows, err := repo.dbHandler.Query(`SELECT m.MessageId, c.DisplayName, m.Body, m.Time, m.ChannelId
				 	    		 	   FROM message m
				 	    		 	   JOIN client c ON m.ClientId = c.ClientId
				 	    		 	   WHERE ClientId=?
				 	    		 	   LIMIT ?`, clientId, max)

	if err != nil && err != ErrNoRows {
		return nil, err
	}

	defer rows.Close()

	var (
		id          int64
		author      string
		body        string
		dateTimeStr string
		channelId   int64
	)

	for rows.Next() {
		if err := rows.Scan(&id, &body, &dateTimeStr, &clientId); err != nil {
			return nil, err
		}

		dateTime, err := time.Parse(dbDateFormat, dateTimeStr)
		if err != nil {
			return nil, err
		}

		message := &domain.Message{
			Id:        id,
			Author:    author,
			Body:      body,
			Time:      dateTime,
			ClientId:  clientId,
			ChannelId: channelId,
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (repo *DbMessageRepo) Store(msg *domain.Message) error {
	res, err := repo.dbHandler.Execute(`INSERT INTO message (Body, Time, ClientId, ChannelId)
						   			    VALUES (?,?,?,?)`,
		msg.Body, msg.Time.Format(dbDateFormat), msg.ClientId, msg.ChannelId)

	if err != nil {
		return err
	}

	msgId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	msg.Id = msgId

	return nil
}
