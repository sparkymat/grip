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
	setCellFn      SetCellFn
	View           View
	Direction      ScrollDirection
	Size           int
	scrollPosition int
	rect           Rect
	visibleRect    Rect
}

func (s *ScrollView) Initialize(setCellFn SetCellFn) {
	s.setCellFn = setCellFn
	s.scrollPosition = 0
	s.View.Initialize(setCellFn)
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

func (s *ScrollView) OnEvent(e event.Event) {
	s.View.OnEvent(e)
}

func (s *ScrollView) Resize(rect, visibleRect Rect) {
	s.rect = rect
	s.visibleRect = visibleRect

	switch s.Direction {
	case ScrollDirectionHorizontal:
		if s.scrollPosition+s.rect.Size.Width >= s.Size {
			s.scrollPosition -= (s.Size - s.rect.Size.Width)
		}
		s.View.Resize(
			Rect{
				Point{s.rect.Origin.X - s.scrollPosition, s.rect.Origin.Y},
				Size{s.Size, s.rect.Size.Height},
			},
			s.rect,
		)
		break
	case ScrollDirectionVertical:
		if s.scrollPosition > s.Size-s.rect.Size.Height {
			s.scrollPosition -= (s.Size - s.rect.Size.Height)
		}
		s.View.Resize(
			Rect{
				Point{s.rect.Origin.X, s.rect.Origin.Y - s.scrollPosition},
				Size{s.rect.Size.Width, s.Size},
			},
			s.rect,
		)
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
		s.setCellFn(Point{x, y}, ColoredRune{ch, fg, bg})
	}
}
