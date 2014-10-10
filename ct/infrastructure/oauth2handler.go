package infrastructure

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"code.google.com/p/go.net/context"
	"github.com/golang/oauth2"

	config "github.com/janicduplessis/projectgo/ct/config"
	"github.com/janicduplessis/projectgo/ct/interfaces"
)

const (
	authUrl        = "https://accounts.google.com/o/oauth2/auth"
	tokenUrl       = "https://accounts.google.com/o/oauth2/token"
	profileInfoURL = "https://www.googleapis.com/plus/v1/people/me"
)

var scopes = []string{"profile", "email"}

type OAuth2Handler struct {
	url    string
	config *oauth2.Config
}

func (handler *OAuth2Handler) Init() {
	conf, err := oauth2.NewConfig(&oauth2.Options{
		ClientID:     config.OAuth2ClientId,
		ClientSecret: config.OAuth2ClientSecret,
		RedirectURL:  config.SiteUrl + interfaces.UrlOAuth2Login,
		Scopes:       scopes,
	},
		authUrl,
		tokenUrl)
	if err != nil {
		log.Fatal(err)
	}

	handler.url = conf.AuthCodeURL("state", "online", "auto")
	handler.config = conf
}

func (handler *OAuth2Handler) GetProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) (*interfaces.OAuth2Profile, error) {
	// Get the auth code and create a transport object
	authCode := r.FormValue("code")
	t, err := handler.config.NewTransportWithCode(authCode)
	if err != nil {
		return nil, err
	}
	// Use the transport object to make request with an http.Client
	client := http.Client{Transport: t}
	// Request the user profile, we use google profile api
	response, err := client.Get(profileInfoURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(body))
	userData := make(map[string]interface{})
	err = json.Unmarshal(body, &userData)
	if err != nil {
		return nil, err
	}

	// Get what we need from the user profile
	oauth2Profile := &interfaces.OAuth2Profile{
	/*Id:           userData["id"],
	DisplayName:  userData["name"],
	FirstName:    userData["given_name"],
	LastName:     userData["family_name"],
	ProfileImage: userData["picture"],*/
	}

	return oauth2Profile, nil
}

func (handler *OAuth2Handler) GetUrl() (string, error) {
	return handler.url, nil
}
