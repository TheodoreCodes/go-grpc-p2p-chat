package chat

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/sortkeys"
	"github.com/proh-gram-er/go-grpc-p2p-chat/database"
	"github.com/proh-gram-er/go-grpc-p2p-chat/events"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

type server struct{}

var (
	client        ChatClient
	ctx           context.Context
	listenPort    string
	activeContact database.Contact
)

func (s *server) SendMessage(ctx context.Context, in *Message) (*Response, error) {
	events.PublishEvent("message:receive", events.ReceiveMessageEvent{Message: in.Content})
	return &Response{Received: true}, nil
}

func Listen() {
	lis, err := net.Listen("tcp", strings.Join([]string{"localhost", listenPort}, ":"))

	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	RegisterChatServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %s", err)
	}
}

func Connect(address *string) {
	conn, err := grpc.Dial(*address, grpc.WithInsecure())

	if err != nil {
		log.Panicf("Can't connect to peer: %s", err)
	}

	//defer conn.Close()

	client = NewChatClient(conn)

	ctx, _ = context.WithCancel(context.TODO())

	//defer cancel()
}

func SetListenPort(port string) {
	listenPort = port
}

func SetActiveContact(address *string) {
	activeContact, _ = database.GetContact(address)
}

func GetActiveConversation() string {
	conversation, err := database.GetConversation(activeContact.Address)

	if err != nil {
		log.Panic(err)
	}

	var displayConversation []string
	var sender string

	var keys []int64

	for k := range conversation {
		keys = append(keys, k)
	}

	sortkeys.Int64s(keys)

	for _, k := range keys {
		msg := conversation[k]

		if msg.Sender == 0 {
			sender = "Me"
		} else {
			sender = activeContact.Name
		}
		displayConversation = append(displayConversation, fmt.Sprintf("%s: %s", sender, msg.Content))
	}

	return strings.Join(displayConversation, "\n")

}
