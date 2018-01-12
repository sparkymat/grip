package event

type Type int

const (
	EventKeyPress Type = iota
	EventResize
	EventTick
	EventDrawCellRequest
	EventShowModal
	EventHideModal
	EventRefresh
)
