package grip

import (
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

type modal struct {
	emitEvent func(event.Type, interface{})
	grid      Grid
	width     uint32
	height    uint32
	title     View
	body      View
	footer    View
	onEvent   func(*App, event.Event)
}

func NewModal(app *App, width uint32, height uint32, title View, body View, footer View, onEvent func(*App, event.Event)) modal {
	modalGrid := Grid{
		ColumnSizes: []size.Size{size.Auto},
		RowSizes:    []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(3)},
	}

	modalGrid.AddView(title, Area{0, 0, 0, 0})
	modalGrid.AddView(body, Area{0, 0, 1, 1})
	modalGrid.AddView(footer, Area{0, 0, 2, 2})

	return modal{
		width:   width,
		height:  height,
		grid:    modalGrid,
		title:   title,
		body:    body,
		footer:  footer,
		onEvent: onEvent,
	}
}

func (m *modal) Initialize(emit func(event.Type, interface{})) {
	m.emitEvent = emit
	m.grid.Initialize(emit)
}

func (m *modal) Draw() {
	m.grid.Draw()
}

func (m *modal) Resize(x, y, width, height uint32) {
	m.grid.Resize(x, y, width, height)
}

func (m *modal) OnLoad() {
}

func (m *modal) OnEvent(app *App, ev event.Event) {
	if m.onEvent != nil {
		m.onEvent(app, ev)
	}
}
