package main

import (
	"context"
	"fmt"
	"hello-grpc/greet/greetpb"
	"io"
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
	//doUnary(c)

	doServerStreaming(c)
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

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	resStream, err := c.GreetManyTimes(context.Background(), &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "rkritchat",
			LastName:  "rojanaphruk",
		},
	})
	if err != nil {
		log.Fatalf("err while invoke GreetManyTimes: %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//end of the stream
			break
		}
		if err != nil {
			log.Fatalf("err while recieve: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}
