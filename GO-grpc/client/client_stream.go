package main

import (
	"context"
	"log"
	"time"

	pb "github.com/jyotishmoy12/go-grpc/proto"
)

func callSayHelloClientStream(client pb.GreetServiceClient, names *pb.NamesList) {
	log.Printf("Client streaming started")

	stream, err := client.SayHelloClientStreaming(context.Background())
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

     for _, name :=range names.Names{
		req := &pb.HelloRequest{
		Name: name,
	}
	if err := stream.Send(req); err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	log.Printf("Sent request with name: %s", name)
	time.Sleep(2 * time.Second)
}
	 res, err := stream.CloseAndRecv()
	 log.Println("Closed stream and waiting for server response...")

	 log.Printf("Client streaming finished")
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}
	log.Printf("Received response: %s", res.Messages)
}
