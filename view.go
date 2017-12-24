package grip

import "github.com/sparkymat/grip/event"

type View interface {
	Draw()
	OnEvent(e event.Event)
	OnLoad()
	RegisteredEvents() []event.Type
	Resize(x, y, width, height uint32)
	SetApp(app *App)
}
