package interfaces

import (
	"code.google.com/p/go.net/context"
	"net/http"

	"github.com/janicduplessis/projectgo/ct/usecases"
)

const (
	KeyUser string = "User"
)

type Webservice interface {
	SendJson(w http.ResponseWriter, obj interface{})
	ReadJson(w http.ResponseWriter, r *http.Request, obj interface{}) error
	Error(w http.ResponseWriter, err error)
	StartSession(ctx context.Context, w http.ResponseWriter, r *http.Request, user *usecases.User) error
	EndSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error
	AddHandler(url string, authenticated bool, fn func(context.Context, http.ResponseWriter, *http.Request))
}
