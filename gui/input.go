package gui

import (
	"github.com/jroimartin/gocui"
	"github.com/proh-gram-er/go-grpc-p2p-chat/events"
)

type Input struct {
	name      string
	x, y      int
	w         int
	h         int
	maxLength int
	tainted   bool
}

func NewInput(name string, x, y, w, h, maxLength int) *Input {
	return &Input{name: name, x: x, y: y, w: w, h: h, maxLength: maxLength, tainted: false}
}

func (i *Input) Layout(g *gocui.Gui) error {
	v, err := g.SetView(i.name, i.x, i.y, i.x+i.w, i.y+i.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	v.Editor = i
	v.Editable = true

	return nil
}

func (i *Input) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	cx, _ := v.Cursor()
	ox, _ := v.Origin()
	limit := ox+cx+1 > i.maxLength

	if !i.tainted {
		v.Clear()
		i.tainted = true
	}

	switch {

	case ch != 0 && mod == 0 && !limit:
		v.EditWrite(ch)
	case key == gocui.KeySpace && !limit:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)

	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)

	case key == gocui.KeyEnter:
		vb := v.ViewBuffer()

		if vb != "" {

			events.PublishEvent("message:send", events.SendMessageEvent{Message: vb})
			v.Clear()

			// ignore error because the coords are hard coded
			_ = v.SetCursor(0, 0)
		}
	}
}
