package main

import (
	"log"
	"net"
     pb "github.com/jyotishmoy12/go-grpc/proto"
	"google.golang.org/grpc"
)


const (
	port =":8080"
)

type helloServer struct {
	pb.UnimplementedGreetServiceServer
}

func main(){
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	// Register your gRPC services here
	pb.RegisterGreetServiceServer(grpcServer, &helloServer{})
    
	log.Printf("server started at %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}