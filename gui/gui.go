package gui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/proh-gram-er/go-grpc-p2p-chat/chat"
	"github.com/proh-gram-er/go-grpc-p2p-chat/database"
	"log"
	"strings"
)

var (
	g             *gocui.Gui
	contactsCount int
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()

		if cy < contactsCount-1 {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func highlightContacts(gui *gocui.Gui, view *gocui.View) error {
	_, err := g.SetCurrentView("contactsList")

	return err
}

func selectContact(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		return err
	}

	contact := strings.Split(l, "(")
	name := strings.TrimSpace(contact[0])
	address := strings.TrimSuffix(contact[1], ")")

	chat.Connect(&address)

	view, err := g.View("conversation")

	chat.SetActiveContact(&address)

	view.Clear()
	fmt.Fprintln(view, chat.GetActiveConversation())

	view.Title = name

	v, err = g.SetCurrentView("messageInput")

	if err != nil {
		log.Panic(err)
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	msgInput := NewInput("messageInput", maxX/4, 9*maxY/10, 3*maxX/4-1, maxY/10-1, 256)

	g.SetManager(msgInput)

	g.Cursor = true

	if v, err := g.SetView("contactsList", 0, 0, maxX/4-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Contacts"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		contacts, err := database.GetContacts()

		if err != nil {
			log.Panic(err)
		}

		if err != nil {
			return err
		}

		for _, contact := range contacts {
			fmt.Fprintln(v, fmt.Sprintf("%s (%s)", contact.Name, contact.Address))
		}

		contactsCount = len(contacts)

		if _, err := g.SetCurrentView("contactsList"); err != nil {
			log.Panic(err)
		}
	}

	if v, err := g.SetView("conversation", maxX/4, 0, maxX-1, 9*maxY/10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Autoscroll = true
	}

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("contactsList", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("contactsList", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("contactsList", gocui.KeyEnter, gocui.ModNone, selectContact); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, highlightContacts); err != nil {
		return err
	}

	return nil
}

func InitGui() {
	var err error

	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil && err != gocui.ErrUnknownView {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	err = layout(g)
	if err != nil {
		log.Panic(err)
	}

	err = keybindings(g)
	if err != nil {
		log.Panic(err)
	}

	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Panic(err)
	}
}
