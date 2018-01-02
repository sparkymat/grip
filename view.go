package grip

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type ViewID string

const AppRoot ViewID = "app-root"
const ModalRoot ViewID = "modal-root"

const WildCardPath ViewID = "*"

type View interface {
	Initialize(app *App, layer Layer)
	OnEvent(app *App, e event.Event)
	Resize(rect Rect, visibleRect Rect)
	Draw()
	SetCellIfVisible(int, int, rune, termbox.Attribute, termbox.Attribute)
}
