package main

import (
	"context"
	"flag"
	"github.com/proh-gram-er/go-grpc-p2p-chat/events"
	"github.com/proh-gram-er/go-grpc-p2p-chat/gui"
	pb "github.com/proh-gram-er/go-grpc-p2p-chat/protobuf"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

var (
	peer   *string
	client pb.ChatClient
	ctx    context.Context
)

type server struct{}

func (s *server) SendMessage(ctx context.Context, in *pb.Message) (*pb.Response, error) {
	events.PublishEvent("message:receive", events.ReceiveMessageEvent{Message: in.Content})
	return &pb.Response{Received: true}, nil
}

func listen(port string) {
	lis, err := net.Listen("tcp", strings.Join([]string{"localhost", port}, ":"))

	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Panicf("failed to serve: %s", err)
	}
}

func main() {
	port := flag.String("port", "8080", "Port on which the app listens for incoming messages")
	peer = flag.String("peer", "", "Address of peer (ip:port)")

	flag.Parse()

	conn, err := grpc.Dial(*peer, grpc.WithInsecure())

	if err != nil {
		log.Panicf("Can't connect to peer: %s", err)
	}

	defer conn.Close()

	client = pb.NewChatClient(conn)

	var cancel context.CancelFunc

	ctx, cancel = context.WithCancel(context.TODO())

	defer cancel()

	go listen(*port)

	subscribeListeners()

	go events.Run()

	gui.InitGui()
}
