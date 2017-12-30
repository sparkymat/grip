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
	height          int
	width           int
	x               int
	y               int
	progessX        int
	speedX          int
}

func (a *ActivityView) Initialize(emit func(eventType event.Type, data interface{})) {
	a.emitEvent = emit

	a.progessX = a.x + 1
	a.speedX = 10

	timer := time.NewTicker(time.Millisecond * 100)
	go func() {
		for _ = range timer.C {
			a.progessX += (a.width / a.speedX)
			if a.progessX >= a.x+a.width-1 {
				a.progessX = a.x + a.width - 2
				a.speedX *= -1
			} else if a.progessX <= a.x+1 {
				a.progessX = a.x + 1
				a.speedX *= -1
			}
			a.Draw()
		}
	}()
}

func (a *ActivityView) Resize(x, y, width, height int) {
	a.x = x
	a.y = y
	a.width = width
	a.height = height
}

func (a *ActivityView) Draw() {
	for j := a.y; j <= (a.y + a.height - 1); j++ {
		termbox.SetCell(a.x, j, '[', a.ForegroundColor, a.BackgroundColor)
		termbox.SetCell(a.x+a.width-1, j, ']', a.ForegroundColor, a.BackgroundColor)

		for i := a.x + 1; i < (a.x + a.width - 1); i++ {
			termbox.SetCell(i, j, ' ', a.ForegroundColor, a.BackgroundColor)
		}
		if a.progessX >= a.x && a.progessX <= a.x+a.width-1 {
			termbox.SetCell(a.progessX, j, '=', a.ForegroundColor, a.BackgroundColor)
		}

		if len(a.Text) > 0 {
			displayText := fmt.Sprintf(" %v ", a.Text)
			textX := a.x + (a.width-len(displayText))/2
			for i := textX; i < textX+len(displayText); i++ {
				char := rune(displayText[i-textX])
				termbox.SetCell(i, j, char, a.ForegroundColor, a.BackgroundColor)
			}
		}
	}
}

func (a *ActivityView) OnEvent(app *App, e event.Event) {
}
