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
		[]size.Size{size.Auto, size.WithPercent(30)},
		[]size.Size{size.Auto, size.WithPoints(1)},
	)

	sidebarGrid := NewGrid(
		[]size.Size{size.Auto},
		[]size.Size{size.WithPoints(1), size.WithPoints(1), size.Auto},
	)

	sidebarGrid.AddView(&TextView{
		Text:            "Name: Adam",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorCyan,
	}, Area{0, 0, 0, 0})

	sidebarGrid.AddView(&TextView{
		Text:            "Class: Warlock",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorYellow,
	}, Area{0, 0, 1, 1})

	sidebarGrid.AddView(&TextView{
		Text:            "SidebarArea - Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eu consectetur lacus. Sed tincidunt eros non ultrices commodo. Sed ornare id dolor sed ultricies. Duis in est at nulla pretium mattis ac quis quam. Maecenas nibh nisi, rhoncus quis iaculis sit amet, semper et diam. Aenean pharetra ex non mi placerat rhoncus. Vivamus erat ante, suscipit vitae aliquet id, congue et dolor. Curabitur sed tortor tortor. Duis non sem et lacus ultrices finibus quis quis felis. Integer non elementum ante. Vestibulum vel augue ut tortor condimentum pulvinar eu blandit leo. Donec nibh nibh, tincidunt vitae risus a, consectetur suscipit felis. Quisque elementum velit nec mauris tristique, id malesuada tellus dictum.",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorBlue,
	}, Area{0, 0, 2, 2})

	mainGrid.AddView(&TextView{
		Text:            "MainArea",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorGreen,
	}, Area{0, 0, 0, 0})

	mainGrid.AddView(&sidebarGrid, Area{1, 1, 0, 1})

	mainGrid.AddView(&TextView{
		Text:            "CommandArea",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorRed,
	}, Area{0, 0, 1, 1})

	app.SetRootNode(&mainGrid)

	eventChannel := make(chan termbox.Event)
	go HandleEvents(eventChannel)

	app.Run(eventChannel)
}
