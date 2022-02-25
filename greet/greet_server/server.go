package main

import (
	"context"
	"fmt"
	"hello-grpc/greet/greetpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("run on port 50051")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was inwoked with %v\n", req)
	firstname := req.GetGreeting().GetFirstName()
	result := "Hello " + firstname
	return &greetpb.GreetResponse{
		Result: result,
	}, nil
}
