# grip
A UI toolkit based on termbox. Most of the paradigms used here are inspired by CSS Grid Layout.

### Currently Implemented
* Grid
* TextView

### Example

``` go
	app := grip.New()

	mainGrid := grip.NewGrid(
		[]size.Size{size.Auto, size.WithPercent(30)},
		[]size.Size{size.Auto, size.WithPoints(1)},
		grip.Area{0, 0, 0, 0},
	)

	sidebarGrid := grip.NewGrid(
		[]size.Size{size.Auto},
		[]size.Size{size.WithPoints(1), size.WithPoints(1), size.Auto},
		grip.Area{1, 1, 0, 1},
	)
	sidebarGrid.AddView(&grip.TextView{
		Text:            "Name: Adam",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorCyan,
		Area:            grip.Area{0, 0, 0, 0},
	})
	sidebarGrid.AddView(&grip.TextView{
		Text:            "Class: Warlock",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorYellow,
		Area:            grip.Area{0, 0, 1, 1},
	})
	sidebarGrid.AddView(&grip.TextView{
		Text:            "SidebarArea - Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed eu consectetur lacus. Sed tincidunt eros non ultrices commodo. Sed ornare id dolor sed ultricies. Duis in est at nulla pretium mattis ac quis quam. Maecenas nibh nisi, rhoncus quis iaculis sit amet, semper et diam. Aenean pharetra ex non mi placerat rhoncus. Vivamus erat ante, suscipit vitae aliquet id, congue et dolor. Curabitur sed tortor tortor. Duis non sem et lacus ultrices finibus quis quis felis. Integer non elementum ante. Vestibulum vel augue ut tortor condimentum pulvinar eu blandit leo. Donec nibh nibh, tincidunt vitae risus a, consectetur suscipit felis. Quisque elementum velit nec mauris tristique, id malesuada tellus dictum.",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorBlue,
		Area:            grip.Area{0, 0, 2, 2},
	})

	mainGrid.AddView(&grip.TextView{
		Text:            "MainArea",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorGreen,
		Area:            grip.Area{0, 0, 0, 0},
	})
	mainGrid.AddView(&sidebarGrid)
	mainGrid.AddView(&TextView{
		Text:            "CommandArea",
		ForegroundColor: termbox.ColorWhite,
		BackgroundColor: termbox.ColorRed,
		Area:            grip.Area{0, 0, 1, 1},
	})

	app.SetRootNode(&mainGrid)

	eventChannel := make(chan termbox.Event)
	go HandleEvents(eventChannel)

	app.Run(eventChannel)
```

produces the following

![Grid example](/screenshots/screenshot1.png?raw=true)
