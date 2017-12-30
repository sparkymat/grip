package grip

import "github.com/sparkymat/grip/event"

type ViewID string

const AppRootID ViewID = "app-root"
const ModalRootID ViewID = "modal-root"

type View interface {
	Initialize(emitEvent func(eventType event.Type, data interface{}))
	OnEvent(app *App, e event.Event)
	Resize(x, y, width, height, visibleX, visibleY, visibleWidth, visibleHeight int)
	Draw()
}
