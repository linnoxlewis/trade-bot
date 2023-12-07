package telegram

type Type int

const (
	Unknown Type = iota
	Message
)

type Event struct {
	Type Type
	Text string
	Meta interface{}
}
