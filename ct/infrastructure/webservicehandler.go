package infrastructure

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"code.google.com/p/go.net/context"
	"github.com/gorilla/sessions"

	"github.com/janicduplessis/projectgo/ct/interfaces"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

const (
	sessionKey  = "IAsOAlsdkawpkodpwaoADas"
	sessionName = "ct-session"
)

type WebserviceHandler struct {
	Logger usecases.Logger

	store *sessions.CookieStore
}

type AuthNeededError struct {
	Response string
}

func NewWebserviceHandler(logger usecases.Logger) *WebserviceHandler {
	store := sessions.NewCookieStore([]byte(sessionKey))
	gob.Register(&usecases.User{})

	return &WebserviceHandler{
		Logger: logger,
		store:  store,
	}
}

func (handler *WebserviceHandler) ReadJson(w http.ResponseWriter, r *http.Request, obj interface{}) error {
	// Read the request's body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, obj)
	if err != nil {
		return err
	}
	return nil
}

func (handler *WebserviceHandler) SendJson(w http.ResponseWriter, obj interface{}) {
	//Convert object to json
	bytes, err := json.Marshal(obj)
	if err != nil {
		handler.Error(w, err)
		return
	}
	//Set content type
	w.Header().Set("Content-Type", "application/json")
	//Write response body
	w.Write(bytes)
}

func (handler *WebserviceHandler) Redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (handler *WebserviceHandler) Log(msg string) {
	handler.Logger.Log(msg)
}

func (handler *WebserviceHandler) Error(w http.ResponseWriter, err error) {
	handler.Logger.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func (handler *WebserviceHandler) StartSession(ctx context.Context, w http.ResponseWriter, r *http.Request, user *usecases.User) error {
	session, err := handler.store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Values["User"] = user
	session.Save(r, w)

	return nil
}

func (handler *WebserviceHandler) EndSession(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	session, err := handler.store.Get(r, sessionName)
	if err != nil {
		return err
	}
	session.Values["User"] = nil
	session.Save(r, w)

	return nil
}

func (handler *WebserviceHandler) AddHandler(url string, authenticated bool, fn func(context.Context, http.ResponseWriter, *http.Request)) {

	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		// init the main context
		ctx := context.Background()

		if authenticated {
			session, err := handler.store.Get(r, sessionName)
			if err != nil {
				handler.Error(w, err)
				return
			}
			if session.Values["User"] == nil {
				handler.SendJson(w, AuthNeededError{Response: "AUTH_NEEDED_ERROR"})
				return
			}
			user, ok := session.Values["User"].(*usecases.User)
			if !ok {
				handler.SendJson(w, AuthNeededError{Response: "AUTH_NEEDED_ERROR"})
				return
			}

			ctx = context.WithValue(ctx, interfaces.KeyUser, user)
		}

		fn(ctx, w, r)
	})

}
