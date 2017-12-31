package grip

import "github.com/sparkymat/grip/event"

type ViewID string

const AppRoot ViewID = "app-root"
const ModalRoot ViewID = "modal-root"

const WildCardPath ViewID = "*"

type View interface {
	Initialize(emitEvent func(eventType event.Type, data interface{}))
	OnEvent(app *App, e event.Event)
	Resize(x, y, width, height, visibleX, visibleY, visibleWidth, visibleHeight int)
	Draw()
}
