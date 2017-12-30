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
	height          int
	width           int
	x               int
	y               int
}

func (i *ImageView) Initialize(emit func(eventType event.Type, data interface{})) {
	i.emitEvent = emit
}

func (i *ImageView) Resize(x, y, width, height int) {
	i.x = x
	i.y = y
	i.width = width
	i.height = height

	i.scaleAscii = asciiart.Convert2AsciiOfWidth(i.Image, int(i.width)-1)
}

func (v *ImageView) Draw() {
	for j := v.y; j <= (v.y + v.height - 1); j++ {
		for i := v.x + 1; i < (v.x + v.width - 1); i++ {
			var r rune = ' '
			if j*v.width+i < len(v.scaleAscii) {
				r = rune(v.scaleAscii[j*v.width+i])
			}
			termbox.SetCell(i, j, r, v.ForegroundColor, v.BackgroundColor)
		}
	}
}

func (i *ImageView) OnEvent(app *App, e event.Event) {
}
