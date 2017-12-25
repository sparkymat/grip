package grip

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

type TableView struct {
	grid            Grid
	foregroundColor termbox.Attribute
	backgroundColor termbox.Attribute
}

func NewTableView(columnSizes []size.Size, rowSizes []size.Size, foregroundColor termbox.Attribute, backgroundColor termbox.Attribute) TableView {
	gridColumnSizes := []size.Size{size.WithPoints(1)}
	for _, originalSize := range columnSizes {
		gridColumnSizes = append(gridColumnSizes, originalSize)
		gridColumnSizes = append(gridColumnSizes, size.WithPoints(1))
	}

	gridRowSizes := []size.Size{size.WithPoints(1)}
	for _, originalSize := range rowSizes {
		gridRowSizes = append(gridRowSizes, originalSize)
		gridRowSizes = append(gridRowSizes, size.WithPoints(1))
	}

	return TableView{
		grid: Grid{
			ColumnSizes: gridColumnSizes,
			RowSizes:    gridRowSizes,
		},
		foregroundColor: foregroundColor,
		backgroundColor: backgroundColor,
	}
}

func (t *TableView) AddView(v View, a Area) {
	t.grid.AddView(v, Area{a.ColumnStart*2 + 1, a.ColumnEnd*2 + 1, a.RowStart*2 + 1, a.RowEnd*2 + 1})
}

func (t *TableView) SetApp(app *App) {
	t.grid.SetApp(app)
}

func (t *TableView) RegisteredEvents() []event.Type {
	return []event.Type{}
}

func (t *TableView) OnLoad() {
}

func (t *TableView) OnEvent(e event.Event) {
}

func (t *TableView) Resize(x, y, width, height uint32) {
	t.grid.Resize(x, y, width, height)
}

func (t *TableView) Draw() {
	junctionXs := []uint32{}
	junctionYs := []uint32{}

	// Vertical
	i := t.grid.x
	for idx, width := range t.grid.columnWidths {
		if idx%2 == 0 {
			junctionXs = append(junctionXs, i)
			for j := t.grid.y; j < t.grid.y+t.grid.height; j++ {
				termbox.SetCell(int(i), int(j), '│', t.foregroundColor, t.backgroundColor)
			}
		}

		i += width
	}

	// Horizontal
	j := t.grid.y
	for idx, height := range t.grid.rowHeights {
		if idx%2 == 0 {
			junctionYs = append(junctionYs, j)
			for i := t.grid.x; i < t.grid.x+t.grid.width; i++ {
				termbox.SetCell(int(i), int(j), '─', t.foregroundColor, t.backgroundColor)
			}
		}

		j += height
	}

	for _, i := range junctionXs {
		for _, j := range junctionYs {
			if i == t.grid.x && j == t.grid.y {
				termbox.SetCell(int(i), int(j), '┌', t.foregroundColor, t.backgroundColor)
			} else if i == t.grid.x && j == t.grid.y+t.grid.height-1 {
				termbox.SetCell(int(i), int(j), '└', t.foregroundColor, t.backgroundColor)
			} else if i == t.grid.x+t.grid.width-1 && j == t.grid.y {
				termbox.SetCell(int(i), int(j), '┐', t.foregroundColor, t.backgroundColor)
			} else if i == t.grid.x+t.grid.width-1 && j == t.grid.y+t.grid.height-1 {
				termbox.SetCell(int(i), int(j), '┘', t.foregroundColor, t.backgroundColor)
			} else if i == t.grid.x {
				termbox.SetCell(int(i), int(j), '├', t.foregroundColor, t.backgroundColor)
			} else if i == t.grid.x+t.grid.width-1 {
				termbox.SetCell(int(i), int(j), '┤', t.foregroundColor, t.backgroundColor)
			} else if j == t.grid.y {
				termbox.SetCell(int(i), int(j), '┬', t.foregroundColor, t.backgroundColor)
			} else if j == t.grid.y+t.grid.height-1 {
				termbox.SetCell(int(i), int(j), '┴', t.foregroundColor, t.backgroundColor)
			} else {
				termbox.SetCell(int(i), int(j), '┼', t.foregroundColor, t.backgroundColor)
			}
		}
	}

	t.grid.Draw()
}
