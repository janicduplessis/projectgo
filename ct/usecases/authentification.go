package usecases

import (
	"github.com/janicduplessis/projectgo/ct/domain"

	"fmt"
	"regexp"
	"strings"
)

// UserRepository interfaces a user repo
type UserRepository interface {
	Create(info *RegisterInfo) (*User, error)
	CreateGoogle(info *GoogleRegisterInfo) (*User, error)
	FindByName(name string) (*User, error)
	FindByNameWithHash(name string) (*User, string, error)
	FindByGoogleId(id string) (*User, error)
	UpdatePasswordHash(user *User, hash string) error
}

// Crypto interfaces a crypto handler
type Crypto interface {
	CompareHashAndPassword(hash string, password string) (bool, error)
	GenerateFromPassword(password string) (string, error)
}

// User represents an authentified user
type User struct {
	Id       int64
	Username string
	Client   *domain.Client
}

// RegisterInfo info used during registration
type RegisterInfo struct {
	Username  string
	Password  string
	FirstName string
	LastName  string
	Email     string
}

// GoogleRegisterInfo info for google login/registration
type GoogleRegisterInfo struct {
	Id          string
	DisplayName string
	FirstName   string
	LastName    string
	Email       string
}

// LoginInfo info for normal login
type LoginInfo struct {
	Username string
	Password string
}

// AuthentificationInteractor use cases for authentification
type AuthentificationInteractor struct {
	UserRepository UserRepository
	Crypto         Crypto
	Logger         Logger

	emailRegexp *regexp.Regexp
}

// NewAuthentificationInteractor ctor
func NewAuthentificationInteractor(userRepository UserRepository, crypto Crypto, logger Logger) *AuthentificationInteractor {
	emailRegexp := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,6}$`)

	return &AuthentificationInteractor{
		UserRepository: userRepository,
		Crypto:         crypto,
		Logger:         logger,
		emailRegexp:    emailRegexp,
	}
}

// Login logs in the user
func (ai *AuthentificationInteractor) Login(info *LoginInfo) (*User, error) {
	ai.Logger.Log(fmt.Sprintf("Login request. Username: %s", info.Username))

	user, hash, err := ai.UserRepository.FindByNameWithHash(info.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		ai.Logger.Log(fmt.Sprintf("Bad username. Username: %s", info.Username))
		return nil, nil
	}

	ok, err := ai.Crypto.CompareHashAndPassword(hash, info.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		ai.Logger.Log(fmt.Sprintf("Bad password. Username: %s", info.Username))
		return nil, nil
	}

	ai.Logger.Log(fmt.Sprintf("Successful login. Username: %s", info.Username))

	return user, nil
}

// Register registers the user
func (ai *AuthentificationInteractor) Register(info *RegisterInfo) (*User, error) {
	info.Email = strings.ToLower(info.Email)

	// Validate infos
	if len(info.Username) < 4 || len(info.Username) > 40 {
		return nil, ErrInvalidRegisterInfo
	}
	if len(info.Password) < 6 || len(info.Password) > 40 {
		return nil, ErrInvalidRegisterInfo
	}
	if len(info.FirstName) < 1 || len(info.FirstName) > 40 {
		return nil, ErrInvalidRegisterInfo
	}
	if len(info.LastName) < 1 || len(info.LastName) > 40 {
		return nil, ErrInvalidRegisterInfo
	}
	if len(info.Email) < 1 || len(info.Email) > 40 {
		return nil, ErrInvalidRegisterInfo
	}
	// Simple email regex
	if !ai.emailRegexp.MatchString(info.Email) {
		return nil, ErrInvalidRegisterInfo
	}

	// Hash the password
	hash, err := ai.Crypto.GenerateFromPassword(info.Password)
	if err != nil {
		return nil, err
	}

	info.Password = hash

	return ai.UserRepository.Create(info)
}

// GoogleLogin logs the user using google+ profile
func (ai *AuthentificationInteractor) GoogleLogin(info *GoogleRegisterInfo) (*User, error) {
	user, err := ai.UserRepository.FindByGoogleId(info.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user, err = ai.UserRepository.CreateGoogle(info)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
