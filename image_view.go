package grip

import (
	"image"

	termbox "github.com/nsf/termbox-go"
	asciiart "github.com/sparkymat/goasciiart"
	"github.com/sparkymat/grip/event"
)

type ImageView struct {
	emitEvent       func(event.Type, interface{})
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	Image           image.Image
	scaleAscii      []byte
	rect            Rect
	visibleRect     Rect
}

func (i *ImageView) Initialize(emit func(eventType event.Type, data interface{})) {
	i.emitEvent = emit
}

func (i *ImageView) Resize(rect, visibleRect Rect) {
	i.rect = rect
	i.visibleRect = visibleRect

	i.scaleAscii = asciiart.Convert2AsciiOfWidth(i.Image, int(i.rect.Width)-1)
}

func (v *ImageView) Draw() {
	for j := v.rect.Y; j <= (v.rect.Y + v.rect.Height - 1); j++ {
		for i := v.rect.X + 1; i < (v.rect.X + v.rect.Width - 1); i++ {
			var r rune = ' '
			if j*v.rect.Width+i < len(v.scaleAscii) {
				r = rune(v.scaleAscii[j*v.rect.Width+i])
			}
			v.SetCellIfVisible(i, j, r, v.ForegroundColor, v.BackgroundColor)
		}
	}
}

func (i *ImageView) OnEvent(app *App, e event.Event) {
}

func (i *ImageView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	if i.visibleRect.Contains(x, y) {
		termbox.SetCell(x, y, ch, fg, bg)
	}
}
