package grip

import termbox "github.com/nsf/termbox-go"

type Layer int

const (
	AppLayer Layer = iota
	ModalLayer
)

type ColoredRune struct {
	Ch              rune
	ForegroundColor termbox.Attribute
	BackgroundColor termbox.Attribute
}

type DrawCellRequest struct {
	Position Point
	Rune     ColoredRune
}
