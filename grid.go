package grip

import (
	"github.com/kr/pretty"
	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/size"
)

type grid struct {
	columnSizes  []size.Size
	rowSizes     []size.Size
	columnWidths []uint32
	rowHeights   []uint32
	views        []View
	windowX      uint32
	windowY      uint32
	windowWidth  uint32
	windowHeight uint32
}

func NewGrid(columnSizes []size.Size, rowSizes []size.Size) grid {
	width, height := termbox.Size()
	columnWidths := distributeLength(uint32(width), columnSizes)
	rowHeights := distributeLength(uint32(height), rowSizes)
	return grid{columnSizes, rowSizes, columnWidths, rowHeights, []View{}, 0, 0, uint32(width), uint32(height)}
}

func (g *grid) AddView(v View) {
	g.views = append(g.views, v)
}

func (g *grid) Resize(x, y, width, height uint32) {
	g.columnWidths = distributeLength(uint32(width), g.columnSizes)
	g.rowHeights = distributeLength(uint32(height), g.rowSizes)

	for _, view := range g.views {
		area := view.GetArea()
		var xOffset uint32 = 0
		var yOffset uint32 = 0
		var viewWidth uint32 = 0
		var viewHeight uint32 = 0
		var i uint32

		for i = 0; i < area.ColumnStart; i++ {
			xOffset += g.columnWidths[i]
		}

		for i = area.ColumnStart; i <= area.ColumnEnd; i++ {
			viewWidth += g.columnWidths[i]
		}

		for i = 0; i < area.RowStart; i++ {
			yOffset += g.rowHeights[i]
		}

		for i = area.RowStart; i <= area.RowEnd; i++ {
			viewHeight += g.rowHeights[i]
		}

		pretty.Logf("Resize(%d, %d, %d, %d)\n", x+xOffset, y+yOffset, viewWidth, viewHeight)
		view.Resize(x+xOffset, y+yOffset, viewWidth, viewHeight)
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
		for lengthIndex, lengthSize := range sizes {
			if lengthSize.Unit == size.Fraction {
				distributedLengths[lengthIndex] = (remainingLength * lengthSize.Value) / totalFractions
			}
		}
	}

	return distributedLengths
}

func (g grid) Draw() {
	for _, view := range g.views {
		view.Draw()
	}
}
