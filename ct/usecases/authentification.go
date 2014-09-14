package usecases

import (
	"github.com/janicduplessis/projectgo/ct/domain"

	"fmt"
)

type UserRepository interface {
	Create(info *RegisterInfo) (*User, error)
	FindByName(name string) (*User, error)
	FindByNameWithHash(name string) (*User, string, error)
	UpdatePasswordHash(user *User, hash string) error
}

type Crypto interface {
	CompareHashAndPassword(hash string, password string) (bool, error)
	GenerateFromPassword(password string) (string, error)
}

type User struct {
	Id       int64
	Username string
	Client   *domain.Client
}

type RegisterInfo struct {
	Username  string
	Password  string
	FirstName string
	LastName  string
	Email     string
}

type LoginInfo struct {
	Username string
	Password string
}

type AuthentificationInteractor struct {
	UserRepository UserRepository
	Crypto         Crypto
	Logger         Logger
}

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

func (ai *AuthentificationInteractor) Register(info *RegisterInfo) (*User, error) {
	// Hash the password
	hash, err := ai.Crypto.GenerateFromPassword(info.Password)
	if err != nil {
		return nil, err
	}

	info.Password = hash

	// TODO: Validate infos

	return ai.UserRepository.Create(info)
}
