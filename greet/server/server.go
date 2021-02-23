package main

import (
	"context"
	"fmt"
	"grpc-course/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("Received a GreetRequest from : ", req.Greeting)
	fname := req.GetGreeting().GetFirstName()
	result := "Hello, " + fname
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, resStream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("Received a GreetManyTimesRequest from: ", req.Greeting)
	fname := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello, " + fname + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		resStream.Send(res)
		time.Sleep(time.Millisecond * 1000)
	}
	return nil
}

func (s *server) LongGreet(reqStream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet instance invoked")

	result := ""

	for {
		req, err := reqStream.Recv()
		if err == io.EOF {
			return reqStream.SendAndClose(&greetpb.LongGreetResponse{Result: result})
		} else if err != nil {
			log.Fatalf("Error receiving LongGreet request stream: %v", err)
		}

		result += "\nHello, " + req.GetGreeting().GetFirstName()
	}
}

func (s *server) GreetEveryone(reqStream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone instance invoked")

	for {
		req, err := reqStream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Fatalf("Error while receiveing message: %v", err)
			return err
		}

		fmt.Println(req.GetGreeting())
		result := "Hello, " + req.GetGreeting().GetFirstName()
		if err := reqStream.Send(&greetpb.GreetEveryoneResponse{Result: result}); err != nil {
			log.Fatalf("Error while receiveing message: %v", err)
		}
	}
}

func main() {
	fmt.Println("Hello World, I am a Server!!!")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf(" Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
