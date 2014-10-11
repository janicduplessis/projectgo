package interfaces

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"code.google.com/p/go.net/context"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

const (
	urlLogin       string = "/login"
	urlRegister    string = "/register"
	UrlOAuth2Login string = "/oauth2login"
	urlLogout      string = "/logout"

	urlGetLoginModel string = "/models/getLoginModel"
)

type AuthentificationInteractor interface {
	Login(info *usecases.LoginInfo) (*usecases.User, error)
	Register(info *usecases.RegisterInfo) (*usecases.User, error)
	GoogleLogin(info *usecases.GoogleRegisterInfo) (*usecases.User, error)
}

type AuthentificationWebserviceHandler struct {
	Webservice                 Webservice
	OAuth2                     OAuth2
	AuthentificationInteractor AuthentificationInteractor
	ChatInteractor             ChatInteractor
}

func NewAuthentificationWebservice(ws Webservice, oauth2 OAuth2, ai AuthentificationInteractor, ci ChatInteractor) *AuthentificationWebserviceHandler {
	wsHandler := &AuthentificationWebserviceHandler{
		Webservice: ws,
		OAuth2:     oauth2,
		AuthentificationInteractor: ai,
		ChatInteractor:             ci,
	}

	ws.AddHandler(urlGetLoginModel, false, wsHandler.GetLoginModel)
	ws.AddHandler(urlLogin, false, wsHandler.Login)
	ws.AddHandler(UrlOAuth2Login, false, wsHandler.OAuth2Login)
	ws.AddHandler(urlRegister, false, wsHandler.Register)
	ws.AddHandler(urlLogout, true, wsHandler.Logout)

	return wsHandler
}

func (handler *AuthentificationWebserviceHandler) GetLoginModel(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	url, err := handler.OAuth2.GetUrl()
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	model := &LoginModel{
		GoogleLoginURL: url,
	}
	response := &ModelResponse{
		Model: model,
	}

	handler.Webservice.SendJson(w, response)
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
			Id:          user.Id,
			Username:    user.Username,
			DisplayName: user.Client.DisplayName,
			FirstName:   user.Client.FirstName,
			LastName:    user.Client.LastName,
			Email:       user.Client.Email,
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
	var errorMessage string
	info := new(usecases.RegisterInfo)
	if err := handler.Webservice.ReadJson(w, r, info); err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	user, err := handler.AuthentificationInteractor.Register(info)
	if err != nil {
		if err == usecases.ErrUserAlreadyExists {
			errorMessage = err.Error()
		} else {
			handler.Webservice.Error(w, err)
			return
		}
	}

	var userModel *UserModel

	if user != nil {
		userModel = &UserModel{
			Id:          user.Id,
			Username:    user.Username,
			DisplayName: user.Client.DisplayName,
			FirstName:   user.Client.FirstName,
			LastName:    user.Client.LastName,
			Email:       user.Client.Email,
		}

		handler.Webservice.StartSession(ctx, w, r, user)
	}

	response := &RegisterResponseModel{
		Result: user != nil,
		User:   userModel,
		Error:  errorMessage,
	}

	handler.Webservice.SendJson(w, response)
}

func (handler *AuthentificationWebserviceHandler) OAuth2Login(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	info, err := handler.OAuth2.GetProfile(ctx, w, r)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	user, err := handler.AuthentificationInteractor.GoogleLogin(info)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	handler.Webservice.StartSession(ctx, w, r, user)

	userModel := &UserModel{
		Id:          user.Id,
		DisplayName: user.Client.DisplayName,
		FirstName:   user.Client.FirstName,
		LastName:    user.Client.LastName,
		Email:       user.Client.Email,
	}

	userModelJSON, err := json.Marshal(userModel)

	// Adds a cookie with escaped JSON formatted user info for the client
	cookie := new(http.Cookie)
	cookie.Name = "ctuserinfo"
	cookie.Value = url.QueryEscape(string(userModelJSON))
	cookie.Expires = time.Now().AddDate(0, 1, 0)
	http.SetCookie(w, cookie)

	handler.Webservice.Redirect(w, r, "/")
}

func (handler *AuthentificationWebserviceHandler) Logout(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value(KeyUser).(*usecases.User)

	if err := handler.ChatInteractor.Disconnect(user.Id); err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	handler.Webservice.EndSession(ctx, w, r)

	reponse := &LogoutResponseModel{
		Result: true,
	}

	handler.Webservice.SendJson(w, reponse)
}
