package interfaces

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

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
	FileStore      FileStore
}

func NewHomeWebservice(ws Webservice, hi HomeInteractor, imageUtils ImageUtils, fileStore FileStore) *HomeWebserviceHandler {
	wsHandler := &HomeWebserviceHandler{
		Webservice:     ws,
		HomeInteractor: hi,
		ImageUtils:     imageUtils,
		FileStore:      fileStore,
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

	data, err := handler.FileStore.Open(imageUrl)
	if err != nil {
		if err == ErrNoFile {
			data, err = ioutil.ReadFile("images/no-profile.png")
		} else {
			handler.Webservice.Error(w, err)
			return
		}
	}

	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", path.Ext(imageUrl)[1:]))
	w.Write(data)
}

func (handler *HomeWebserviceHandler) SetProfileImage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value(KeyUser).(*usecases.User)

	imgData, hndl, err := r.FormFile("File")
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	handler.Webservice.Log(fmt.Sprintf("Upload file %s", hndl.Filename))

	image, err := handler.ImageUtils.Load(imgData)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	image = handler.ImageUtils.Resize(image, 192, 192)

	data, err := handler.ImageUtils.Save(image, ".png")

	err = handler.FileStore.Create(fmt.Sprintf("upload/profile_%d.png", user.Id), data)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	handler.Webservice.SendJson(w, setProfileImageResponse{Result: true})
}
