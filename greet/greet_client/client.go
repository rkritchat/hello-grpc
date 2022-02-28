package main

import (
	"context"
	"fmt"
	"hello-grpc/greet/greetpb"
	"io"
	"log"
	"time"

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

	//doServerStreaming(c)

	//doClientStreaming(c)

	doBiDirectional(c)
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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a client streaming RPC...")
	firstname := []string{"kritchat", "aaaaa", "bbbb", "cccc"}

	stream, err := c.LongGreeting(context.Background())
	if err != nil {
		log.Fatalf("err while invoke long greeting")
	}
	for _, val := range firstname {
		fmt.Printf("start sent %v to stream\n", val)
		err = stream.Send(&greetpb.LongGreetingRequest{
			Greeting: &greetpb.Greeting{
				FirstName: val,
			},
		})
		if err != nil {
			log.Fatalf("err while stream.send: %v", err)
		}
		time.Sleep(time.Second)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("err while close and recv: %v", err)
	}
	fmt.Printf("Response from LongGreeting: %v", resp)
}

func doBiDirectional(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a client streaming RPC...")
	firstname := []string{"kritchat", "any", "anyname", "john"}
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("err while init stream: %v", err)
		return
	}

	wait := make(chan struct{})
	//send data to server
	go func() {
		for _, val := range firstname {
			fmt.Printf("sending message %v\n", val)
			err = stream.Send(&greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: val,
				},
			})
			if err != nil {
				log.Fatalf("err while send data to Server: %v", err)
				return
			}
			time.Sleep(time.Duration(2) * time.Second)
		}
		err = stream.CloseSend()
		if err != nil {
			log.Fatalf("err while Close send: %v", err)
			return
		}
	}()

	//consume data from server
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(wait) //close the channel
				return
			}
			if err != nil {
				log.Fatalf("err while recv: %v", err)
				return
			}
			fmt.Printf("<--- %v\n", resp.GetResult())
		}
	}()
	<-wait
}
