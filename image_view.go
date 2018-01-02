package grip

import (
	"image"

	termbox "github.com/nsf/termbox-go"
	asciiart "github.com/sparkymat/goasciiart"
	"github.com/sparkymat/grip/event"
)

type ImageView struct {
	app             *App
	layer           Layer
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	Image           image.Image
	scaleAscii      []byte
	rect            Rect
	visibleRect     Rect
}

func (i *ImageView) Initialize(app *App, layer Layer) {
	i.app = app
	i.layer = layer
}

func (i *ImageView) Resize(rect, visibleRect Rect) {
	i.rect = rect
	i.visibleRect = visibleRect

	if i.rect.Width > 0 {
		i.scaleAscii = asciiart.Convert2AsciiOfWidth(i.Image, int(i.rect.Width)-1)
	}
}

func (v *ImageView) Draw() {
	if v.scaleAscii == nil || len(v.scaleAscii) == 0 {
		return
	}

	for j := v.rect.Y; j <= (v.rect.Y + v.rect.Height - 1); j++ {
		for i := v.rect.X + 1; i < (v.rect.X + v.rect.Width - 1); i++ {
			var r rune = ' '
			position := (j-v.rect.Y)*v.rect.Width + (i - v.rect.X)
			if position < len(v.scaleAscii) {
				r = rune(v.scaleAscii[position])
			}
			v.SetCellIfVisible(i, j, r, v.ForegroundColor, v.BackgroundColor)
		}
	}
}

func (i *ImageView) OnEvent(app *App, e event.Event) {
}

func (i *ImageView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	if i.visibleRect.Contains(x, y) {
		i.app.SetCell(i.layer, x, y, ch, fg, bg)
	}
}
