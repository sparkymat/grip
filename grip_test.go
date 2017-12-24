package grip

import (
	"os"
	"testing"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

type TestEventHandler struct{}

func (t TestEventHandler) OnEvent(e event.Event) {
	switch e.Type {
	case event.SystemKeyPress:
		termboxEvent := e.Data.(termbox.Event)
		if termboxEvent.Type == termbox.EventKey && termboxEvent.Key == termbox.KeyEsc {
			termbox.Close()
			os.Exit(0)
		}
		break
	}
}

func TestSanity(t *testing.T) {
	app := New()

	app.RegisterEvents(
		event.SystemKeyPress,
		event.SystemTick,
	)

	mainGrid := NewGrid(
		[]size.Size{size.Auto, size.WithPercent(30)},
		[]size.Size{size.Auto, size.WithPoints(1)},
	)

	sidebarGrid := NewGrid(
		[]size.Size{size.Auto},
		[]size.Size{size.WithPoints(1), size.WithPoints(1), size.WithPoints(1), size.WithPoints(1), size.Auto},
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

	progress := ProgressView{
		Type:            ProgressViewTypePercentage,
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorMagenta,
		MinimumValue:    0,
		MaximumValue:    1000,
		CurrentValue:    666,
	}

	progressTimer := time.NewTimer(time.Second * 2)
	go func() {
		<-progressTimer.C
		progress.CurrentValue += 200
		progress.Draw()
	}()
	sidebarGrid.AddView(&progress, Area{0, 0, 2, 2})

	sidebarGrid.AddView(&ActivityView{
		BackgroundColor: termbox.ColorRed,
		ForegroundColor: termbox.ColorWhite,
		Text:            "Loading...",
	}, Area{0, 0, 3, 3})

	sidebarGrid.AddView(&TextView{
		Text:            "SidebarArea - Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eu consectetur lacus. Sed tincidunt eros non ultrices commodo. Sed ornare id dolor sed ultricies. Duis in est at nulla pretium mattis ac quis quam. Maecenas nibh nisi, rhoncus quis iaculis sit amet, semper et diam. Aenean pharetra ex non mi placerat rhoncus. Vivamus erat ante, suscipit vitae aliquet id, congue et dolor. Curabitur sed tortor tortor. Duis non sem et lacus ultrices finibus quis quis felis. Integer non elementum ante. Vestibulum vel augue ut tortor condimentum pulvinar eu blandit leo. Donec nibh nibh, tincidunt vitae risus a, consectetur suscipit felis. Quisque elementum velit nec mauris tristique, id malesuada tellus dictum.",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorBlue,
	}, Area{0, 0, 4, 4})

	mainGrid.AddView(&TextView{
		Text:            "MainArea",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorGreen,
	}, Area{0, 0, 0, 0})

	mainGrid.AddView(&sidebarGrid, Area{1, 1, 0, 1})

	mainGrid.AddView(&InputView{
		TextView: TextView{
			Text:            "",
			ForegroundColor: termbox.ColorWhite,
			BackgroundColor: termbox.ColorRed,
		},
	}, Area{0, 0, 1, 1})

	app.SetRootNode(&mainGrid)

	v := TestEventHandler{}
	app.RegisterEventListener(event.SystemKeyPress, v)

	app.Run()
}
