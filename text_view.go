package grip

import "github.com/nsf/termbox-go"

type TextView struct {
	BackgroundColor termbox.Attribute
	ForegroundColor termbox.Attribute
	Text            string
	area            area
	height          uint32
	width           uint32
	x               uint32
	y               uint32
}

func (t *TextView) GetArea() area {
	return t.area
}

func (t *TextView) Resize(x, y, width, height uint32) {
	t.x = x
	t.y = y
	t.width = width
	t.height = height
}

func (t *TextView) Draw() {
	for i := t.x; i <= (t.x + t.width - 1); i++ {
		for j := t.y; j <= (t.y + t.height - 1); j++ {
			position := (j-t.y)*t.width + (i - t.x)
			char := ' '
			if position < uint32(len(t.Text)) {
				char = rune(t.Text[position])
			}
			termbox.SetCell(int(i), int(j), char, t.ForegroundColor, t.BackgroundColor)
		}
	}
}
