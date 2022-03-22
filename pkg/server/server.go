package server

import (
	"io"
	"log"

	"github.com/mj-hagonoy/go-chat/protos/chat"
)

type ChatServiceServer struct {
	chat.UnimplementedChatServiceServer
	Channel map[string][]chan *chat.Message
}

func (server *ChatServiceServer) Join(ch *chat.Channel, msgStream chat.ChatService_JoinServer) error {
	msgChannel := make(chan *chat.Message)
	server.Channel[ch.Name] = append(server.Channel[ch.Name], msgChannel)

	for {
		select {
		case <-msgStream.Context().Done():
			return nil
		case msg := <-msgChannel:
			log.Printf("message received: %v\n", msg)
			msgStream.Send(msg)
		}
	}
}

func (server *ChatServiceServer) Send(msgStream chat.ChatService_SendServer) error {
	msg, err := msgStream.Recv()

	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	status := &chat.Status{Status: chat.Status_SENT}
	msgStream.SendAndClose(status)

	go func() {
		streams := server.Channel[msg.Channel.Name]
		for _, msgChan := range streams {
			msgChan <- msg
		}
	}()

	return nil
}
