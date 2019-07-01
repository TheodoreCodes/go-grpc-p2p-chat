package gui

import (
	"github.com/jroimartin/gocui"
	"log"
)

var g *gocui.Gui

func setFocus(name string) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		_, err := g.SetCurrentView(name)
		return err
	}
}

func guiLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	input := NewInput("messageInput", 0, 9*maxY/10, maxX-1, maxY/10-1, 256)
	focus := gocui.ManagerFunc(setFocus("messageInput"))

	g.SetManager(input, focus)

	if v, err := g.SetView("conversation", 0, 0, maxX-1, 9*maxY/10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Editable = true
	}

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("messageInput")
		if err != nil {
			return err
		}

		v.Title = "Message"

		return nil
	})

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func InitGui() {
	var err error

	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	err = guiLayout(g)
	if err != nil {
		log.Panic(err)
	}

	err = keybindings(g)
	if err != nil {
		log.Panic(err)
	}

	subscribeListeners()

	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
}
