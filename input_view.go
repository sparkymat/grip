package grip

import (
	"bytes"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type InputView struct {
	TextView TextView
}

func (t *InputView) Initialize(register func(eventType event.Type, handler EventHandler), emit func(eventType event.Type, data interface{}) error) {
	register(event.SystemKeyPress, t)

	t.TextView.Initialize(register, emit)
}

func (t *InputView) Resize(x, y, width, height uint32) {
	t.TextView.Resize(x, y, width, height)
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

func (t *InputView) OnLoad() {
}
