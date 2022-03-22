package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mj-hagonoy/go-chat/pkg/client"
	"github.com/mj-hagonoy/go-chat/protos/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var channelName = flag.String("channel", "default", "Channel name for chatting")
	var senderName = flag.String("sender", "default", "Senders name")
	var senderID = flag.String("senderId", "1", "Senders name")
	var tcpServer = flag.String("server", ":5400", "Tcp server")
	flag.Parse()

	fmt.Println("--- CLIENT APP ---")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(*tcpServer, opts...)
	if err != nil {
		log.Fatalf("Fail to dial: %v", err)
	}

	defer conn.Close()
	type ctxValue string
	ctx := context.Background()
	ctx = context.WithValue(ctx, "channelName", *channelName)
	ctx = context.WithValue(ctx, "sendersName", *senderName)
	ctx = context.WithValue(ctx, "sendersId", *senderID)

	chatClient := chat.NewChatServiceClient(conn)

	go client.JoinChannel(ctx, chatClient)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		go client.SendMessage(ctx, chatClient, scanner.Text())
	}
}
