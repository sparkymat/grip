package grip

import (
	"errors"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

type EventHandler func(*App, event.Event)
type DrawCellFn func(int, int, rune, termbox.Attribute, termbox.Attribute)
type EmitEventFn func(event.Type, interface{})

type App struct {
	container            ViewContainer
	globalEventListeners map[event.Type][]EventHandler
	modalRect            Rect
	modalContainer       ViewContainer
	modalVisible         bool
}

func (a *App) RegisterGlobalEventListener(eventType event.Type, handler EventHandler) {
	if a.globalEventListeners == nil {
		a.globalEventListeners = make(map[event.Type][]EventHandler)
	}

	if _, ok := a.globalEventListeners[eventType]; !ok {
		a.globalEventListeners[eventType] = []EventHandler{}
	}

	listeners := append(a.globalEventListeners[eventType], handler)
	a.globalEventListeners[eventType] = listeners
}

func (a *App) EmitGlobalEvent(eventType event.Type, data interface{}) error {
	if _, ok := a.globalEventListeners[eventType]; !ok {
		return errors.New("Unregistered event")
	}

	for _, registeredHandler := range a.globalEventListeners[eventType] {
		go registeredHandler(a, event.Event{eventType, data})
	}

	return nil
}

func (a *App) EmitEvent(eventType event.Type, data interface{}) {
	if a.container != nil {
		go a.container.OnEvent(a, event.Event{Type: eventType, Data: data})
	}
}

func (a *App) EmitModalEvent(eventType event.Type, data interface{}) {
	if a.modalContainer != nil {
		go a.modalContainer.OnEvent(a, event.Event{Type: eventType, Data: data})
	}
}

func (a *App) SetContainer(container *Grid) {
	a.container = container
	a.container.Initialize(a.MainSetCell, a.EmitEvent)
}

func (a *App) SetModal(m *modal) {
	modalGrid := Grid{
		ColumnSizes: []size.Size{size.Auto, m.width, size.Auto},
		RowSizes:    []size.Size{size.Auto, m.height, size.Auto},
	}

	modalGrid.AddView("modal", m, Area{1, 1, 1, 1})

	a.modalContainer = &modalGrid
	a.modalContainer.Initialize(a.ModalSetCell, a.EmitModalEvent)
}

func (a App) Run() error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetOutputMode(termbox.Output256)

	defer termbox.Close()

	// Draw initial
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	a.Refresh()
	termbox.Flush()

	// FIXME: FLush every 16 ms ?
	refreshTicker := time.NewTicker(time.Millisecond * 16)
	go func() {
		for t := range refreshTicker.C {
			a.EmitEvent(event.SystemTick, t)
			a.EmitModalEvent(event.SystemTick, t)
			a.EmitGlobalEvent(event.SystemTick, t)
			termbox.Flush()
		}
	}()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			a.Refresh()
			termbox.Flush()
			break
		case termbox.EventKey:

			if a.modalVisible {
				a.EmitModalEvent(event.SystemKeyPress, ev)
			} else {
				a.EmitGlobalEvent(event.SystemKeyPress, ev)
				a.EmitEvent(event.SystemKeyPress, ev)
			}
			break
		case termbox.EventError:
			panic(ev.Err)
		}
	}

	refreshTicker.Stop()

	return err
}

func (a *App) Refresh() {
	width, height := termbox.Size()

	if a.container != nil {
		a.container.Resize(
			Rect{
				X:      0,
				Y:      0,
				Width:  width,
				Height: height,
			},
			Rect{
				X:      0,
				Y:      0,
				Width:  width,
				Height: height,
			},
		)
		a.container.Draw()
	}

	a.modalRect = Rect{0, 0, 0, 0}
	if a.modalVisible && a.modalContainer != nil {
		a.modalContainer.Resize(
			Rect{
				X:      0,
				Y:      0,
				Width:  width,
				Height: height,
			},
			Rect{
				X:      0,
				Y:      0,
				Width:  width,
				Height: height,
			},
		)
		a.modalRect = Rect{
			X:      a.modalContainer.(*Grid).columnWidths[0],
			Y:      a.modalContainer.(*Grid).rowHeights[0],
			Width:  a.modalContainer.(*Grid).columnWidths[1],
			Height: a.modalContainer.(*Grid).rowHeights[1],
		}
		a.modalContainer.Draw()
	}
}

func (a *App) ShowModal() error {
	if a.modalContainer == nil {
		return errors.New("No modal to show")
	}

	a.modalVisible = true
	a.Refresh()

	return nil
}

func (a *App) HideModal() {
	a.modalVisible = false
	a.Refresh()
}

func (a *App) Find(path ...ViewID) (View, error) {
	if path == nil || len(path) == 0 {
		return nil, errors.New("View not found")
	}

	currentID := path[0]
	remainingPath := path[1:]

	if currentID == WildCardPath {
		if a.modalVisible {
			if a.modalContainer != nil {
				return a.modalContainer.Find(path...)
			}
		} else {
			if a.container != nil {
				return a.container.Find(path...)
			}
		}
	} else if currentID == AppRoot && a.container != nil {
		return a.container.Find(remainingPath...)
	} else if currentID == ModalRoot && a.modalVisible && a.modalContainer != nil {
		return a.modalContainer.Find(remainingPath...)
	}

	return nil, errors.New("View not found")
}

func (a *App) MainSetCell(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	if !a.modalVisible || a.modalContainer == nil || !a.modalRect.Contains(x, y) {
		termbox.SetCell(x, y, ch, fg, bg)
	}
}

func (a *App) ModalSetCell(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	termbox.SetCell(x, y, ch, fg, bg)
}
