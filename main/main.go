package main

import (
	"flag"
	"github.com/proh-gram-er/go-grpc-p2p-chat/chat"
	"github.com/proh-gram-er/go-grpc-p2p-chat/database"
	"github.com/proh-gram-er/go-grpc-p2p-chat/events"
	"github.com/proh-gram-er/go-grpc-p2p-chat/gui"
)

func main() {
	database.Init()

	_ = database.AddContact(database.Contact{
		Name:    "Don",
		Address: "127.0.0.1:8080",
	})

	_ = database.AddContact(database.Contact{
		Name:    "Jon",
		Address: "127.0.0.1:9080",
	})

	gui.SubscribeListeners()
	chat.SubscribeListeners()

	go events.Run()

	port := flag.String("port", "8080", "")
	flag.Parse()

	chat.SetListenPort(*port)
	go chat.Listen()

	gui.InitGui()
}
