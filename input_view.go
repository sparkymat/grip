package grip

import (
	"bytes"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type InputView struct {
	app      *App
	layer    Layer
	TextView TextView
	enabled  bool
}

func (t *InputView) Initialize(app *App, layer Layer) {
	t.app = app
	t.layer = layer

	t.TextView.Initialize(app, layer)
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
			switch ev.Key {
			case termbox.KeyEsc:
				t.Disable()
			default:
				if ev.Ch != 0 && t.enabled {
					var buffer bytes.Buffer
					buffer.WriteString(t.TextView.Text)
					buffer.WriteRune(ev.Ch)
					t.TextView.Text = buffer.String()
				}
			}
			t.Draw()
			return
		}
		return
	}
}

func (t *InputView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	t.TextView.SetCellIfVisible(x, y, ch, fg, bg)
}

func (t *InputView) Enable() {
	t.enabled = true
}

func (t *InputView) Disable() {
	t.enabled = false
}
