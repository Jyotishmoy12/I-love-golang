package main

import (
	"io"
	"log"

	pb "github.com/jyotishmoy12/go-grpc/proto"
)

func (s *helloServer) SayHelloClientStreaming(stream pb.GreetService_SayHelloClientStreamingServer) error {
	var messages []string

	for {
		log.Println("Waiting to receive from client...")
		req, err := stream.Recv()
		if err == io.EOF{
			log.Println("Client finished sending, now sending back the response...")
			return stream.SendAndClose(&pb.MessagesList{Messages: messages})
		}
		if err != nil {
			return err
		}
		log.Printf("Got request with name: %v", req.Name)
		messages = append(messages, "Hello" + req.Name)
	}
}