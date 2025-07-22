package main

import (
	"context"
	"io"
	"log"

	pb "github.com/jyotishmoy12/go-grpc/proto"
)

func caySayHelloServerStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Streaming started")
	stream, err := client.SayHelloServerStreaming(context.Background(), names)

	if err != nil {
		log.Printf("Error while calling SayHelloServerStreaming: %v", err)
		return
	}

	for {
		message, err := stream.Recv()
		
		if err == io.EOF{
			log.Printf("Streaming finished")
			break
		}
		if err != nil {
			log.Printf("Error while receiving message: %v", err)
		} 
		log.Printf("Received message: %s", message.Message)
	}
	log.Printf("Streaming completed")
}