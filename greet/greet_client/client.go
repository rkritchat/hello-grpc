package main

import (
	"context"
	"fmt"
	"hello-grpc/greet/greetpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("starting to do a Unary RPC...")
	req := greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "rkritchat",
			LastName:  "rojanaphruk",
		},
	}
	resp, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatalf("error while invoke Greet rpc: %v", err)
	}

	fmt.Printf("resp: %v", resp.Result)
}
