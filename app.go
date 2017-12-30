package grip

import (
	"errors"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

type EventHandler interface {
	OnEvent(app *App, e event.Event)
}

type App struct {
	container            View
	globalEventListeners map[event.Type][]EventHandler
	modalContainer       View
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

	for _, registeredView := range a.globalEventListeners[eventType] {
		go registeredView.OnEvent(a, event.Event{eventType, data})
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
	a.container.Initialize(a.EmitEvent)
}

func (a *App) SetModal(m *modal) {
	modalGrid := Grid{
		ColumnSizes: []size.Size{size.Auto, size.WithPoints(m.width), size.Auto},
		RowSizes:    []size.Size{size.Auto, size.WithPoints(m.height), size.Auto},
	}

	modalGrid.AddView(m, Area{1, 1, 1, 1})

	a.modalContainer = &modalGrid
	a.modalContainer.Initialize(a.EmitModalEvent)
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
	width, height := termbox.Size()

	a.container.Resize(0, 0, uint32(width), uint32(height))
	a.container.OnLoad()
	a.container.Draw()

	if a.modalVisible {
		a.modalContainer.OnLoad()
		a.modalContainer.Resize(0, 0, uint32(width), uint32(height))
		a.modalContainer.Draw()
	}

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
			width, height := termbox.Size()

			a.container.Resize(0, 0, uint32(width), uint32(height))
			a.container.Draw()

			if a.modalVisible && a.modalContainer != nil {
				a.modalContainer.Resize(0, 0, uint32(width), uint32(height))
				a.modalContainer.Draw()
			}

			termbox.Flush()
			break
		case termbox.EventKey:
			a.EmitGlobalEvent(event.SystemKeyPress, ev)

			if a.modalVisible {
				a.EmitModalEvent(event.SystemKeyPress, ev)
			} else {
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

func (a *App) ShowModal() error {
	if a.modalContainer == nil {
		return errors.New("No modal to show")
	}

	a.modalVisible = true
	width, height := termbox.Size()
	a.modalContainer.Resize(0, 0, uint32(width), uint32(height))
	a.modalContainer.OnLoad()
	a.modalContainer.Draw()

	return nil
}

func (a *App) HideModal() {
	a.modalVisible = false
	a.container.Draw()
}
