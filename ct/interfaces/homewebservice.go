package interfaces

import (
	"fmt"
	"net/http"
	"os"

	"code.google.com/p/go.net/context"

	"github.com/janicduplessis/projectgo/ct/domain"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

const (
	urlGetProfileModel       string = "/models/getProfileModel"
	urlGetClientProfileModel string = "/models/getClientProfileModel"
	urlGetProfileImage       string = "/getProfileImage"
	urlSetProfileImage       string = "/setProfileImage"
)

type HomeInteractor interface {
	GetClient(clientId int64) (*domain.Client, error)
}

type HomeWebserviceHandler struct {
	Webservice     Webservice
	HomeInteractor HomeInteractor
	ImageUtils     ImageUtils
}

func NewHomeWebservice(ws Webservice, hi HomeInteractor, imageUtils ImageUtils) *HomeWebserviceHandler {
	wsHandler := &HomeWebserviceHandler{
		Webservice:     ws,
		HomeInteractor: hi,
		ImageUtils:     imageUtils,
	}

	ws.AddHandler(urlGetProfileModel, true, wsHandler.GetProfileModel)
	ws.AddHandler(urlGetClientProfileModel, true, wsHandler.GetClientProfileModel)
	ws.AddHandler(urlGetProfileImage, true, wsHandler.GetProfileImage)
	ws.AddHandler(urlSetProfileImage, true, wsHandler.SetProfileImage)

	return wsHandler
}

func (handler *HomeWebserviceHandler) GetProfileModel(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value(KeyUser).(*usecases.User)
	client, err := handler.HomeInteractor.GetClient(user.Id)
	if err != nil {
		handler.Webservice.Error(w, err)
	}

	model := &profileModel{
		Username:     user.Username,
		DisplayName:  client.DisplayName,
		FirstName:    client.FirstName,
		LastName:     client.LastName,
		Email:        client.Email,
		ProfileImage: fmt.Sprintf("/getProfileImage?clientId=%d", client.Id),
	}
	response := ModelResponse{
		Model: model,
	}
	handler.Webservice.SendJson(w, response)
}

func (handler *HomeWebserviceHandler) GetClientProfileModel(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	request := new(getClientProfileModelRequest)
	err := handler.Webservice.ReadJson(w, r, request)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	client, err := handler.HomeInteractor.GetClient(request.ClientId)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	model := &getClientProfileModelResponse{
		Username:     client.DisplayName,
		ProfileImage: fmt.Sprintf("/getProfileImage?clientId=%d", client.Id),
	}
	response := ModelResponse{
		Model: model,
	}
	handler.Webservice.SendJson(w, response)
}

func (handler *HomeWebserviceHandler) GetProfileImage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Get request with query params
	clientId := r.FormValue("clientId")
	imageUrl := fmt.Sprintf("upload/profile_%s.png", clientId)
	_, err := os.Stat(imageUrl)
	if os.IsNotExist(err) {
		imageUrl = "images/no-profile.png"
	}

	http.ServeFile(w, r, imageUrl)
}

func (handler *HomeWebserviceHandler) SetProfileImage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value(KeyUser).(*usecases.User)

	file, hndl, err := r.FormFile("File")
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	handler.Webservice.Log(fmt.Sprintf("Upload file %s", hndl.Filename))

	image, err := handler.ImageUtils.Load(file)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	image = handler.ImageUtils.Resize(image, 192, 192)

	err = handler.ImageUtils.Save(image, fmt.Sprintf("upload/profile_%d.png", user.Id))
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	handler.Webservice.SendJson(w, setProfileImageResponse{Result: true})
}
