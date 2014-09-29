package interfaces

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

func NewHomeWebservice(ws Webservice) *HomeWebserviceHandler {
	wsHandler := &HomeWebserviceHandler{
		Webservice: ws,
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
		ProfileImage: "/getProfileImage",
	}
	response := ModelResponse{
		Model: model,
	}
	handler.Webservice.SendJson(w, response)
}

func (handler *HomeWebserviceHandler) GetProfileImage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := ctx.Value(KeyUser).(*usecases.User)
	imageUrl := fmt.Sprintf("upload/profile_%d.png", user.Id)
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

	ext := strings.ToLower(filepath.Ext(hndl.Filename))
	// TODO: support more formats
	if ext != ".png" {
		handler.Webservice.Error(w, errors.New("Invalid file format"))
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("upload/profile_%d.png", user.Id), data, 0777)
	if err != nil {
		handler.Webservice.Error(w, err)
		return
	}

	handler.Webservice.SendJson(w, SetProfileImageResponse{Result: true})
}
