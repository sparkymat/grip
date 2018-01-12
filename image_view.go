package grip

import (
	"image"

	termbox "github.com/nsf/termbox-go"
	asciiart "github.com/sparkymat/goasciiart"
	"github.com/sparkymat/grip/event"
)

type ImageView struct {
	setCellFn       SetCellFn
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	Image           image.Image
	scaleAscii      []byte
	rect            Rect
	visibleRect     Rect
}

func (i *ImageView) Initialize(setCellFn SetCellFn) {
	i.setCellFn = setCellFn
}

func (i *ImageView) Resize(rect, visibleRect Rect) {
	i.rect = rect
	i.visibleRect = visibleRect

	if i.rect.Size.Width > 0 {
		i.scaleAscii = asciiart.Convert2AsciiOfWidth(i.Image, i.rect.Size.Width-1)
	}
}

func (v *ImageView) Draw() {
	if v.scaleAscii == nil || len(v.scaleAscii) == 0 {
		return
	}

	for j := v.rect.Origin.Y; j <= (v.rect.Origin.Y + v.rect.Size.Height - 1); j++ {
		for i := v.rect.Origin.X + 1; i < (v.rect.Origin.X + v.rect.Size.Width - 1); i++ {
			var r rune = ' '
			position := (j-v.rect.Origin.Y)*v.rect.Size.Width + (i - v.rect.Origin.X)
			if position < len(v.scaleAscii) {
				r = rune(v.scaleAscii[position])
			}
			v.SetCellIfVisible(i, j, r, v.ForegroundColor, v.BackgroundColor)
		}
	}
}

func (i *ImageView) OnEvent(e event.Event) {
}

func (i *ImageView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	if i.visibleRect.Contains(x, y) {
		i.setCellFn(Point{x, y}, ColoredRune{ch, fg, bg})
	}
}
