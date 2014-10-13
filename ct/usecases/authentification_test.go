package usecases

import (
	"log"
	"testing"

	"github.com/janicduplessis/projectgo/ct/domain"
)

type FakeUserRepository struct {
	Users     map[int64]*User
	Passwords map[int64]string
}

var id int64 = 1

func NewFakeUserRepository() *FakeUserRepository {
	users := make(map[int64]*User)
	pass := make(map[int64]string)
	return &FakeUserRepository{
		Users:     users,
		Passwords: pass,
	}
}

func (repo *FakeUserRepository) Create(info *RegisterInfo) (*User, error) {
	for _, user := range repo.Users {
		if user.Username == info.Username {
			return nil, ErrUserAlreadyExists
		}
	}
	client := &domain.Client{
		Id:          id,
		FirstName:   info.FirstName,
		LastName:    info.LastName,
		Email:       info.Email,
		DisplayName: info.Username,
	}
	user := &User{
		Id:       id,
		Username: info.Username,
		Client:   client,
	}
	id++

	repo.Users[id] = user
	repo.Passwords[id] = info.Password

	return repo.Users[id], nil
}

func (repo *FakeUserRepository) CreateGoogle(info *GoogleRegisterInfo) (*User, error) {
	return nil, nil
}

func (repo *FakeUserRepository) FindByName(name string) (*User, error) {
	for _, user := range repo.Users {
		if user.Username == name {
			return user, nil
		}
	}
	return nil, nil
}

func (repo *FakeUserRepository) FindByNameWithHash(name string) (*User, string, error) {
	for id, user := range repo.Users {
		if user.Username == name {
			return user, repo.Passwords[id], nil
		}
	}
	return nil, "", nil
}

func (repo *FakeUserRepository) FindByGoogleId(id string) (*User, error) {
	return nil, nil
}

func (repo *FakeUserRepository) UpdatePasswordHash(user *User, hash string) error {
	panic("UpdatePasswordHash not implemented")
}

type FakeCryptoHandler struct {
}

func (ch *FakeCryptoHandler) CompareHashAndPassword(hash string, password string) (bool, error) {
	return hash == password, nil
}

func (ch *FakeCryptoHandler) GenerateFromPassword(password string) (string, error) {
	return password, nil
}

type FakeLoggerHandler struct {
}

func (handler *FakeLoggerHandler) Log(message string) {
	log.Println(message)
}

func (handler *FakeLoggerHandler) Error(err error) {
	log.Println(err.Error())
}

func InitAuthentificationInteractor() *AuthentificationInteractor {
	authInteractor := NewAuthentificationInteractor(NewFakeUserRepository(), new(FakeCryptoHandler), new(FakeLoggerHandler))

	//Fake data
	authInteractor.UserRepository.Create(&RegisterInfo{
		Username:  "Test",
		Password:  "test1234",
		FirstName: "Mr",
		LastName:  "Test",
		Email:     "test@test.com",
	})

	return authInteractor
}

// TEST CASES
func TestLogin(t *testing.T) {
	authInteractor := InitAuthentificationInteractor()
	// Good login
	info := &LoginInfo{
		Username: "Test",
		Password: "test1234",
	}
	user, err := authInteractor.Login(info)
	if user == nil {
		t.Errorf("user is nil, want user with name %s", info.Username)
	}
	if err != nil {
		t.Errorf("error is not nil, want nil. Error: %s", err.Error())
	}

	// Bad login
	info = &LoginInfo{
		Username: "Test",
		Password: "waaaaaaaaaaaaaaa",
	}
	user, err = authInteractor.Login(info)
	if user != nil {

	}
	if err != nil {
		t.Errorf("error is not nil, want nil. Error: %s", err.Error())
	}
}

func TestRegister(t *testing.T) {
	authInteractor := InitAuthentificationInteractor()

	// Good register
	info := &RegisterInfo{
		Username:  "MrPotato",
		Password:  "1234test",
		FirstName: "Mr",
		LastName:  "Potato",
		Email:     "potato@test.com",
	}

	user, err := authInteractor.Register(info)
	if user == nil {
		t.Errorf("user is nil, want user with name %s", info.Username)
	}
	if err != nil {
		t.Errorf("error is not nil, want nil. Error: %s", err.Error())
	}

	// Bad, user already exists
	info = &RegisterInfo{
		Username:  "Test",
		Password:  "1234test",
		FirstName: "Mr",
		LastName:  "Potato",
		Email:     "potato@test.com",
	}

	user, err = authInteractor.Register(info)
	if user != nil {
		t.Errorf("user is not nil, returned user with name %s, want nil", user.Username)
	}
	if err == nil {
		t.Errorf("error is nil, want ErrUserAlreadyExists")
	} else if err != ErrUserAlreadyExists {
		t.Errorf("error is %s, want ErrUserAlreadyExists", err.Error())
	}

	// Bad, invalid info
	// Username too short
	info = &RegisterInfo{
		Username:  "Te",
		Password:  "1234test",
		FirstName: "Mr",
		LastName:  "Potato",
		Email:     "potato@test.com",
	}

	user, err = authInteractor.Register(info)
	if user != nil {
		t.Errorf("user is not nil, returned user with name %s, want nil", user.Username)
	}
	if err == nil {
		t.Errorf("error is nil, want ErrInvalidRegisterInfo")
	} else if err != ErrInvalidRegisterInfo {
		t.Errorf("error is %s, want ErrInvalidRegisterInfo", err.Error())
	}

	// Invalid email
	info = &RegisterInfo{
		Username:  "IHaveBadEmail",
		Password:  "1234test",
		FirstName: "Mr",
		LastName:  "Potato",
		Email:     "potatotest.com",
	}

	user, err = authInteractor.Register(info)
	if user != nil {
		t.Errorf("user is not nil, returned user with name %s, want nil", user.Username)
	}
	if err == nil {
		t.Errorf("error is nil, want ErrInvalidRegisterInfo")
	} else if err != ErrInvalidRegisterInfo {
		t.Errorf("error is %s, want ErrInvalidRegisterInfo", err.Error())
	}

	// Missing field
	info = &RegisterInfo{
		Username:  "IForgotSomething",
		Password:  "1234test",
		FirstName: "Mr",
		Email:     "potato@test.com",
	}

	user, err = authInteractor.Register(info)
	if user != nil {
		t.Errorf("user is not nil, returned user with name %s, want nil", user.Username)
	}
	if err == nil {
		t.Errorf("error is nil, want ErrInvalidRegisterInfo")
	} else if err != ErrInvalidRegisterInfo {
		t.Errorf("error is %s, want ErrInvalidRegisterInfo", err.Error())
	}
}
