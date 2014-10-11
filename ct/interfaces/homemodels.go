package interfaces

// Models
type profileModel struct {
	Username     string
	DisplayName  string
	FirstName    string
	LastName     string
	Email        string
	ProfileImage string
}

// Requests
type getClientProfileModelRequest struct {
	ClientId int64
}

// Responses
type getClientProfileModelResponse struct {
	Username     string
	ProfileImage string
}

type setProfileImageResponse struct {
	Result bool
}
