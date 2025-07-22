package main

import (
	"log"
	"time"

	pb "github.com/jyotishmoy12/go-grpc/proto"
)

func (s *helloServer) SayHelloServerStreaming(req *pb.NamesList, stream pb.GreetService_SayHelloServerStreamingServer) error{
	log.Printf("got request with names: %v", req.Names)

	for _, name := range req.Names {
		res := &pb.HelloResponse{
			Message: "Hello " + name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(time.Second * 2) // Simulate some processing delay
	}

	return nil


}