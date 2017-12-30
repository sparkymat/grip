package grip

import (
	"image"
	"math/rand"
	"os"
	"testing"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/sparkymat/grip/event"
	"github.com/sparkymat/grip/size"
)

type TestEventHandler struct {
	App *App
}

func (t TestEventHandler) OnEvent(app *App, e event.Event) {
	switch e.Type {
	case event.SystemKeyPress:
		termboxEvent := e.Data.(termbox.Event)
		if termboxEvent.Type == termbox.EventKey && termboxEvent.Key == termbox.KeyEsc {
			termbox.Close()
			os.Exit(0)
		} else if termboxEvent.Type == termbox.EventKey && termboxEvent.Key == termbox.KeyF1 {
			t.App.Alert("Hello", func() {})
		}
		break
	}
}

func TestSanity(t *testing.T) {
	app := App{}

	mainGrid := Grid{
		ColumnSizes: []size.Size{size.Auto, size.WithPercent(30)},
		RowSizes:    []size.Size{size.Auto, size.WithPoints(1)},
	}

	sidebarGrid := NewTableView(
		[]size.Size{size.Auto},
		[]size.Size{size.WithPoints(1), size.WithPoints(1), size.WithPoints(1), size.WithPoints(1), size.Auto, size.Auto},
		termbox.ColorBlue,
		termbox.ColorDefault,
	)

	sidebarGrid.AddView(&TextView{
		Text:            "Name: Adam",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorCyan,
		TextAlignment:   TextAlignmentCenter,
	}, Area{0, 0, 0, 0})

	sidebarGrid.AddView(&TextView{
		Text:            "Class: Warlock",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorYellow,
		TextAlignment:   TextAlignmentRight,
	}, Area{0, 0, 1, 1})

	progress := ProgressView{
		Type:            ProgressViewTypePercentage,
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorMagenta,
		MinimumValue:    0,
		MaximumValue:    1000,
		CurrentValue:    0,
	}

	progressTimer := time.NewTicker(time.Millisecond * 250)
	go func() {
		for _ = range progressTimer.C {
			progress.CurrentValue += (rand.Int31() % 200)
			if progress.CurrentValue > 1000 {
				progress.CurrentValue = 0
			}
			progress.Draw()
		}
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

	sampleTable := NewTableView(
		[]size.Size{size.Auto, size.Auto, size.Auto},
		[]size.Size{size.Auto, size.Auto, size.Auto},
		termbox.ColorDefault,
		termbox.ColorDefault,
	)

	sidebarGrid.AddView(&sampleTable, Area{0, 0, 5, 5})

	f, err := os.Open("test.jpg")
	if err != nil {
		t.Error(err.Error())
	}

	img, _, err := image.Decode(f)
	if err != nil {
		t.Error(err.Error())
	}

	mainGrid.AddView(&ImageView{
		Image:           img,
		ForegroundColor: termbox.ColorDefault,
		BackgroundColor: termbox.ColorDefault,
	}, Area{0, 0, 0, 0})

	mainGrid.AddView(&sidebarGrid, Area{1, 1, 0, 1})

	mainGrid.AddView(&InputView{
		TextView: TextView{
			Text:            "",
			ForegroundColor: termbox.ColorWhite,
			BackgroundColor: termbox.ColorRed,
		},
	}, Area{0, 0, 1, 1})

	app.SetContainer(&mainGrid)

	v := TestEventHandler{App: &app}
	app.RegisterGlobalEventListener(event.SystemKeyPress, v)

	app.Run()
}
