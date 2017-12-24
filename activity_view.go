package grip

import (
	"fmt"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type ActivityView struct {
	app             *app
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	Text            string
	height          uint32
	width           uint32
	x               uint32
	y               uint32
	progessX        uint32
	speedX          int32
}

func (a *ActivityView) SetApp(app *app) {
	a.app = app
}

func (a *ActivityView) RegisteredEvents() []event.Type {
	return []event.Type{}
}

func (a *ActivityView) Resize(x, y, width, height uint32) {
	a.x = x
	a.y = y
	a.width = width
	a.height = height
}

func (a *ActivityView) Draw() {
	for j := a.y; j <= (a.y + a.height - 1); j++ {
		termbox.SetCell(int(a.x), int(j), '[', a.ForegroundColor, a.BackgroundColor)
		termbox.SetCell(int(a.x+a.width-1), int(j), ']', a.ForegroundColor, a.BackgroundColor)

		for i := a.x + 1; i < (a.x + a.width - 1); i++ {
			termbox.SetCell(int(i), int(j), ' ', a.ForegroundColor, a.BackgroundColor)
		}
		if a.progessX >= a.x && a.progessX <= a.x+a.width-1 {
			termbox.SetCell(int(a.progessX), int(j), '=', a.ForegroundColor, a.BackgroundColor)
		}

		if len(a.Text) > 0 {
			displayText := fmt.Sprintf(" %v ", a.Text)
			textX := a.x + (a.width-uint32(len(displayText)))/2
			for i := textX; i < textX+uint32(len(displayText)); i++ {
				char := rune(displayText[i-textX])
				termbox.SetCell(int(i), int(j), char, a.ForegroundColor, a.BackgroundColor)
			}
		}
	}
}

func (a *ActivityView) OnLoad() {
	a.progessX = a.x
	a.speedX = 10

	timer := time.NewTicker(time.Millisecond * 100)
	go func() {
		for _ = range timer.C {
			a.progessX = uint32(int32(a.progessX) + int32(a.width)/a.speedX)
			if a.progessX >= a.x+a.width {
				a.progessX = a.x + a.width - 1
				a.speedX *= -1
			} else if a.progessX <= a.x {
				a.progessX = a.x
				a.speedX *= -1
			}
			a.Draw()
		}
	}()
}

func (a *ActivityView) OnEvent(e event.Event) {
}
