package grip

type View interface {
	Resize(x, y, width, height uint32)
	Draw()
}
