package interfaces

import (
	"net/http"

	"code.google.com/p/go.net/context"
)

type OAuth2 interface {
	GetProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) (*OAuth2Profile, error)
	GetUrl() (string, error)
}

type OAuth2Profile struct {
	Id           string
	DisplayName  string
	FirstName    string
	LastName     string
	ProfileImage string
}
