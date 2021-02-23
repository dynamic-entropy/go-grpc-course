package main

import (
	"context"
	"fmt"
	"grpc-course/calculate/calculatepb"
	"log"
	"math"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *server) SquareRoot(ctx context.Context, req *calculatepb.SquareRootRequest) (*calculatepb.SquareRootResponse, error) {
	fmt.Println("SquareRoot request received")

	num := req.GetNum()
	if num < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Input to Real Root Functions cannot be smaller than 0",
		)
	}
	res := &calculatepb.SquareRootResponse{
		RootNum: math.Sqrt(num),
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
