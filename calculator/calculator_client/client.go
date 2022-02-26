package main

import (
	"context"
	"fmt"
	"hello-grpc/calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello from client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	plus(cc)
	minus(cc)
	multiply(cc)
}

func plus(cc *grpc.ClientConn) {
	fmt.Println("start call plus")
	c := calculatorpb.NewCalculatorServiceClient(cc)
	resp, err := c.Sum(context.Background(), &calculatorpb.CalRequest{
		FirstNumber:  5,
		SecondNumber: 2,
	})
	if err != nil {
		log.Fatalf("err while invoke Sum: %v", err)
	}
	fmt.Println(resp)
}

func minus(cc *grpc.ClientConn) {
	fmt.Println()
	fmt.Println("start call minus")
	c := calculatorpb.NewCalculatorServiceClient(cc)
	resp, err := c.Minus(context.Background(), &calculatorpb.CalRequest{
		FirstNumber:  10,
		SecondNumber: 2,
	})

	if err != nil {
		log.Fatalf("err while invoke minus: %v", err)
	}
	fmt.Println(resp)
}

func multiply(cc *grpc.ClientConn) {
	fmt.Println()
	fmt.Println("start call multiply")
	c := calculatorpb.NewCalculatorServiceClient(cc)
	resp, err := c.Multiply(context.Background(), &calculatorpb.CalRequest{
		FirstNumber:  5,
		SecondNumber: 6,
	})
	if err != nil {
		log.Fatalf("err while invoke Multiply: %v", err)
	}
	fmt.Println(resp)
}
