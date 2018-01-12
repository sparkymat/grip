package grip

import (
	"errors"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

type SetCellFn func(Point, ColoredRune)
type EventHandler func(app *App, ev event.Event)

type App struct {
	View          ViewContainer
	eventChannel  chan event.Event
	modalView     ViewContainer
	showModal     bool
	modalRect     Rect
	eventHandlers map[event.Type][]EventHandler
}

func (a *App) Run() {
	err := termbox.Init()
	if err != nil {
		panic(err.Error())
	}
	termbox.SetOutputMode(termbox.Output256)

	defer termbox.Close()

	a.eventChannel = make(chan (event.Event), 20) // For nested calls which write to the channel
	go a.eventLoop()

	ticker := time.NewTicker(time.Millisecond * 40) // for 25fps
	go func() {
		for t := range ticker.C {
			a.DispatchEvent(event.EventTick, t)
		}
	}()

	// This will push to the channel (draw()) which will block if channel is busy
	a.View.Initialize(a.SetCellApp)
	a.Refresh()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			width, height := termbox.Size()
			a.DispatchEvent(event.EventResize, Size{width, height})
		case termbox.EventKey:
			a.DispatchEvent(event.EventKeyPress, ev)
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func (a *App) Refresh() {
	a.eventChannel <- event.Event{event.EventRefresh, nil}
}

func (a *App) eventLoop() {
	width, height := termbox.Size()
	windowSize := Size{width, height}
	windowPosition := Point{0, 0}

	for ev := range a.eventChannel {
		switch ev.Type {
		case event.EventRefresh:
			a.View.Resize(Rect{windowPosition, windowSize}, Rect{windowPosition, windowSize})
			a.View.Draw()

			if a.showModal && a.modalView != nil {
				a.modalView.Resize(Rect{windowPosition, windowSize}, Rect{windowPosition, windowSize})
				a.modalView.Draw()
			}
		case event.EventResize:
			windowSize = ev.Data.(Size)
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			a.Refresh()
		case event.EventKeyPress:
			a.handleEvent(ev)
		case event.EventTick:
			termbox.Flush()
			a.handleEvent(ev)
		case event.EventDrawCellRequest:
			request := ev.Data.(DrawCellRequest)
			termbox.SetCell(request.Position.X, request.Position.Y, request.Rune.Ch, request.Rune.ForegroundColor, request.Rune.BackgroundColor)
		case event.EventShowModal:
			modal := ev.Data.(*modal)
			modalGrid := Grid{
				ColumnSizes: []size.Size{size.Auto, modal.width, size.Auto},
				RowSizes:    []size.Size{size.Auto, modal.height, size.Auto},
			}
			modalGrid.AddView("modal", modal, Area{1, 1, 1, 1})
			a.modalView = &modalGrid

			a.modalView.Initialize(a.SetCellModal)
			a.showModal = true
			a.modalView.Resize(Rect{windowPosition, windowSize}, Rect{windowPosition, windowSize})

			a.Refresh()
		case event.EventHideModal:
			a.showModal = false
			a.modalView = nil
			a.Refresh()
		default:
			a.handleEvent(ev)
		}
	}
}

func (a *App) handleEvent(ev event.Event) {
	if a.showModal {
		a.modalView.OnEvent(ev)
	} else {
		for _, handler := range a.eventHandlers[ev.Type] {
			handler(a, ev)
		}
		a.View.OnEvent(ev)
	}
}

func (a *App) DispatchEvent(ev event.Type, data interface{}) {
	a.eventChannel <- event.Event{ev, data}
}

func (a *App) RegisterEventListener(ev event.Type, handler EventHandler) {
	if a.eventHandlers == nil {
		a.eventHandlers = make(map[event.Type][]EventHandler)
	}

	if _, hasHandlers := a.eventHandlers[ev]; !hasHandlers {
		a.eventHandlers[ev] = []EventHandler{}
	}

	a.eventHandlers[ev] = append(a.eventHandlers[ev], handler)
}

func (a *App) ShowModal(modal *modal) {
	a.eventChannel <- event.Event{event.EventShowModal, modal}
}

func (a *App) HideModal() {
	a.eventChannel <- event.Event{event.EventHideModal, nil}
}

func (a *App) SetCellModal(position Point, ch ColoredRune) {
	termbox.SetCell(position.X, position.Y, ch.Ch, ch.ForegroundColor, ch.BackgroundColor)
}

func (a *App) SetCellApp(position Point, ch ColoredRune) {
	termbox.SetCell(position.X, position.Y, ch.Ch, ch.ForegroundColor, ch.BackgroundColor)
}

func (a *App) Find(path ...ViewID) (View, error) {
	if path == nil || len(path) == 0 {
		return nil, errors.New("View not found")
	}

	currentID := path[0]
	remainingPath := path[1:]

	if currentID == WildCardPath {
		if a.showModal {
			if a.modalView != nil {
				return a.modalView.Find(path...)
			}
		} else {
			if a.View != nil {
				return a.View.Find(path...)
			}
		}
	} else if currentID == AppRoot && a.View != nil {
		return a.View.Find(remainingPath...)
	} else if currentID == ModalRoot && a.showModal && a.modalView != nil {
		return a.modalView.Find(remainingPath...)
	}

	return nil, errors.New("View not found")
}
