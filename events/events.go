package events

import "time"

type EventType string

type Event interface{}

type Handler func(args Event)

var eventMap = make(map[EventType][]Handler)

var eventsChan = make(
	chan struct {
		Event
		EventType
	})

func Subscribe(h Handler, e EventType) {
	if len(eventMap[e]) == 0 {
		eventMap[e] = make([]Handler, 0)
	}

	eventMap[e] = append(eventMap[e], h)
}

func notify(t EventType, e Event) {
	for _, h := range eventMap[t] {
		go h(e)
	}
}

func PublishEvent(t EventType, e Event) {
	eventsChan <- struct {
		Event
		EventType
	}{Event: e, EventType: t}
}

// to be run as a goroutine
func Run() {
	for {
		select {
		case e := <-eventsChan:
			notify(e.EventType, e.Event)
		default:
			// ignore
		}

		time.Sleep(time.Millisecond)
	}

}
