package main

import (
	"context"
	"fmt"
	"grpc-course/greet/greetpb"
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
