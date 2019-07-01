package main

import (
	pb "github.com/proh-gram-er/go-grpc-p2p-chat/protobuf"

	"github.com/proh-gram-er/go-grpc-p2p-chat/events"
	"time"
)

func sendMsg(e events.Event) {
	if e, ok := e.(events.SendMessageEvent); ok {
		_, err := client.SendMessage(ctx, &pb.Message{Content: e.Message})

		if err != nil {
			events.PublishEvent("modal:display", events.DisplayModalEvent{Message: "Connection to peer lost. Retrying"})
			time.Sleep(time.Millisecond)
			sendMsg(e)
		} else {
			events.PublishEvent("modal:hide", events.EmptyMessage{})
			events.PublishEvent("message:sent", events.SentMessageEvent{e.Message})
		}

	} else {
		// @TODO if e is not of SendMessageEvent type
		// ignore for the time being
	}
}

func subscribeListeners() {
	events.Subscribe(sendMsg, "message:send")
}
