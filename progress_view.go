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
	setCellFn       SetCellFn
	MinimumValue    int
	CurrentValue    int
	MaximumValue    int
	Type            ProgressViewType
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	rect            Rect
	visibleRect     Rect
}

func (p *ProgressView) Initialize(setCellFn SetCellFn) {
	p.setCellFn = setCellFn
}

func (p *ProgressView) Resize(rect, visibleRect Rect) {
	p.rect = rect
	p.visibleRect = visibleRect
}

func (p *ProgressView) Draw() {
	for j := p.rect.Origin.Y; j <= (p.rect.Origin.Y + p.rect.Size.Height - 1); j++ {
		p.SetCellIfVisible(p.rect.Origin.X, j, '[', p.ForegroundColor, p.BackgroundColor)
		for i := p.rect.Origin.X + 1; i < (p.rect.Origin.X + p.rect.Size.Width - 1); i++ {
			p.SetCellIfVisible(i, j, ' ', p.ForegroundColor, p.BackgroundColor)
		}
		p.SetCellIfVisible(p.rect.Origin.X+p.rect.Size.Width-1, j, ']', p.ForegroundColor, p.BackgroundColor)

		var fractionComplete float32 = 0.0
		if p.MaximumValue != p.MinimumValue {
			fractionComplete = float32(p.CurrentValue) / float32(p.MaximumValue-p.MinimumValue)
			minX := p.rect.Origin.X + 1
			maxX := p.rect.Origin.X + p.rect.Size.Width - 2
			currentMaxX := minX + int(float32(maxX-minX)*fractionComplete)
			for i := minX; i <= currentMaxX; i++ {
				p.SetCellIfVisible(i, j, '=', p.ForegroundColor, p.BackgroundColor)
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
		textX := p.rect.Origin.X + (p.rect.Size.Width-len(displayText))/2
		for i := textX; i < textX+len(displayText); i++ {
			char := rune(displayText[i-textX])
			p.SetCellIfVisible(i, j, char, p.ForegroundColor, p.BackgroundColor)
		}
	}
}

func (p *ProgressView) OnEvent(e event.Event) {
}

func (p *ProgressView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	if p.visibleRect.Contains(x, y) {
		p.setCellFn(Point{x, y}, ColoredRune{ch, fg, bg})
	}
}
