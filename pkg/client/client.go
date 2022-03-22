package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/mj-hagonoy/go-chat/protos/chat"
)

func JoinChannel(ctx context.Context, client chat.ChatServiceClient) {
	sendersName := getSendersName(ctx)
	channel := chat.Channel{Name: getChannelName(ctx), SendersName: sendersName}
	stream, err := client.Join(ctx, &channel)
	if err != nil {
		log.Fatalf("client.Join: %v", err)
	}
	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive message from channel joining. \nErr: %v", err)
			}
			if sendersName != in.Sender.Name {
				fmt.Printf("MESSAGE: (%v) -> %v \n", in.Sender, in.Message)
			}
		}
	}()

	<-waitc
}

func SendMessage(ctx context.Context, client chat.ChatServiceClient, message string) {
	sendersName := getSendersName(ctx)
	stream, err := client.Send(ctx)
	if err != nil {
		log.Printf("Cannot send message: error: %v", err)
	}
	msg := chat.Message{
		Channel: &chat.Channel{
			Name:        getChannelName(ctx),
			SendersName: sendersName},
		Message: message,
		Sender:  &chat.Sender{Id: getSendersId(ctx), Name: getSendersName(ctx)},
	}
	if err := stream.Send(&msg); err != nil {
		log.Printf("stream.Send: %v", err)
	}

	ack, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("stream.CloseAndRecv: %v", err)
	}
	log.Printf("message received: %v", ack.Status)
}

type ctxValue string

func getChannelName(ctx context.Context) string {
	name := ctx.Value("channelName")
	if name == nil {
		return "default"
	}
	return name.(string)
}

func getSendersName(ctx context.Context) string {
	name := ctx.Value("sendersName")
	if name == nil {
		return "default"
	}
	return name.(string)
}

func getSendersId(ctx context.Context) int32 {
	id := ctx.Value("sendersId")
	if id == nil {
		return int32(time.Now().Unix())
	}
	ID, e := strconv.Atoi(id.(string))
	if e != nil {
		return int32(time.Now().Unix())
	}
	return int32(ID)
}
