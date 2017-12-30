package grip

import (
	"github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type TextAlignment int

const TextAlignmentLeft TextAlignment = 0
const TextAlignmentCenter TextAlignment = 1
const TextAlignmentRight TextAlignment = 2

type TextView struct {
	emitEvent       func(eventType event.Type, data interface{})
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	TextAlignment   TextAlignment
	Text            string
	Area            Area
	height          int
	width           int
	x               int
	y               int
	visibleHeight   int
	visibleWidth    int
	visibleX        int
	visibleY        int
}

func (t *TextView) Initialize(emit func(eventType event.Type, data interface{})) {
	t.emitEvent = emit
}

func (t *TextView) Resize(x, y, width, height, visibleX, visibleY, visibleWidth, visibleHeight int) {
	t.x = x
	t.y = y
	t.width = width
	t.height = height
	t.visibleX = visibleX
	t.visibleY = visibleY
	t.visibleWidth = visibleWidth
	t.visibleHeight = visibleHeight
}

func (t *TextView) Draw() {
	for j := t.y; j <= (t.y + t.height - 1); j++ {
		startPosition := (j - t.y) * t.width
		endPosition := startPosition + t.width - 1
		if endPosition > len(t.Text)-1 {
			endPosition = len(t.Text) - 1
		}

		var line string
		if startPosition < len(t.Text) {
			line = t.Text[startPosition : endPosition+1]
		}

		if len(line) < t.width {
			for i := t.x; i <= (t.x + t.width - 1); i++ {
				termbox.SetCell(i, j, ' ', t.ForegroundColor, t.BackgroundColor)
			}
		}

		if len(line) > 0 {
			textStart := t.x
			textEnd := textStart + len(line) - 1

			if t.TextAlignment == TextAlignmentCenter {
				textStart = t.x + (t.width-len(line))/2
				textEnd = textStart + len(line)
			} else if t.TextAlignment == TextAlignmentRight {
				textEnd = t.x + t.width - 1
				textStart = textEnd - len(line) + 1
			}

			for i := textStart; i <= textEnd; i++ {
				if i-textStart < len(line) {
					char := rune(line[i-textStart])
					termbox.SetCell(i, j, char, t.ForegroundColor, t.BackgroundColor)
				}
			}
		}
	}
}

func (t *TextView) OnEvent(app *App, e event.Event) {
}
