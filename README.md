# grip
A UI toolkit based on termbox. Most of the paradigms used here are inspired by CSS Grid Layout.

### Currently Implemented

* Event sending
* Find view (by path)
* Views
  * ActivityView
  * Grid
  * ImageView
  * Modal (+ Alert()/Confirm() helpers )
  * ProgressView
  * TableView (similar to Grid, but with borders)
  * TextView
  * InputView (incomplete)

### Example

```g
---
func OnEvent(app *App, e event.Event) {
	switch e.Type {
	case event.EventKeyPress:
		termboxEvent := e.Data.(termbox.Event)
		if termboxEvent.Type == termbox.EventKey {
			switch termboxEvent.Key {
			case termbox.KeyEsc:
				app.Confirm("Are you sure you want to quit?", func(app *App) {
					termbox.Close()
					os.Exit(0)
				}, func(app *App) {})
			case termbox.KeyF1:
				app.Alert("Hello", func(app *App) {})
			case termbox.KeyArrowLeft:
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
			case termbox.KeyArrowRight:
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
			case termbox.KeyArrowUp:
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
			case termbox.KeyArrowDown:
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
			case termbox.KeySpace:
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
			case termbox.KeyF2:
				inputView, err := app.Find(WildCardPath, "main-input")
				if err != nil {
					panic("Unable to find main-input")
				}

				mainInputView, isInputView := inputView.(*InputView)
				if !isInputView {
					panic("Unable to find main-input of type InputView")
				}
				mainInputView.Enable()
			}
		}
	}
}

  // App code
	app := App{}

	mainGrid := Grid{
		ColumnSizes:     []size.Size{size.Auto, size.WithPercent(30)},
		RowSizes:        []size.Size{size.Auto, size.WithPoints(1)},
		HasBackground:   true,
		BackgroundColor: termbox.ColorDefault,
	}

	sidebarGrid := NewTableView(
		[]size.Size{size.Auto},
		[]size.Size{size.WithPoints(1), size.WithPoints(1), size.WithPoints(1), size.WithPoints(1), size.WithPoints(1), size.Auto, size.Auto, size.Auto},
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
			app.Refresh()
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

	sidebarGrid.AddView("vertical-scroll", &ScrollView{
		Direction: ScrollDirectionVertical,
		View: &TextView{Text: `SidebarArea - Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eu consectetur lacus. Sed tincidunt eros non ultrices commodo. Sed ornare id dolor sed ultricies. Duis in est at nulla pretium mattis ac quis quam. Maecenas nibh nisi, rhoncus quis iaculis sit amet, semper et diam. Aenean pharetra ex non mi placerat rhoncus. Vivamus erat ante, suscipit vitae aliquet id, congue et dolor. Curabitur sed tortor tortor. Duis non sem et lacus ultrices finibus quis quis felis. Integer non elementum ante. Vestibulum vel augue ut tortor condimentum pulvinar eu blandit leo. Donec nibh nibh, tincidunt vitae risus a, consectetur suscipit felis. Quisque elementum velit nec mauris tristique, id malesuada tellus dictum.
			`, ForegroundColor: termbox.ColorWhite, BackgroundColor: termbox.ColorBlue},
		Size: 200,
	}, Area{0, 0, 5, 5})

	sidebarGrid.AddView("horizontal-scroll", &ScrollView{
		Direction: ScrollDirectionHorizontal,
		View: &TextView{Text: `SidebarArea - Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eu consectetur lacus. Sed tincidunt eros non ultrices commodo. Sed ornare id dolor sed ultricies. Duis in est at nulla pretium mattis ac quis quam. Maecenas nibh nisi, rhoncus quis iaculis sit amet, semper et diam. Aenean pharetra ex non mi placerat rhoncus. Vivamus erat ante, suscipit vitae aliquet id, congue et dolor. Curabitur sed tortor tortor. Duis non sem et lacus ultrices finibus quis quis felis. Integer non elementum ante. Vestibulum vel augue ut tortor condimentum pulvinar eu blandit leo. Donec nibh nibh, tincidunt vitae risus a, consectetur suscipit felis. Quisque elementum velit nec mauris tristique, id malesuada tellus dictum.
			`, ForegroundColor: termbox.ColorWhite, BackgroundColor: termbox.ColorGreen},
		Size: 200,
	}, Area{0, 0, 6, 6})

	sampleTable := NewTableView(
		[]size.Size{size.Auto, size.Auto, size.Auto},
		[]size.Size{size.Auto, size.Auto, size.Auto},
		termbox.ColorDefault,
		termbox.ColorDefault,
	)

	sidebarGrid.AddView("table", &sampleTable, Area{0, 0, 7, 7})

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
	}, Area{1, 1, 0, 0})

	gifView, err := NewAnimatedImageViewForGifFile("test.gif")
	if err != nil {
		panic(err.Error())
	}

	mainLayoutGrid.AddView("test-gif", &gifView, Area{0, 0, 0, 1})

	mainLayoutGrid.AddView("main-input", &InputView{
		TextView: TextView{
			Text:            "",
			ForegroundColor: termbox.ColorRed,
			BackgroundColor: termbox.ColorBlack,
		},
	}, Area{1, 1, 1, 1})

	mainGrid.AddView("main-layout", &mainLayoutGrid, Area{0, 0, 0, 0})
	mainGrid.AddView("sidebar-grid", &sidebarGrid, Area{1, 1, 0, 1})

	app.View = &mainGrid
	app.RegisterEventListener(event.EventKeyPress, OnEvent)

	app.Run()
---
```

produces the following

![Grid example](/screenshots/screenshot1.png?raw=true)
