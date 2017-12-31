package grip

import (
	"fmt"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type ActivityView struct {
	emitEvent       func(event.Type, interface{})
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	Text            string
	progessX        int
	speedX          int
	rect            Rect
	visibleRect     Rect
}

func (a *ActivityView) Initialize(emit func(eventType event.Type, data interface{})) {
	a.emitEvent = emit

	a.progessX = a.rect.X + 1
	a.speedX = 10

	timer := time.NewTicker(time.Millisecond * 100)
	go func() {
		for _ = range timer.C {
			a.progessX += (a.rect.Width / a.speedX)
			if a.progessX >= a.rect.X+a.rect.Width-1 {
				a.progessX = a.rect.X + a.rect.Width - 2
				a.speedX *= -1
			} else if a.progessX <= a.rect.X+1 {
				a.progessX = a.rect.X + 1
				a.speedX *= -1
			}
			a.Draw()
		}
	}()
}

func (a *ActivityView) Resize(rect, visibleRect Rect) {
	a.rect = rect
	a.visibleRect = visibleRect
}

func (a *ActivityView) Draw() {
	for j := a.rect.Y; j <= (a.rect.Y + a.rect.Height - 1); j++ {
		termbox.SetCell(a.rect.X, j, '[', a.ForegroundColor, a.BackgroundColor)
		termbox.SetCell(a.rect.X+a.rect.Width-1, j, ']', a.ForegroundColor, a.BackgroundColor)

		for i := a.rect.X + 1; i < (a.rect.X + a.rect.Width - 1); i++ {
			termbox.SetCell(i, j, ' ', a.ForegroundColor, a.BackgroundColor)
		}
		if a.progessX >= a.rect.X && a.progessX <= a.rect.X+a.rect.Width-1 {
			termbox.SetCell(a.progessX, j, '=', a.ForegroundColor, a.BackgroundColor)
		}

		if len(a.Text) > 0 {
			displayText := fmt.Sprintf(" %v ", a.Text)
			textX := a.rect.X + (a.rect.Width-len(displayText))/2
			for i := textX; i < textX+len(displayText); i++ {
				char := rune(displayText[i-textX])
				termbox.SetCell(i, j, char, a.ForegroundColor, a.BackgroundColor)
			}
		}
	}
}

func (a *ActivityView) OnEvent(app *App, e event.Event) {
}
