package event

type Type int

const (
	EventKeyPress Type = iota
	EventTick
	EventDrawCellRequest
	EventShowModal
	EventHideModal
)
