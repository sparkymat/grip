package grip

import (
	"bytes"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type InputView struct {
	drawCell  DrawCellFn
	emitEvent EmitEventFn
	TextView  TextView
}

func (t *InputView) Initialize(drawCell DrawCellFn, emitEvent EmitEventFn) {
	t.drawCell = drawCell
	t.emitEvent = emitEvent

	t.TextView.Initialize(drawCell, emitEvent)
}

func (t *InputView) Resize(rect, visibleRect Rect) {
	t.TextView.Resize(rect, visibleRect)
}

func (t *InputView) Draw() {
	t.TextView.Draw()
}

func (t *InputView) OnEvent(app *App, e event.Event) {
	switch e.Type {
	case event.SystemKeyPress:
		ev := e.Data.(termbox.Event)
		switch ev.Type {
		case termbox.EventKey:
			var buffer bytes.Buffer
			buffer.WriteString(t.TextView.Text)
			buffer.WriteRune(ev.Ch)
			t.TextView.Text = buffer.String()
			t.Draw()
			return
		}
		return
	}
}

func (t *InputView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	t.TextView.SetCellIfVisible(x, y, ch, fg, bg)
}
