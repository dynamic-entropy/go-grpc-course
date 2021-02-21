package main

import (
	"context"
	"fmt"
	"grpc-course/calculate/calculatepb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	calculatepb.UnimplementedCalculatorServiceServer
}

func (s *server) Sum(ctx context.Context, req *calculatepb.SumRequest) (*calculatepb.SumResponse, error) {
	fmt.Println("Sum request received")
	sum := req.GetNum1() + req.GetNum2()
	res := &calculatepb.SumResponse{
		Sum: sum,
	}

	return res, nil
}

func main() {
	fmt.Println("Server listening on port 50051")
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Server failed to listen: %v", err)
	}

	s := grpc.NewServer()

	calculatepb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Server failed to serve: %v", err)
	}

}
