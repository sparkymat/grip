package grip

import "github.com/sparkymat/grip/event"

type View interface {
	Resize(x, y, width, height uint32)
	Draw()
	SetApp(app *app)
	RegisteredEvents() []event.Type
	OnEvent(e event.Event)
}
