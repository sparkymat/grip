package grip

import "github.com/sparkymat/grip/event"

type View interface {
	Initialize(registerEvent func(eventType event.Type, handler EventHandler), emitEvent func(eventType event.Type, data interface{}) error)
	OnEvent(app *App, e event.Event)
	OnLoad()
	Resize(x, y, width, height uint32)
	Draw()
}
