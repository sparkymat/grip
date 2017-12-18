package grip

import (
	termbox "github.com/nsf/termbox-go"
)

type app struct {
	rootNode *grid
}

func New() app {
	return app{}
}

func (a *app) SetRootNode(node *grid) {
	a.rootNode = node
}

func (a app) Run(eventChannel chan termbox.Event) error {
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.SetOutputMode(termbox.Output256)

	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	width, height := termbox.Size()
	a.rootNode.Resize(0, 0, uint32(width), uint32(height))
	a.rootNode.Draw()

	termbox.Flush()

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			width, height := termbox.Size()
			a.rootNode.Resize(0, 0, uint32(width), uint32(height))
			a.rootNode.Draw()
			termbox.Flush()
		case termbox.EventKey:
			eventChannel <- ev
		case termbox.EventError:
			panic(ev.Err)
		}
	}

	return err
}
