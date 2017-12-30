package grip

import (
	"github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type TextAlignment uint32

const TextAlignmentLeft TextAlignment = 0
const TextAlignmentCenter TextAlignment = 1
const TextAlignmentRight TextAlignment = 2

type TextView struct {
	app             *App
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	TextAlignment   TextAlignment
	Text            string
	Area            Area
	height          uint32
	width           uint32
	x               uint32
	y               uint32
}

func (t *TextView) Initialize(register func(eventType event.Type, handler EventHandler), emit func(eventType event.Type, data interface{}) error) {
}

func (t *TextView) Resize(x, y, width, height uint32) {
	t.x = x
	t.y = y
	t.width = width
	t.height = height
}

func (t *TextView) Draw() {
	for j := t.y; j <= (t.y + t.height - 1); j++ {
		startPosition := (j - t.y) * t.width
		endPosition := startPosition + t.width - 1
		if endPosition > uint32(len(t.Text))-1 {
			endPosition = uint32(len(t.Text)) - 1
		}

		var line string
		if startPosition < uint32(len(t.Text)) {
			line = t.Text[startPosition : endPosition+1]
		}

		if uint32(len(line)) < t.width {
			for i := t.x; i <= (t.x + t.width - 1); i++ {
				termbox.SetCell(int(i), int(j), ' ', t.ForegroundColor, t.BackgroundColor)
			}
		}

		if len(line) > 0 {
			textStart := t.x
			textEnd := textStart + uint32(len(line)) - 1

			if t.TextAlignment == TextAlignmentCenter {
				textStart = t.x + (t.width-uint32(len(line)))/2
				textEnd = textStart + uint32(len(line))
			} else if t.TextAlignment == TextAlignmentRight {
				textEnd = t.x + t.width - 1
				textStart = textEnd - uint32(len(line)) + 1
			}

			for i := textStart; i <= textEnd; i++ {
				if i-textStart < uint32(len(line)) {
					char := rune(line[i-textStart])
					termbox.SetCell(int(i), int(j), char, t.ForegroundColor, t.BackgroundColor)
				}
			}
		}
	}
}

func (t *TextView) OnEvent(app *App, e event.Event) {
}

func (t *TextView) OnLoad() {
}
