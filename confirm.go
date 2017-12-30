package grip

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

func (a *App) Confirm(message string, onConfirm func(*App), onDismiss func(*App)) {
	body := Grid{
		ColumnSizes:     []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(1)},
		RowSizes:        []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(1)},
		HasBackground:   true,
		BackgroundColor: termbox.ColorWhite,
	}

	body.AddView("body-text", &TextView{
		Text:            message,
		ForegroundColor: termbox.ColorBlack,
		BackgroundColor: termbox.ColorWhite,
		TextAlignment:   TextAlignmentCenter,
	}, Area{1, 1, 1, 1})

	footer := Grid{
		ColumnSizes: []size.Size{size.Auto, size.Auto},
		RowSizes:    []size.Size{size.Auto},
	}

	yesBox := Grid{
		ColumnSizes:     []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(1)},
		RowSizes:        []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(1)},
		HasBackground:   true,
		BackgroundColor: termbox.ColorGreen,
	}

	yesBox.AddView("yes-text", &TextView{
		Text:            "[Y]es",
		ForegroundColor: termbox.ColorBlack | termbox.AttrBold,
		BackgroundColor: termbox.ColorGreen,
		TextAlignment:   TextAlignmentCenter,
	}, Area{1, 1, 1, 1})

	noBox := Grid{
		ColumnSizes:     []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(1)},
		RowSizes:        []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(1)},
		HasBackground:   true,
		BackgroundColor: termbox.ColorRed,
	}
	noBox.AddView("no-text", &TextView{
		Text:            "[N]o",
		ForegroundColor: termbox.ColorBlack | termbox.AttrBold,
		BackgroundColor: termbox.ColorRed,
		TextAlignment:   TextAlignmentCenter,
	}, Area{1, 1, 1, 1})

	footer.AddView("no-box", &noBox, Area{0, 0, 0, 0})
	footer.AddView("yes-box", &yesBox, Area{1, 1, 0, 0})

	m := NewModal(
		a,
		size.WithPercent(40),
		size.WithPercent(25),
		&TextView{Text: "Confirm", ForegroundColor: termbox.ColorBlack | termbox.AttrBold, BackgroundColor: termbox.ColorWhite, TextAlignment: TextAlignmentCenter},
		&body,
		&footer,
		func(app *App, ev event.Event) {
			if ev.Type == event.SystemKeyPress {
				termboxEvent := ev.Data.(termbox.Event)

				if termboxEvent.Key == termbox.KeyEsc || termboxEvent.Ch == 'n' || termboxEvent.Ch == 'N' {
					app.HideModal()
					onDismiss(app)
				} else if termboxEvent.Ch == 'y' || termboxEvent.Ch == 'Y' {
					app.HideModal()
					onConfirm(app)
				}
			}
		},
	)

	a.SetModal(&m)
	a.ShowModal()
}
