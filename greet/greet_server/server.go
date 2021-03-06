package main

import (
	"context"
	"fmt"
	"hello-grpc/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

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

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("Greet Many times function invoke with %v\n", req)
	firstname := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstname + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		err := stream.Send(res)
		if err != nil {
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreeting(stream greetpb.GreetService_LongGreetingServer) error {
	fmt.Printf("LongGreeting with stream\n")
	var str string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//send some thing back to client, if eof
			err = stream.SendAndClose(&greetpb.LongGreetingResponse{
				Result: str,
			})
			if err != nil {
				log.Fatalf("err while send response back: %v", err)
				return err
			}
			break
		}
		if err != nil {
			log.Fatalf("err while recv message: %v", err)
			return err
		}
		fmt.Println(req.GetGreeting().GetFirstName())
		str += "Hello " + req.GetGreeting().GetFirstName() + "! "
	}
	return nil
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone with BiDi\n")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("err while recv: %v", err)
			return err
		}
		firstname := req.GetGreeting().GetFirstName()
		result := "Hello " + firstname + "! "
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("err while sending data to Client: %v", err)
			return err
		}
	}
}
