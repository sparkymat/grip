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
	emitEvent       EmitEventFn
	drawCell        DrawCellFn
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	TextAlignment   TextAlignment
	Text            string
	Area            Area
	rect            Rect
	visibleRect     Rect
}

func (t *TextView) Initialize(drawCell DrawCellFn, emitEvent EmitEventFn) {
	t.drawCell = drawCell
	t.emitEvent = emitEvent
}

func (t *TextView) Resize(rect, visibleRect Rect) {
	t.rect = rect
	t.visibleRect = visibleRect
}

func (t *TextView) Draw() {
	for j := t.rect.Y; j <= (t.rect.Y + t.rect.Height - 1); j++ {
		startPosition := (j - t.rect.Y) * t.rect.Width
		endPosition := startPosition + t.rect.Width - 1
		if endPosition > len(t.Text)-1 {
			endPosition = len(t.Text) - 1
		}

		line := ""
		if startPosition < len(t.Text) {
			line = t.Text[startPosition : endPosition+1]
		}

		if len(line) < t.rect.Width {
			for i := t.rect.X; i <= (t.rect.X + t.rect.Width - 1); i++ {
				t.SetCellIfVisible(i, j, ' ', t.ForegroundColor, t.BackgroundColor)
			}
		}

		if len(line) > 0 {
			textStart := t.rect.X
			textEnd := textStart + len(line) - 1

			if t.TextAlignment == TextAlignmentCenter {
				textStart = t.rect.X + (t.rect.Width-len(line))/2
				textEnd = textStart + len(line)
			} else if t.TextAlignment == TextAlignmentRight {
				textEnd = t.rect.X + t.rect.Width - 1
				textStart = textEnd - len(line) + 1
			}

			for i := textStart; i <= textEnd; i++ {
				if i-textStart < len(line) {
					char := rune(line[i-textStart])
					t.SetCellIfVisible(i, j, char, t.ForegroundColor, t.BackgroundColor)
				}
			}
		}
	}
}

func (t *TextView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	if t.visibleRect.Contains(x, y) {
		t.drawCell(x, y, ch, fg, bg)
	}
}

func (t *TextView) OnEvent(app *App, e event.Event) {
}
