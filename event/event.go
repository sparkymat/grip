package event

type Event struct {
	Type Type
	Data interface{}
}

type EventHandler interface {
	OnEvent(e Event)
}
