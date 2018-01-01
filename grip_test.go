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

func OnEvent(app *App, e event.Event) {
	switch e.Type {
	case event.SystemKeyPress:
		termboxEvent := e.Data.(termbox.Event)
		if termboxEvent.Type == termbox.EventKey && termboxEvent.Key == termbox.KeyEsc {
			app.Confirm("Are you sure you want to quit?", func(app *App) {
				termbox.Close()
				os.Exit(0)
			}, func(app *App) {})
		} else if termboxEvent.Type == termbox.EventKey && termboxEvent.Key == termbox.KeyF1 {
			app.Alert("Hello", func(app *App) {})
		} else if termboxEvent.Type == termbox.EventKey && termboxEvent.Key == termbox.KeyArrowLeft {
			scrollViewView, err := app.Find(WildCardPath, "horizontal-scroll")
			if err != nil {
				panic("Unable to find horizontal-scroll")
			}

			scrollView, isScroll := scrollViewView.(*ScrollView)

			if !isScroll {
				panic("Unable to find horizontal-scroll of type ScrollView")
			}

			scrollView.ScrollTo(scrollView.GetScrollPosition() - 1)
			scrollView.Draw()
		} else if termboxEvent.Type == termbox.EventKey && termboxEvent.Key == termbox.KeyArrowRight {
			scrollViewView, err := app.Find(WildCardPath, "horizontal-scroll")
			if err != nil {
				panic("Unable to find horizontal-scroll")
			}

			scrollView, isScroll := scrollViewView.(*ScrollView)

			if !isScroll {
				panic("Unable to find horizontal-scroll of type ScrollView")
			}

			scrollView.ScrollTo(scrollView.GetScrollPosition() + 1)
			scrollView.Draw()
		} else if termboxEvent.Type == termbox.EventKey && termboxEvent.Key == termbox.KeyArrowUp {
			scrollViewView, err := app.Find(WildCardPath, "vertical-scroll")
			if err != nil {
				panic("Unable to find vertical-scroll")
			}

			scrollView, isScroll := scrollViewView.(*ScrollView)

			if !isScroll {
				panic("Unable to find vertical-scroll of type ScrollView")
			}

			scrollView.ScrollTo(scrollView.GetScrollPosition() - 1)
			scrollView.Draw()
		} else if termboxEvent.Type == termbox.EventKey && termboxEvent.Key == termbox.KeyArrowDown {
			scrollViewView, err := app.Find(WildCardPath, "vertical-scroll")
			if err != nil {
				panic("Unable to find vertical-scroll")
			}

			scrollView, isScroll := scrollViewView.(*ScrollView)

			if !isScroll {
				panic("Unable to find vertical-scroll of type ScrollView")
			}

			scrollView.ScrollTo(scrollView.GetScrollPosition() + 1)
			scrollView.Draw()
		} else if termboxEvent.Type == termbox.EventKey && termboxEvent.Key == termbox.KeySpace {
			manualProgressView, err := app.Find(WildCardPath, "manual-progress")
			if err != nil {
				panic("Unable to find manual-progress")
			}

			manualProgress, isProgress := manualProgressView.(*ProgressView)

			if !isProgress {
				panic("Unable to find manual-progress of type ProgressView")
			}

			manualProgress.CurrentValue += (rand.Int() % 200)
			if manualProgress.CurrentValue > 1000 {
				manualProgress.CurrentValue = 0
			}
			manualProgress.Draw()
		}
		break
	}
}

