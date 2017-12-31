package grip

import (
	"errors"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

type Grid struct {
	emitEvent       func(event.Type, interface{})
	HasBackground   bool
	BackgroundColor termbox.Attribute
	ColumnSizes     []size.Size
	RowSizes        []size.Size
	app             *App
	cells           map[ViewID]cell
	columnWidths    []int
	rowHeights      []int
	rect            Rect
	visibleRect     Rect
}

func (g *Grid) OnEvent(app *App, e event.Event) {
	for _, cell := range g.cells {
		cell.view.OnEvent(app, e)
	}
}

func (g *Grid) Initialize(emit func(eventType event.Type, data interface{})) {
	g.emitEvent = emit
	for _, cell := range g.cells {
		cell.view.Initialize(emit)
	}
}

func (g *Grid) AddView(id ViewID, v View, a Area) error {
	if g.cells == nil {
		g.cells = make(map[ViewID]cell)
	}

	if _, ok := g.cells[id]; ok {
		panic("This view ID is already registered")
	}

	g.cells[id] = cell{v, a}

	return nil
}

func (g *Grid) Resize(rect, visibleRect Rect) {
	g.rect = rect
	g.visibleRect = visibleRect

	g.columnWidths = distributeLength(rect.Width, g.ColumnSizes)
	g.rowHeights = distributeLength(rect.Height, g.RowSizes)

	for _, cell := range g.cells {
		var xOffset int = 0
		var yOffset int = 0
		var viewWidth int = 0
		var viewHeight int = 0
		var i int

		for i = 0; i < cell.area.ColumnStart; i++ {
			xOffset += g.columnWidths[i]
		}

		for i = cell.area.ColumnStart; i <= cell.area.ColumnEnd; i++ {
			viewWidth += g.columnWidths[i]
		}

		for i = 0; i < cell.area.RowStart; i++ {
			yOffset += g.rowHeights[i]
		}

		for i = cell.area.RowStart; i <= cell.area.RowEnd; i++ {
			viewHeight += g.rowHeights[i]
		}

		cell.view.Resize(
			Rect{
				X:      rect.X + xOffset,
				Y:      rect.Y + yOffset,
				Width:  viewWidth,
				Height: viewHeight,
			},
			Rect{
				X:      rect.X + xOffset,
				Y:      rect.Y + yOffset,
				Width:  viewWidth,
				Height: viewHeight,
			})
	}
}

func distributeLength(totalLength int, sizes []size.Size) []int {
	var distributedLengths = make([]int, len(sizes))
	var totalFractions int = 0
	var usedLength int = 0

	for lengthIndex, lengthSize := range sizes {
		if lengthSize.Unit == size.Fraction {
			totalFractions += lengthSize.Value
			distributedLengths[lengthIndex] = 0
		} else if lengthSize.Unit == size.Point {
			distributedLengths[lengthIndex] = lengthSize.Value
		} else if lengthSize.Unit == size.Percent {
			distributedLengths[lengthIndex] = (totalLength * lengthSize.Value) / 100
		}
		usedLength += distributedLengths[lengthIndex]
	}
	remainingLength := totalLength - usedLength

	if usedLength < totalLength {
		var totalFractionalLength int = 0
		var lastFractionIndex int = -1
		for lengthIndex, lengthSize := range sizes {
			if lengthSize.Unit == size.Fraction {
				lastFractionIndex = lengthIndex
				distributedLengths[lengthIndex] = (remainingLength * lengthSize.Value) / totalFractions
				totalFractionalLength += (remainingLength * lengthSize.Value) / totalFractions
			}
		}
		if usedLength+totalFractionalLength < totalLength && lastFractionIndex != -1 {
			distributedLengths[lastFractionIndex] += (totalLength - (usedLength + totalFractionalLength))
		}
	}

	return distributedLengths
}

func (g *Grid) Draw() {
	if g.HasBackground {
		for j := g.rect.Y; j < g.rect.Y+g.rect.Height; j++ {
			for i := g.rect.X; i < g.rect.X+g.rect.Width; i++ {
				g.SetCellIfVisible(i, j, ' ', g.BackgroundColor, g.BackgroundColor)
			}
		}
	}

	for _, cell := range g.cells {
		cell.view.Draw()
	}
}

func (g *Grid) Find(path ...ViewID) (View, error) {
	if path == nil || len(path) == 0 {
		return nil, errors.New("View not found")
	}
	currentID := path[0]
	remainingPath := path[1:]

	if currentID == WildCardPath {
		matchedView, err := g.Find(remainingPath...)
		if err == nil {
			return matchedView, nil
		}

		for _, cell := range g.cells {
			if matchedViewContainer, isViewContainer := cell.view.(ViewContainer); isViewContainer {
				matchedView, err := matchedViewContainer.Find(path...)
				if err == nil {
					return matchedView, nil
				}
			}
		}

		return nil, errors.New("View not found")
	} else {
		if matchedCell, viewFound := g.cells[currentID]; viewFound {
			if len(remainingPath) == 0 {
				return matchedCell.view, nil
			} else if matchedViewContainer, isViewContainer := matchedCell.view.(ViewContainer); isViewContainer {
				return matchedViewContainer.Find(remainingPath...)
			} else {
				return nil, errors.New("View not found")
			}
		} else {
			return nil, errors.New("View not found")
		}
	}

	return nil, nil
}

func (g *Grid) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	if g.visibleRect.Contains(x, y) {
		termbox.SetCell(x, y, ch, fg, bg)
	}
}
