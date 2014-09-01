package ct

type Message struct {
	Id        int
	Author    string
	Body      string
	ChannelId int64

	Channel *Channel
}

func (m *Message) String() string {
	return m.Author + ": " + m.Body
}
