package grip

import (
	"os"
	"testing"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/size"
)

func HandleEvents(eventChannel <-chan termbox.Event) {
	event := <-eventChannel

	switch event.Type {
	case termbox.EventKey:
		if event.Key == termbox.KeyEsc {
			termbox.Close()
			os.Exit(0)
		}
	}
}

func TestSanity(t *testing.T) {
	app := New()

	mainGrid := NewGrid(
		[]size.Size{size.Auto, size.WithPoints(30)},
		[]size.Size{size.Auto, size.WithPoints(1)},
	)
	mainGrid.AddView(&TextView{
		Text:            "MainArea",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorGreen,
		area:            area{0, 0, 0, 0},
	})
	mainGrid.AddView(&TextView{
		Text:            "SidebarArea",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorBlue,
		area:            area{1, 1, 0, 1},
	})
	mainGrid.AddView(&TextView{
		Text:            "CommandArea",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorRed,
		area:            area{0, 0, 1, 1},
	})

	app.SetRootNode(&mainGrid)

	eventChannel := make(chan termbox.Event)
	go HandleEvents(eventChannel)

	app.Run(eventChannel)
}
