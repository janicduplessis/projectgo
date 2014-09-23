package interfaces

// Models
type ClientModel struct {
	Id   int64
	Name string
}

type MessageModel struct {
	Author string
	Body   string
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
type SendMessageResponse struct {
	Body      string
	Author    string
	ChannelId int64
	ClientId  int64
}

type ChannelsResponse struct {
	List    []*ChannelModel
	Current int64
}

type JoinChannelResponse struct {
	Messages []MessageModel
	Result   bool
}
