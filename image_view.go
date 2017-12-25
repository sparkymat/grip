package grip

import (
	"image"

	termbox "github.com/nsf/termbox-go"
	asciiart "github.com/sparkymat/goasciiart"
	"github.com/sparkymat/grip/event"
)

type ImageView struct {
	app             *App
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	Image           image.Image
	scaleAscii      []byte
	height          uint32
	width           uint32
	x               uint32
	y               uint32
}

func (i *ImageView) SetApp(app *App) {
	i.app = app
}

func (i *ImageView) RegisteredEvents() []event.Type {
	return []event.Type{}
}

func (i *ImageView) Resize(x, y, width, height uint32) {
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
			if j*v.width+i < uint32(len(v.scaleAscii)) {
				r = rune(v.scaleAscii[j*v.width+i])
			}
			termbox.SetCell(int(i), int(j), r, v.ForegroundColor, v.BackgroundColor)
		}
	}
}

func (i *ImageView) OnEvent(e event.Event) {
}

func (i *ImageView) OnLoad() {
}
