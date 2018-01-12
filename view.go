package grip

import "github.com/sparkymat/grip/event"

type ViewID string

const WildCardPath ViewID = "*"
const AppRoot ViewID = "app-root"
const ModalRoot ViewID = "modal-root"

type View interface {
	Initialize(setCellFn SetCellFn)
	Draw()
	Resize(rect, visibleRect Rect)
	OnEvent(ev event.Event)
}
