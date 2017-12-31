package grip

import (
	"errors"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
)

type ScrollDirection int

const ScrollDirectionHorizontal ScrollDirection = 0
const ScrollDirectionVertical ScrollDirection = 1

type ScrollView struct {
	View           View
	Direction      ScrollDirection
	Size           int
	scrollPosition int
	emitEvent      func(event.Type, interface{})
	rect           Rect
	visibleRect    Rect
}

func (s *ScrollView) Initialize(emit func(eventType event.Type, data interface{})) {
	s.emitEvent = emit
	s.scrollPosition = 0
}

func (s *ScrollView) GetScrollPosition() int {
	return s.scrollPosition
}

func (s *ScrollView) ScrollTo(scrollPosition int) {
	if scrollPosition < 0 {
		s.scrollPosition = 0
	} else if scrollPosition >= s.Size {
		scrollPosition = s.Size - 1
	} else {
		s.scrollPosition = scrollPosition
	}
}

func (s *ScrollView) OnEvent(app *App, e event.Event) {
	s.View.OnEvent(app, e)
}

func (s *ScrollView) Resize(rect, visibleRect Rect) {
	s.rect = rect
	s.visibleRect = visibleRect

	switch s.Direction {
	case ScrollDirectionHorizontal:
		s.View.Resize(Rect{X: s.rect.X - s.scrollPosition, Y: s.rect.Y, Width: s.Size, Height: s.rect.Height}, s.rect)
		break
	case ScrollDirectionVertical:
		s.View.Resize(Rect{X: s.rect.X, Y: s.rect.Y - s.scrollPosition, Width: s.rect.Width, Height: s.Size}, s.rect)
		break
	}
}

func (s *ScrollView) Draw() {
	s.Resize(s.rect, s.visibleRect)
	s.View.Draw()
}

func (s *ScrollView) Find(path ...ViewID) (View, error) {
	viewContainer, isViewContainer := s.View.(ViewContainer)

	if !isViewContainer {
		return nil, errors.New("View not found")
	} else {
		return viewContainer.Find(path...)
	}
}

func (s *ScrollView) SetCellIfVisible(x int, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
	if s.visibleRect.Contains(x, y) {
		termbox.SetCell(x, y, ch, fg, bg)
	}
}