func TestSanity(t *testing.T) {
	app := App{}

	mainGrid := Grid{
		ColumnSizes:     []size.Size{size.Auto, size.WithPercent(30)},
		RowSizes:        []size.Size{size.Auto, size.WithPoints(1)},
		HasBackground:   true,
		BackgroundColor: termbox.ColorDefault,
	}

	sidebarGrid := NewTableView(
		[]size.Size{size.Auto},
		[]size.Size{size.WithPoints(1), size.WithPoints(1), size.WithPoints(1), size.WithPoints(1), size.WithPoints(1), size.Auto, size.Auto},
		termbox.ColorBlue,
		termbox.ColorDefault,
	)

	mainLayoutGrid := Grid{
		ColumnSizes: []size.Size{size.Auto, size.Auto},
		RowSizes:    []size.Size{size.Auto, size.Auto},
	}

	sidebarGrid.AddView("name-text", &TextView{
		Text:            "Name: Adam",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorCyan,
		TextAlignment:   TextAlignmentCenter,
	}, Area{0, 0, 0, 0})

	sidebarGrid.AddView("class-text", &TextView{
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
			progress.CurrentValue += (rand.Int() % 200)
			if progress.CurrentValue > 1000 {
				progress.CurrentValue = 0
			}
			progress.Draw()
		}
	}()
	sidebarGrid.AddView("progress", &progress, Area{0, 0, 2, 2})

	progress2 := ProgressView{
		Type:            ProgressViewTypePercentage,
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorMagenta,
		MinimumValue:    0,
		MaximumValue:    1000,
		CurrentValue:    0,
	}
	sidebarGrid.AddView("manual-progress", &progress2, Area{0, 0, 3, 3})

	sidebarGrid.AddView("activity", &ActivityView{
		BackgroundColor: termbox.ColorRed,
		ForegroundColor: termbox.ColorWhite,
		Text:            "Loading...",
	}, Area{0, 0, 4, 4})

	sidebarGrid.AddView("long-text", &TextView{
		Text:            "SidebarArea - Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eu consectetur lacus. Sed tincidunt eros non ultrices commodo. Sed ornare id dolor sed ultricies. Duis in est at nulla pretium mattis ac quis quam. Maecenas nibh nisi, rhoncus quis iaculis sit amet, semper et diam. Aenean pharetra ex non mi placerat rhoncus. Vivamus erat ante, suscipit vitae aliquet id, congue et dolor. Curabitur sed tortor tortor. Duis non sem et lacus ultrices finibus quis quis felis. Integer non elementum ante. Vestibulum vel augue ut tortor condimentum pulvinar eu blandit leo. Donec nibh nibh, tincidunt vitae risus a, consectetur suscipit felis. Quisque elementum velit nec mauris tristique, id malesuada tellus dictum.",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorBlue,
	}, Area{0, 0, 5, 5})

	sampleTable := NewTableView(
		[]size.Size{size.Auto, size.Auto, size.Auto},
		[]size.Size{size.Auto, size.Auto, size.Auto},
		termbox.ColorDefault,
		termbox.ColorDefault,
	)

	sidebarGrid.AddView("table", &sampleTable, Area{0, 0, 6, 6})

	f, err := os.Open("test.jpg")
	if err != nil {
		t.Error(err.Error())
	}

	img, _, err := image.Decode(f)
	if err != nil {
		t.Error(err.Error())
	}

	mainLayoutGrid.AddView("test-image", &ImageView{
		Image:           img,
		ForegroundColor: termbox.ColorDefault,
		BackgroundColor: termbox.ColorDefault,
	}, Area{0, 0, 0, 0})
	mainLayoutGrid.AddView("vertical-scroll", &ScrollView{
		Direction: ScrollDirectionVertical,
		View: &TextView{Text: `
The language was announced in November 2009. It is used in some of Google's production systems, as well as by other firms.[15]

Two major implementations exist:

    Google's Go compiler, "gc", developed as open-source software. It targets various platforms including Linux, OS X, Windows, various BSD and Unix versions, and also (since 2015) mobile devices (including smartphones).[16]
    A second compiler, gccgo, is a GCC frontend.[17][18]

The "gc" toolchain has been self-hosting since version 1.5.[19]

Go originated as an experiment by Google engineers Robert Griesemer, Rob Pike, and Ken Thompson to design a new programming language that would resolve common criticisms of other languages while maintaining their positive characteristics. The developers envisaged the new language as:[20]

    statically typed, scalable to large systems (as Java and C++)
    productive and readable, without too many mandatory keywords and repetition[21] ("light on the page" like dynamic languages)
    not requiring integrated development environments, but supporting them well
    supporting networking and multiprocessing

In later interviews, all three of the language designers cited their shared dislike of C++'s complexity as a primary motivation for designing a new language.[22][23][24]

Google released Go 1.0 in March 2012.[25][26]
		`, BackgroundColor: termbox.ColorGreen, ForegroundColor: termbox.ColorBlack},
		Size: 100,
	}, Area{1, 1, 0, 0})
	mainLayoutGrid.AddView("horizontal-scroll", &ScrollView{
		Direction: ScrollDirectionHorizontal,
		View: &TextView{Text: `
The language was announced in November 2009. It is used in some of Google's production systems, as well as by other firms.[15]

Two major implementations exist:

    Google's Go compiler, "gc", developed as open-source software. It targets various platforms including Linux, OS X, Windows, various BSD and Unix versions, and also (since 2015) mobile devices (including smartphones).[16]
    A second compiler, gccgo, is a GCC frontend.[17][18]

The "gc" toolchain has been self-hosting since version 1.5.[19]

Go originated as an experiment by Google engineers Robert Griesemer, Rob Pike, and Ken Thompson to design a new programming language that would resolve common criticisms of other languages while maintaining their positive characteristics. The developers envisaged the new language as:[20]

    statically typed, scalable to large systems (as Java and C++)
    productive and readable, without too many mandatory keywords and repetition[21] ("light on the page" like dynamic languages)
    not requiring integrated development environments, but supporting them well
    supporting networking and multiprocessing

In later interviews, all three of the language designers cited their shared dislike of C++'s complexity as a primary motivation for designing a new language.[22][23][24]

Google released Go 1.0 in March 2012.[25][26]
		`, BackgroundColor: termbox.ColorBlue, ForegroundColor: termbox.ColorBlack},
		Size: 200,
	}, Area{0, 0, 1, 1})

	mainGrid.AddView("main-layout", &mainLayoutGrid, Area{0, 0, 0, 0})
	mainGrid.AddView("sidebar-grid", &sidebarGrid, Area{1, 1, 0, 1})
	mainGrid.AddView("main-input", &InputView{
		TextView: TextView{
			Text:            "",
			ForegroundColor: termbox.ColorRed,
			BackgroundColor: termbox.ColorBlack,
		},
	}, Area{0, 0, 1, 1})

	app.SetContainer(&mainGrid)

	app.RegisterGlobalEventListener(event.SystemKeyPress, OnEvent)

	app.Run()
}
