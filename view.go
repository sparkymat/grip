package grip

import "github.com/sparkymat/grip/event"

type View interface {
	Initialize(emitEvent func(eventType event.Type, data interface{}))
	OnEvent(app *App, e event.Event)
	Resize(x, y, width, height, visibleX, visibleY, visibleWidth, visibleHeight int)
	Draw()
}
