package grip

import (
	"errors"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type app struct {
	rootNode       *grid
	eventListeners map[event.Type][]event.EventHandler
}

func New() app {
	eventListeners := make(map[event.Type][]event.EventHandler)
	return app{
		eventListeners: eventListeners,
	}
}

func (a *app) RegisterEvents(eventTypes ...event.Type) {
	for _, eventType := range eventTypes {
		a.eventListeners[eventType] = []event.EventHandler{}
	}
}

func (a *app) RegisterEventListener(eventType event.Type, handler event.EventHandler) error {
	if _, ok := a.eventListeners[eventType]; !ok {
		return errors.New("Unregistered event")
	}

	listeners := append(a.eventListeners[eventType], handler)
	a.eventListeners[eventType] = listeners

	return nil
}

func (a *app) BroadcastEvent(eventType event.Type, data interface{}) error {
	if _, ok := a.eventListeners[eventType]; !ok {
		return errors.New("Unregistered event")
	}

	for _, registeredView := range a.eventListeners[eventType] {
		registeredView.OnEvent(event.Event{eventType, data})
	}

	return nil
}

func (a *app) SetRootNode(node *grid) {
	a.rootNode = node
	a.rootNode.SetApp(a)
	for _, eventType := range a.rootNode.RegisteredEvents() {
		a.RegisterEventListener(eventType, a.rootNode)
	}
}

func (a app) Run() error {
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
	a.rootNode.Resize(0, 0, uint32(width), uint32(height))
	a.rootNode.Draw()
	termbox.Flush()

	// FIXME: FLush every 16 ms ?
	refreshTicker := time.NewTicker(time.Millisecond * 16)
	go func() {
		for t := range refreshTicker.C {
			a.BroadcastEvent(event.SystemTick, t)
			termbox.Flush()
		}
	}()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			width, height := termbox.Size()
			a.rootNode.Resize(0, 0, uint32(width), uint32(height))
			a.rootNode.Draw()
			termbox.Flush()
			break
		case termbox.EventKey:
			a.BroadcastEvent(event.SystemKeyPress, ev)
			break
		case termbox.EventError:
			panic(ev.Err)
		}
	}

	refreshTicker.Stop()

	return err
}
