package grip

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type ProgressViewType int

const ProgressViewTypePercentage ProgressViewType = 0
const ProgressViewTypeFraction ProgressViewType = 1

type ProgressView struct {
	emitEvent       func(event.Type, interface{})
	MinimumValue    int
	CurrentValue    int
	MaximumValue    int
	Type            ProgressViewType
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	height          int
	width           int
	x               int
	y               int
}

func (p *ProgressView) Initialize(emit func(eventType event.Type, data interface{})) {
	p.emitEvent = emit
}

func (p *ProgressView) Resize(x, y, width, height int) {
	p.x = x
	p.y = y
	p.width = width
	p.height = height
}

func (p *ProgressView) Draw() {
	for j := p.y; j <= (p.y + p.height - 1); j++ {
		termbox.SetCell(int(p.x), int(j), '[', p.ForegroundColor, p.BackgroundColor)
		for i := p.x + 1; i < (p.x + p.width - 1); i++ {
			termbox.SetCell(int(i), int(j), ' ', p.ForegroundColor, p.BackgroundColor)
		}
		termbox.SetCell(int(p.x+p.width-1), int(j), ']', p.ForegroundColor, p.BackgroundColor)

		var fractionComplete float32 = 0.0
		if p.MaximumValue != p.MinimumValue {
			fractionComplete = float32(p.CurrentValue) / float32(p.MaximumValue-p.MinimumValue)
			minX := p.x + 1
			maxX := p.x + p.width - 2
			currentMaxX := minX + int(float32(maxX-minX)*fractionComplete)
			for i := minX; i <= currentMaxX; i++ {
				termbox.SetCell(i, j, '=', p.ForegroundColor, p.BackgroundColor)
			}
		}

		var displayText string = ""

		switch p.Type {
		case ProgressViewTypeFraction:
			displayText = fmt.Sprintf(" %v/%v ", p.CurrentValue, p.MaximumValue)
			break
		case ProgressViewTypePercentage:
			displayText = fmt.Sprintf(" %.1f%% ", fractionComplete*100)
			break
		}
		textX := p.x + (p.width-len(displayText))/2
		for i := textX; i < textX+len(displayText); i++ {
			char := rune(displayText[i-textX])
			termbox.SetCell(i, j, char, p.ForegroundColor, p.BackgroundColor)
		}
	}
}

func (p *ProgressView) OnEvent(app *App, e event.Event) {
}
