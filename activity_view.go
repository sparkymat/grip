package grip

import (
	"fmt"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type ActivityView struct {
	setCellFn       SetCellFn
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	Text            string
	progessX        int
	speedX          int
	rect            Rect
	visibleRect     Rect
}

func (a *ActivityView) Initialize(setCellFn SetCellFn) {
	a.setCellFn = setCellFn

	a.progessX = a.rect.Origin.X + 1
	a.speedX = 10

	timer := time.NewTicker(time.Millisecond * 100)
	go func() {
		for _ = range timer.C {
			a.progessX += (a.rect.Size.Width / a.speedX)
			if a.progessX >= a.rect.Origin.X+a.rect.Size.Width-1 {
				a.progessX = a.rect.Origin.X + a.rect.Size.Width - 2
				a.speedX *= -1
			} else if a.progessX <= a.rect.Origin.X+1 {
				a.progessX = a.rect.Origin.X + 1
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
	for j := a.rect.Origin.Y; j <= (a.rect.Origin.Y + a.rect.Size.Height - 1); j++ {
		a.SetCellIfVisible(a.rect.Origin.X, j, '[', a.ForegroundColor, a.BackgroundColor)
		a.SetCellIfVisible(a.rect.Origin.X+a.rect.Size.Width-1, j, ']', a.ForegroundColor, a.BackgroundColor)

		for i := a.rect.Origin.X + 1; i < (a.rect.Origin.X + a.rect.Size.Width - 1); i++ {
			a.SetCellIfVisible(i, j, ' ', a.ForegroundColor, a.BackgroundColor)
		}
		if a.progessX >= a.rect.Origin.X && a.progessX <= a.rect.Origin.X+a.rect.Size.Width-1 {
			a.SetCellIfVisible(a.progessX, j, '=', a.ForegroundColor, a.BackgroundColor)
		}

		if len(a.Text) > 0 {
			displayText := fmt.Sprintf(" %v ", a.Text)
			textX := a.rect.Origin.X + (a.rect.Size.Width-len(displayText))/2
			for i := textX; i < textX+len(displayText); i++ {
				char := rune(displayText[i-textX])
				a.SetCellIfVisible(i, j, char, a.ForegroundColor, a.BackgroundColor)
			}
		}
	}
}

func (a *ActivityView) OnEvent(e event.Event) {
}

func (a *ActivityView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	if a.visibleRect.Contains(x, y) {
		a.setCellFn(Point{x, y}, ColoredRune{ch, fg, bg})
	}
}
