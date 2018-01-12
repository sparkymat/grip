package grip

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

type TableView struct {
	setCellFn       SetCellFn
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

func (t *TableView) Initialize(setCellFn SetCellFn) {
	t.setCellFn = setCellFn

	t.grid.Initialize(setCellFn)
}

func (t *TableView) AddView(id ViewID, v View, a Area) {
	t.grid.AddView(id, v, Area{a.ColumnStart*2 + 1, a.ColumnEnd*2 + 1, a.RowStart*2 + 1, a.RowEnd*2 + 1})
}

func (t *TableView) OnEvent(e event.Event) {
}

func (t *TableView) Resize(rect, visibleRect Rect) {
	t.grid.Resize(rect, visibleRect)
}

func (t *TableView) Draw() {
	junctionXs := []int{}
	junctionYs := []int{}

	// Vertical
	i := t.grid.rect.Origin.X
	for idx, width := range t.grid.columnWidths {
		if idx%2 == 0 {
			junctionXs = append(junctionXs, i)
			for j := t.grid.rect.Origin.Y; j < t.grid.rect.Origin.Y+t.grid.rect.Size.Height; j++ {
				t.setCellFn(Point{i, j}, ColoredRune{'│', t.foregroundColor, t.backgroundColor})
			}
		}

		i += width
	}

	// Horizontal
	j := t.grid.rect.Origin.Y
	for idx, height := range t.grid.rowHeights {
		if idx%2 == 0 {
			junctionYs = append(junctionYs, j)
			for i := t.grid.rect.Origin.X; i < t.grid.rect.Origin.X+t.grid.rect.Size.Width; i++ {
				t.SetCellIfVisible(i, j, '─', t.foregroundColor, t.backgroundColor)
			}
		}

		j += height
	}

	for _, i := range junctionXs {
		for _, j := range junctionYs {
			if i == t.grid.rect.Origin.X && j == t.grid.rect.Origin.Y {
				t.SetCellIfVisible(i, j, '┌', t.foregroundColor, t.backgroundColor)
			} else if i == t.grid.rect.Origin.X && j == t.grid.rect.Origin.Y+t.grid.rect.Size.Height-1 {
				t.SetCellIfVisible(i, j, '└', t.foregroundColor, t.backgroundColor)
			} else if i == t.grid.rect.Origin.X+t.grid.rect.Size.Width-1 && j == t.grid.rect.Origin.Y {
				t.SetCellIfVisible(i, j, '┐', t.foregroundColor, t.backgroundColor)
			} else if i == t.grid.rect.Origin.X+t.grid.rect.Size.Width-1 && j == t.grid.rect.Origin.Y+t.grid.rect.Size.Height-1 {
				t.SetCellIfVisible(i, j, '┘', t.foregroundColor, t.backgroundColor)
			} else if i == t.grid.rect.Origin.X {
				t.SetCellIfVisible(i, j, '├', t.foregroundColor, t.backgroundColor)
			} else if i == t.grid.rect.Origin.X+t.grid.rect.Size.Width-1 {
				t.SetCellIfVisible(i, j, '┤', t.foregroundColor, t.backgroundColor)
			} else if j == t.grid.rect.Origin.Y {
				t.SetCellIfVisible(i, j, '┬', t.foregroundColor, t.backgroundColor)
			} else if j == t.grid.rect.Origin.Y+t.grid.rect.Size.Height-1 {
				t.SetCellIfVisible(i, j, '┴', t.foregroundColor, t.backgroundColor)
			} else {
				t.SetCellIfVisible(i, j, '┼', t.foregroundColor, t.backgroundColor)
			}
		}
	}

	t.grid.Draw()
}

func (t *TableView) Find(path ...ViewID) (View, error) {
	return t.grid.Find(path...)
}

func (t *TableView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	if t.grid.visibleRect.Contains(x, y) {
		t.setCellFn(Point{x, y}, ColoredRune{ch, fg, bg})
	}
}
