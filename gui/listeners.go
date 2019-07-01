package gui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/proh-gram-er/go-grpc-p2p-chat/events"
	"log"
)

func updateConversation(e events.Event) {
	g.Update(func(g *gocui.Gui) error {
		v, err := g.View("conversation")
		if err != nil {
			// @TODO handle error
			return err
		}

		var msg string
		if se, ok := e.(events.SentMessageEvent); ok {
			msg = fmt.Sprintf("You: %s", se.Message)
		} else if re, ok := e.(events.ReceiveMessageEvent); ok {
			msg = fmt.Sprintf("Friend: %s", re.Message)
		}
		_, err = fmt.Fprintln(v, msg)

		if err != nil {
			return err
		}

		return nil
	})
}

var isModalDisplayed = false

func displayModal(e events.Event) {
	if !isModalDisplayed {
		maxX, maxY := g.Size()

		var modalMsg string
		if e, ok := e.(events.DisplayModalEvent); ok {
			modalMsg = e.Message

			g.Update(func(g *gocui.Gui) error {
				v, err := g.SetView("modal",
					maxX/2-50,
					maxY/2-2,
					maxX/2+50,
					maxY/2+2,
				)

				if err != nil {
					// @TODO handle error
					if err != gocui.ErrUnknownView {
						return err
					}
				}

				v.Wrap = true
				v.Frame = true

				// explicitly ignore error
				_, _ = fmt.Fprint(v, modalMsg)

				return nil
			})
		}

		isModalDisplayed = true
	}
}

func hideModal(e events.Event) {
	if isModalDisplayed {
		err := g.DeleteView("modal")
		if err != nil && err != gocui.ErrUnknownView {
			log.Panicf("could not close moda: %s", err)
		}

		isModalDisplayed = false
	}
}

func subscribeListeners() {
	events.Subscribe(updateConversation, "message:sent")
	events.Subscribe(updateConversation, "message:receive")

	events.Subscribe(displayModal, "modal:display")
	events.Subscribe(hideModal, "modal:hide")
}
