package grip

import (
	"errors"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type app struct {
	rootNode       *grid
	eventListeners map[event.Type][]View
}

func New() app {
	eventListeners := make(map[event.Type][]View)
	return app{
		eventListeners: eventListeners,
	}
}

func (a *app) RegisterEvents(eventTypes ...event.Type) {
	for _, eventType := range eventTypes {
		a.eventListeners[eventType] = []View{}
	}
}

func (a *app) registerEventListener(eventType event.Type, view View) error {
	if _, ok := a.eventListeners[eventType]; !ok {
		return errors.New("Unregistered event")
	}

	listeners := append(a.eventListeners[eventType], view)
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
		a.registerEventListener(eventType, a.rootNode)
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

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			width, height := termbox.Size()
			a.rootNode.Resize(0, 0, uint32(width), uint32(height))
			a.rootNode.Draw()
			termbox.Flush()
		case termbox.EventKey:
			a.BroadcastEvent(event.GlobalKeyPress, ev)
		case termbox.EventError:
			panic(ev.Err)
		}
	}

	return err
}
