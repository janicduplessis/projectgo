package ct

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Channel struct {
	Id       int64
	Name     string
	Users    map[int64]*User
	Messages []*Message
	Server   *Server
}

type ChannelInfo struct {
	Name string
}

type ChannelUser struct {
	User      *User
	ChannelId int64
}

func NewChannel(name string, server *Server) (*Channel, error) {

	// Check if the channel name is available
	err := server.Db.QueryRow(`SELECT 1
				 	   	  	   FROM channel
				 	   	 	   WHERE Name = ?`, name).Scan(new(int))
	if err != sql.ErrNoRows {
		if err != nil {
			return nil, err
		}

		// Channel already exists
		return nil, nil
	}

	// Create the channel
	res, err := server.Db.Exec(`INSERT INTO channel (Name)
						   		VALUES (?)`, name)

	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	users := make(map[int64]*User)
	messages := []*Message{}

	return &Channel{
		Id:       id,
		Name:     name,
		Users:    users,
		Messages: messages,
		Server:   server,
	}, nil

}

func (c *Channel) Join(user *User) {
	c.Users[user.UserId] = user
}

func (c *Channel) Leave(user *User) {
	delete(c.Users, user.UserId)
}

func (c *Channel) AddMessage(message *Message) {
	c.Messages = append(c.Messages, message)
}
