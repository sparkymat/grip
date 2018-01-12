package grip

type Point struct {
	X int
	Y int
}

type Size struct {
	Width  int
	Height int
}

type Rect struct {
	Origin Point
	Size   Size
}

func (r Rect) Contains(x, y int) bool {
	return x >= r.Origin.X && x < r.Origin.X+r.Size.Width && y >= r.Origin.Y && y < r.Origin.Y+r.Size.Height
}
