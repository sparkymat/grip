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
	columnWidths    []uint32
	rowHeights      []uint32
	x               uint32
	y               uint32
	width           uint32
	height          uint32
}

func (g *Grid) OnLoad() {
	for _, cell := range g.cells {
		cell.view.OnLoad()
	}
}

func (g *Grid) OnEvent(app *App, e event.Event) {
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

func (g *Grid) Resize(x, y, width, height uint32) {
	g.x = x
	g.y = y
	g.width = width
	g.height = height

	g.columnWidths = distributeLength(uint32(width), g.ColumnSizes)
	g.rowHeights = distributeLength(uint32(height), g.RowSizes)

	for _, cell := range g.cells {
		var xOffset uint32 = 0
		var yOffset uint32 = 0
		var viewWidth uint32 = 0
		var viewHeight uint32 = 0
		var i uint32

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

func distributeLength(totalLength uint32, sizes []size.Size) []uint32 {
	var distributedLengths = make([]uint32, len(sizes))
	var totalFractions uint32 = 0
	var usedLength uint32 = 0

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
		var totalFractionalLength uint32 = 0
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
				termbox.SetCell(int(i), int(j), ' ', g.BackgroundColor, g.BackgroundColor)
			}
		}
	}

	for _, cell := range g.cells {
		cell.view.Draw()
	}
}
