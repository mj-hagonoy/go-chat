package main

import (
	"log"
	"net"

	"github.com/mj-hagonoy/go-chat/pkg/server"
	"github.com/mj-hagonoy/go-chat/protos/chat"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:5400")
	if err != nil {
		log.Fatalf("Failed to listed: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chat.RegisterChatServiceServer(grpcServer, &server.ChatServiceServer{
		Channel: make(map[string][]chan *chat.Message),
	})
	grpcServer.Serve(lis)
}
