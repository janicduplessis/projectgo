package interfaces

import (
	"net/http"

	"code.google.com/p/go.net/context"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

const (
	urlLogin    string = "/login"
	urlRegister string = "/register"
)

type AuthentificationInteractor interface {
	Login(info *usecases.LoginInfo) (*usecases.User, error)
	Register(info *usecases.RegisterInfo) (*usecases.User, error)
}

type AuthentificationWebserviceHandler struct {
	Webservice                 Webservice
	AuthentificationInteractor AuthentificationInteractor
}

type LoginResponseModel struct {
	Result bool
	User   *UserModel
}

type RegisterResponseModel struct {
	Result bool
	User   *UserModel
}

type UserModel struct {
	Id        int64
	Username  string
	FirstName string
	LastName  string
	Email     string
}

func NewAuthentificationWebservice(ws Webservice, ai AuthentificationInteractor) *AuthentificationWebserviceHandler {
	wsHandler := &AuthentificationWebserviceHandler{
		Webservice:                 ws,
		AuthentificationInteractor: ai,
	}

	ws.AddHandler(urlLogin, false, wsHandler.Login)
	ws.AddHandler(urlRegister, false, wsHandler.Register)

	return wsHandler
}

func (handler *AuthentificationWebserviceHandler) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	info := new(usecases.LoginInfo)
	if err := handler.Webservice.ReadJson(w, r, info); err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	user, err := handler.AuthentificationInteractor.Login(info)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	var userModel *UserModel

	if user != nil {
		userModel = &UserModel{
			Id:        user.Id,
			Username:  user.Client.Username,
			FirstName: user.Client.FirstName,
			LastName:  user.Client.LastName,
			Email:     user.Client.Email,
		}

		handler.Webservice.StartSession(ctx, w, r, user)
	}

	response := &LoginResponseModel{
		Result: user != nil,
		User:   userModel,
	}

	handler.Webservice.SendJson(w, response)
}

func (handler *AuthentificationWebserviceHandler) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	info := new(usecases.RegisterInfo)
	if err := handler.Webservice.ReadJson(w, r, info); err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	user, err := handler.AuthentificationInteractor.Register(info)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	var userModel *UserModel

	if user != nil {
		userModel = &UserModel{
			Id:        user.Id,
			Username:  user.Client.Username,
			FirstName: user.Client.FirstName,
			LastName:  user.Client.LastName,
			Email:     user.Client.Email,
		}

		handler.Webservice.StartSession(ctx, w, r, user)
	}

	response := &RegisterResponseModel{
		Result: user != nil,
		User:   userModel,
	}

	handler.Webservice.SendJson(w, response)
}
