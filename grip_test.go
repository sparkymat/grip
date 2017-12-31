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
The Lord of the Rings is a film series consisting of three high fantasy adventure films directed by Peter Jackson. They are based on the novel The Lord of the Rings by J. R. R. Tolkien. The films are subtitled The Fellowship of the Ring (2001), The Two Towers (2002) and The Return of the King (2003). They are a New Zealand-American venture produced by WingNut Films and The Saul Zaentz Company and distributed by New Line Cinema.

Considered to be one of the biggest and most ambitious film projects ever undertaken, with an overall budget of $281 million (some sources say $310-$330 million),[2] the entire project took eight years, with the filming for all three films done simultaneously and entirely in New Zealand, Jackson's native country.[3] Each film in the series also had special extended editions released on DVD a year after their respective theatrical releases. While the films follow the book's general storyline, they do omit some of the novel's plot elements and include some additions to and deviations from the source material.
		`, BackgroundColor: termbox.ColorGreen, ForegroundColor: termbox.ColorBlack},
		Size: 100,
	}, Area{1, 1, 0, 0})
	mainLayoutGrid.AddView("horizontal-scroll", &ScrollView{
		Direction: ScrollDirectionHorizontal,
		View: &TextView{Text: `
Futurama is an American animated science fiction comedy series created by Matt Groening for the Fox Broadcasting Company. The series follows the adventures of a late-20th-century New York City pizza delivery boy, Philip J. Fry, who, after being unwittingly cryogenically frozen for one thousand years, finds employment at Planet Express, an interplanetary delivery company in the retro-futuristic 31st century. The series was envisioned by Groening in the mid-1990s while working on The Simpsons; he later brought David X. Cohen aboard to develop storylines and characters to pitch the show to Fox.

In the United States, the series aired on Fox from March 28, 1999, to August 10, 2003, before ceasing production. Futurama also aired in reruns on Cartoon Network's Adult Swim from 2003 to 2007, until the network's contract expired. It was revived in 2007 as four direct-to-video films; the last of which was released in early 2009. Comedy Central entered into an agreement with 20th Century Fox Television to syndicate the existing episodes and air the films as 16 new, half-hour episodes, constituting a fifth season.[1][2]

In June 2009, producing studio 20th Century Fox announced that Comedy Central had picked up the show for 26 new half-hour episodes, which began airing in 2010 and 2011.[3][4] The show was renewed for a seventh season, with the first half airing in June 2012 and the second set for mid-2013.[5][6] It was later revealed that the seventh season would be the final season, as Comedy Central announced that they would not be commissioning any further episodes. The series finale aired on September 4, 2013.[7] While Groening has said he will try to get it picked up by another network,[8] David X. Cohen stated that the episode "Meanwhile" would be the last episode of season 7 and also the series finale.[9] A 42-minute audio-only episode featuring its original cast members was released on September 14, 2017, as an episode of The Nerdist Podcast entitled Futurama: Worlds of Tomorrow Presents, RADIORAMA!.[10]
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
