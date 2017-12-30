package grip

import "github.com/sparkymat/grip/event"

type View interface {
	Initialize(emitEvent func(eventType event.Type, data interface{}))
	OnEvent(app *App, e event.Event)
	OnLoad()
	Resize(x, y, width, height uint32)
	Draw()
}
