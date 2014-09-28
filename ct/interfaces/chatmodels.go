package interfaces

// Models
type ClientModel struct {
	Id   int64
	Name string
}

type MessageModel struct {
	Author    string
	Body      string
	UnixTime  int64
	ChannelId int64
	ClientId  int64
}

type ChannelModel struct {
	Id      int64
	Name    string
	Clients []ClientModel
}

// Requests
type SendMessageRequest struct {
	Message string
}

type JoinChannelRequest struct {
	ChannelId int64
}

type CreateChannelRequest struct {
	Name string
}

// Responses
type ChannelsResponse struct {
	List    []*ChannelModel
	Current int64
}

type JoinChannelResponse struct {
	Messages []MessageModel
	Result   bool
}

type ChannelJoinedResponse struct {
	ChannelId int64
	Client    ClientModel
}
