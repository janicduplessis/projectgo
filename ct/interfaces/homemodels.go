package interfaces

// Models
type ProfileModel struct {
	Username     string
	DisplayName  string
	FirstName    string
	LastName     string
	Email        string
	ProfileImage string
}

// Requests
type GetClientProfileModelRequest struct {
	ClientId int64
}

// Responses
type GetClientProfileModelResponse struct {
	Username     string
	ProfileImage string
}

type SetProfileImageResponse struct {
	Result bool
}
