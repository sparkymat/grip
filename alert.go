package grip

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

func (a *App) Alert(message string, onDismiss func(*App)) {
	body := Grid{
		ColumnSizes:     []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(1)},
		RowSizes:        []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(1)},
		HasBackground:   true,
		BackgroundColor: termbox.ColorWhite,
	}

	body.AddView(&TextView{
		Text:            message,
		ForegroundColor: termbox.ColorBlack,
		BackgroundColor: termbox.ColorWhite,
		TextAlignment:   TextAlignmentCenter,
	}, Area{1, 1, 1, 1})

	footer := Grid{
		ColumnSizes:     []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(1)},
		RowSizes:        []size.Size{size.WithPoints(1), size.Auto, size.WithPoints(1)},
		HasBackground:   true,
		BackgroundColor: termbox.ColorWhite,
	}

	footer.AddView(&TextView{
		Text:            "Press any key...",
		ForegroundColor: termbox.ColorBlack | termbox.AttrBold,
		BackgroundColor: termbox.ColorWhite,
		TextAlignment:   TextAlignmentCenter,
	}, Area{1, 1, 1, 1})

	m := NewModal(
		a,
		size.WithPercent(40),
		size.WithPercent(25),
		&TextView{Text: "Alert", ForegroundColor: termbox.ColorBlack | termbox.AttrBold, BackgroundColor: termbox.ColorWhite, TextAlignment: TextAlignmentCenter},
		&body,
		&footer,
		func(app *App, ev event.Event) {
			if ev.Type == event.SystemKeyPress {
				app.HideModal()
				onDismiss(app)
			}
		},
	)

	a.SetModal(&m)
	a.ShowModal()
}
