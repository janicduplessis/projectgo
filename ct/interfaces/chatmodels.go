package interfaces

type SendMessageRequest struct {
	Message string
}

type JoinChannelRequest struct {
	ChannelId int64
}

type CreateChannelRequest struct {
	Name string
}

type ChannelsResponse struct {
	List    []*ChannelModel
	Current int64
}

type ClientModel struct {
	Id   int64
	Name string
}

type ChannelModel struct {
	Id      int64
	Name    string
	Clients []ClientModel
}

type SendMessageResponse struct {
	Body      string
	Author    string
	ChannelId int64
	ClientId  int64
}
