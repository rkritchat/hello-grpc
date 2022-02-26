package main

import (
	"context"
	"fmt"
	"hello-grpc/calculator/calculatorpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("starting on :50051")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Listen failed: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalRequest) (*calculatorpb.CalResponse, error) {
	fmt.Printf("req: %v\n", req)
	return &calculatorpb.CalResponse{
		Result: req.FirstNumber + req.SecondNumber,
	}, nil
}

func (*server) Minus(ctx context.Context, req *calculatorpb.CalRequest) (*calculatorpb.CalResponse, error) {
	fmt.Printf("req: %v\n", req)
	return &calculatorpb.CalResponse{
		Result: req.FirstNumber - req.SecondNumber,
	}, nil
}

func (*server) Multiply(ctx context.Context, req *calculatorpb.CalRequest) (*calculatorpb.CalResponse, error) {
	fmt.Printf("req: %v\n", req)
	return &calculatorpb.CalResponse{
		Result: req.FirstNumber * req.SecondNumber,
	}, nil
}
