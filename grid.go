package grip

import (
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
	cells           []cell
	columnWidths    []int
	rowHeights      []int
	x               int
	y               int
	width           int
	height          int
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

func (g *Grid) AddView(v View, a Area) {
	g.cells = append(g.cells, cell{v, a})
}

func (g *Grid) Resize(x, y, width, height int) {
	g.x = x
	g.y = y
	g.width = width
	g.height = height

	g.columnWidths = distributeLength(width, g.ColumnSizes)
	g.rowHeights = distributeLength(height, g.RowSizes)

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

		cell.view.Resize(x+xOffset, y+yOffset, viewWidth, viewHeight)
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
		for j := g.y; j < g.y+g.height; j++ {
			for i := g.x; i < g.x+g.width; i++ {
				termbox.SetCell(i, j, ' ', g.BackgroundColor, g.BackgroundColor)
			}
		}
	}

	for _, cell := range g.cells {
		cell.view.Draw()
	}
}
