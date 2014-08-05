package ct

type Message struct {
	Id     int
	Author string
	Body   string
}

func (m *Message) String() string {
	return m.Author + ": " + m.Body
}
