package grip

type View interface {
	GetArea() area
	Resize(x, y, width, height uint32)
	Draw()
}
