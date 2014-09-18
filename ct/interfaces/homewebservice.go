package interfaces

import (
	"net/http"

	"code.google.com/p/go.net/context"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

const (
	urlGetProfileModel string = "/models/getProfileModel"
)

type HomeWebserviceHandler struct {
	Webservice Webservice
}

type ProfileModel struct {
	Username    string
	DisplayName string
	FirstName   string
	LastName    string
	Email       string
}

func NewHomeWebservice(ws Webservice) *HomeWebserviceHandler {
	wsHandler := &HomeWebserviceHandler{
		Webservice: ws,
	}

	ws.AddHandler(urlGetProfileModel, true, wsHandler.GetProfileModel)

	return wsHandler
}

func (handler *HomeWebserviceHandler) GetProfileModel(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value(KeyUser).(*usecases.User)
	model := &ProfileModel{
		Username:    user.Username,
		DisplayName: user.Client.DisplayName,
		FirstName:   user.Client.FirstName,
		LastName:    user.Client.LastName,
		Email:       user.Client.Email,
	}
	response := ModelResponse{
		Model: model,
	}
	handler.Webservice.SendJson(w, response)
}
