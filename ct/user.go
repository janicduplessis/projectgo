package ct

import ()

type User struct {
	UserId    int64
	UserName  string
	FirstName string
	LastName  string
	Email     string

	Client *Client
}

func (u *User) String() string {
	return u.FirstName + " " + u.LastName
}
