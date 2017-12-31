package grip

import (
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

type modal struct {
	emitEvent func(event.Type, interface{})
	grid      Grid
	width     size.Size
	height    size.Size
	title     View
	body      View
	footer    View
	onEvent   func(*App, event.Event)
}

func NewModal(app *App, width size.Size, height size.Size, title View, body View, footer View, onEvent func(*App, event.Event)) modal {
	modalGrid := Grid{
		ColumnSizes: []size.Size{size.Auto},
		RowSizes:    []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(3)},
	}

	modalGrid.AddView("title", title, Area{0, 0, 0, 0})
	modalGrid.AddView("body", body, Area{0, 0, 1, 1})
	modalGrid.AddView("footer", footer, Area{0, 0, 2, 2})

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

func (m *modal) Resize(rect, visibleRect Rect) {
	m.grid.Resize(rect, visibleRect)
}

func (m *modal) OnEvent(app *App, ev event.Event) {
	if m.onEvent != nil {
		m.onEvent(app, ev)
	}
}

func (m *modal) Find(path ...ViewID) (View, error) {
	return m.grid.Find(path...)
}
