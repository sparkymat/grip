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
	app             *app
	MinimumValue    int32
	CurrentValue    int32
	MaximumValue    int32
	Type            ProgressViewType
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	height          uint32
	width           uint32
	x               uint32
	y               uint32
}

func (p *ProgressView) SetApp(app *app) {
	p.app = app
}

func (p *ProgressView) RegisteredEvents() []event.Type {
	return []event.Type{}
}

func (p *ProgressView) Resize(x, y, width, height uint32) {
	p.x = x
	p.y = y
	p.width = width
	p.height = height
}

func (p *ProgressView) Draw() {
	for j := p.y; j <= (p.y + p.height - 1); j++ {
		termbox.SetCell(int(p.x), int(j), '[', p.ForegroundColor, p.BackgroundColor)
		termbox.SetCell(int(p.x+p.width-1), int(j), ']', p.ForegroundColor, p.BackgroundColor)

		var fractionComplete float32 = 0.0
		if p.MaximumValue != p.MinimumValue {
			fractionComplete = float32(p.CurrentValue) / float32(p.MaximumValue-p.MinimumValue)
			minX := p.x + 1
			maxX := p.x + p.width - 2
			currentMaxX := minX + uint32(float32(maxX-minX)*fractionComplete)
			for i := minX; i <= currentMaxX; i++ {
				termbox.SetCell(int(i), int(j), '=', p.ForegroundColor, p.BackgroundColor)
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
		textX := p.x + (p.width-uint32(len(displayText)))/2
		for i := textX; i < textX+uint32(len(displayText)); i++ {
			char := rune(displayText[i-textX])
			termbox.SetCell(int(i), int(j), char, p.ForegroundColor, p.BackgroundColor)
		}
	}
}

func (p *ProgressView) OnEvent(e event.Event) {
}
