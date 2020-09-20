package echoapp

type Event struct {
	Name  string
	ComID uint
}

type EventPaid struct {
	Event
	Order Order
}

type Score struct {
	Event
	UserId       uint
	Score        int
	Source       string
	SourceDetail string
	Note         string
}
