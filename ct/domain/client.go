package domain

type ClientRepository interface {
	Store(client *Client)
	FindById(id int64) *Client
}

type Client struct {
	Id        int64
	Username  string
	FirstName string
	LastName  string
	Email     string
}
