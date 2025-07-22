package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/jyotishmoy12/go-grpc/proto"
)

func callSayHelloBidirectionalStream(client pb.GreetServiceClient, names *pb.NamesList){
	log.Printf("Starting Bidirectional Streaming RPC...")
	stream, err := client.SayHelloBidiStreaming(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}
    // Create a channel to wait for the server's response
	// This will allow us to receive messages from the server while sending our own
	waitc := make(chan struct{})
    // Start a goroutine to receive messages from the server
	// This is necessary because the server can send messages at any time
	go func (){
		for{
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving message: %v", err)
			}
			log.Println(message)
		}
		// Close the channel when done
		// This will signal the main goroutine to finish
		close(waitc)
	}()

	for _, name := range names.Names {
		req := &pb.HelloRequest{
			Name: name,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error while sending message: %v", err)
		}
		time.Sleep(2 * time.Second)
	}
	stream.CloseSend()
	// Wait for the server to finish sending messages
	// This will block until the server closes the stream
	<- waitc
	log.Printf("Bidirectional Streaming RPC finished")
}