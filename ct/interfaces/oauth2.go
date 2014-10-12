package interfaces

import (
	"net/http"

	"code.google.com/p/go.net/context"

	"github.com/janicduplessis/projectgo/ct/usecases"
)

type OAuth2 interface {
	GetProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) (*usecases.GoogleRegisterInfo, error)
	GetScope() string
}

type OAuth2Profile struct {
	Id           string
	DisplayName  string
	FirstName    string
	LastName     string
	ProfileImage string
	Email        string
}
