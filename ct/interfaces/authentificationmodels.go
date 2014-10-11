package interfaces

// Models
type UserModel struct {
	Id          int64
	Username    string
	DisplayName string
	FirstName   string
	LastName    string
	Email       string
}

type LoginModel struct {
	GoogleLoginURL string
}

// Responses
type LoginResponseModel struct {
	Result bool
	User   *UserModel
}

type RegisterResponseModel struct {
	Result bool
	User   *UserModel
	Error  string
}

type LogoutResponseModel struct {
	Result bool
}
