package ct

import ()

type User struct {
	UserId    int
	UserName  string
	FirstName string
	LastName  string
	Email     string

	Client *Client
}
