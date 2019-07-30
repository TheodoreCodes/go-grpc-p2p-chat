package chat

import (
	"github.com/proh-gram-er/go-grpc-p2p-chat/database"
	"github.com/proh-gram-er/go-grpc-p2p-chat/events"
	"log"
	"os"
	"time"
)

func sendMsg(e events.Event) {
	if e, ok := e.(events.SendMessageEvent); ok {
		_, err := client.SendMessage(ctx, &Message{Content: e.Message})

		if err != nil {
			events.PublishEvent("modal:display", events.DisplayModalEvent{Message: "Connection to peer lost. Retrying"})
			time.Sleep(1 * time.Millisecond)
			sendMsg(e)
		} else {
			err = database.UpdateConversation(activeContact.Address, 0, e.Message)

			if err != nil {
				// @TODO figure out how to handle error
				log.Panic(err)
			}

			events.PublishEvent("modal:hide", events.EmptyMessage{})
			events.PublishEvent("conversation:refresh", events.EmptyMessage{})
		}

	} else {
		// @TODO if e is not of SendMessageEvent type
		// ignore for the time being
	}
}

func receiveMsg(e events.Event) {
	if e, ok := e.(events.ReceiveMessageEvent); ok {
		err := database.UpdateConversation(activeContact.Address, 1, e.Message)
		if err != nil {
			// @TODO figure out how to handle error
			os.Exit(1)
		}
		events.PublishEvent("conversation:refresh", events.EmptyMessage{})

	}
}

func SubscribeListeners() {
	events.Subscribe(sendMsg, "message:send")
	events.Subscribe(receiveMsg, "message:receive")
}
