package main

import (
	"io"
	"log"

	pb "github.com/jyotishmoy12/go-grpc/proto"
)

func (s * helloServer) SayHelloBidiStreaming(stream pb.GreetService_SayHelloBidiStreamingServer)error{
	for{
		req, err :=stream.Recv()
		if err == io.EOF{
			return nil // End of stream
		}
		if err != nil {
			return err // Error receiving message
		}
		log.Printf("Got request with name: %v", req.Name)

		res := &pb.HelloResponse{
			Message: "Hello " + req.Name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}
}