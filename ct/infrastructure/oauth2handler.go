package infrastructure

import (
	"errors"
	"log"
	"net/http"

	"code.google.com/p/go.net/context"
	"code.google.com/p/google-api-go-client/plus/v1"
	"github.com/golang/oauth2"

	"github.com/janicduplessis/projectgo/ct/config"
	"github.com/janicduplessis/projectgo/ct/usecases"
)

const (
	authUrl        = "https://accounts.google.com/o/oauth2/auth"
	tokenUrl       = "https://accounts.google.com/o/oauth2/token"
	profileInfoURL = "https://www.googleapis.com/plus/v1/people/me"
)

var scopes = []string{plus.UserinfoEmailScope}

type OAuth2Handler struct {
	config *oauth2.Config
}

func (handler *OAuth2Handler) Init() {
	conf, err := oauth2.NewConfig(&oauth2.Options{
		ClientID:     config.OAuth2ClientId,
		ClientSecret: config.OAuth2ClientSecret,
		RedirectURL:  config.SiteUrl,
		Scopes:       scopes,
	},
		authUrl,
		tokenUrl)
	if err != nil {
		log.Fatal(err)
	}

	handler.config = conf
}

func (handler *OAuth2Handler) GetProfile(ctx context.Context, w http.ResponseWriter, r *http.Request) (*usecases.GoogleRegisterInfo, error) {
	// Get the auth code and create a transport object
	authCode := r.FormValue("Code")
	log.Println(authCode)
	t, err := handler.config.NewTransportWithCode(authCode)
	if err != nil {
		return nil, err
	}
	// Use the transport object to make request with an http.Client
	client := http.Client{Transport: t}

	// Create the plus service object with which we can make an API call
	service, err := plus.New(&client)
	plusUser, err := service.People.Get("me").Do()
	if err != nil {
		return nil, err
	}

	var email string
	if len(plusUser.Emails) > 0 {
		email = plusUser.Emails[0].Value
	} else {
		return nil, errors.New("The google profile has no emails")
	}

	// Get what we need from the user profile
	profile := &usecases.GoogleRegisterInfo{
		Id:          plusUser.Id,
		DisplayName: plusUser.DisplayName,
		FirstName:   plusUser.Name.GivenName,
		LastName:    plusUser.Name.FamilyName,
		Email:       email,
	}
	return profile, nil
}

func (handler *OAuth2Handler) GetScope() string {
	return plus.UserinfoEmailScope
}
