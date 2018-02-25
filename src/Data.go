package radio

const MaximumClients = 5

type Message struct {
	Kind  string
	Value string
}

type PlayingSong struct {
	Id       string
	Position int
}
