package interfaces

import (
	"fmt"
	"net/http"
	"os"

	"code.google.com/p/go.net/context"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

const (
	urlGetProfileModel string = "/models/getProfileModel"
	urlGetProfileImage string = "/getProfileImage"
	urlSetProfileImage string = "/setProfileImage"
)

type HomeWebserviceHandler struct {
	Webservice Webservice
	ImageUtils ImageUtils
}

type ProfileModel struct {
	Username     string
	DisplayName  string
	FirstName    string
	LastName     string
	Email        string
	ProfileImage string
}

type SetProfileImageResponse struct {
	Result bool
}

func NewHomeWebservice(ws Webservice, imageUtils ImageUtils) *HomeWebserviceHandler {
	wsHandler := &HomeWebserviceHandler{
		Webservice: ws,
		ImageUtils: imageUtils,
	}

	ws.AddHandler(urlGetProfileModel, true, wsHandler.GetProfileModel)
	ws.AddHandler(urlGetProfileImage, true, wsHandler.GetProfileImage)
	ws.AddHandler(urlSetProfileImage, true, wsHandler.SetProfileImage)

	return wsHandler
}

func (handler *HomeWebserviceHandler) GetProfileModel(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value(KeyUser).(*usecases.User)
	model := &ProfileModel{
		Username:     user.Username,
		DisplayName:  user.Client.DisplayName,
		FirstName:    user.Client.FirstName,
		LastName:     user.Client.LastName,
		Email:        user.Client.Email,
		ProfileImage: fmt.Sprintf("/getProfileImage?clientId=%d", user.Id),
	}
	response := ModelResponse{
		Model: model,
	}
	handler.Webservice.SendJson(w, response)
}

func (handler *HomeWebserviceHandler) GetProfileImage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

	handler.Webservice.SendJson(w, SetProfileImageResponse{Result: true})
}
