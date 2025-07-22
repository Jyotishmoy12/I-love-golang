package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	 pb "github.com/jyotishmoy12/go-grpc/proto"
)

const (
	port = ":8080"
)

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)

	names := &pb.NamesList{
		Names: []string{"Alice", "Bob", "Charlie", "Jyotishmoy"},
	}

	//callSayHello(client)

	//caySayHelloServerStream(client, names)
	//callSayHelloClientStream(client, names)
	callSayHelloBidirectionalStream(client, names)
}