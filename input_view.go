package grip

import (
	"bytes"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type InputView struct {
	emitEvent func(event.Type, interface{})
	TextView  TextView
}

func (t *InputView) Initialize(emit func(eventType event.Type, data interface{})) {
	t.emitEvent = emit
	t.TextView.Initialize(emit)
}

func (t *InputView) Resize(x, y, width, height, visibleX, visibleY, visibleWidth, visibleHeight int) {
	t.TextView.Resize(x, y, width, height, visibleX, visibleY, visibleWidth, visibleHeight)
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
