package main

import (
	"context"
	"fmt"
	"grpc-course/calculate/calculatepb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to create a client connection: %v", err)
	}
	defer cc.Close()

	c := calculatepb.NewCalculatorServiceClient(cc)

	// addNumbers(c, 5, 9)

	squareRoot(c, 10)  //Valid Argument call
	squareRoot(c, -10) //Invalid Argument call

}

func addNumbers(c calculatepb.CalculatorServiceClient, num1 int32, num2 int32) {

	in := &calculatepb.SumRequest{
		Num1: num1,
		Num2: num2,
	}

	sum, err := c.Sum(context.Background(), in)
	if err != nil {
		log.Fatalf("Failed to get a response from server: %v", err)
	}

	fmt.Printf("Sum of %d & %d is %d ", num1, num2, sum.GetSum())
}

func squareRoot(c calculatepb.CalculatorServiceClient, num float64) {

	in := &calculatepb.SquareRootRequest{
		Num: num,
	}

	res, err := c.SquareRoot(context.Background(), in)
	if err != nil {
		err, ok := status.FromError(err)
		if ok {
			fmt.Printf("Server:: Message: %v, Error Code: %v", err.Message(), err.Code())
			return
		}
		log.Fatalf("Unknown Fatal Error: %v", err)
		return
	}

	fmt.Printf("Response from Server :: Root of %f = %f \n", num, res.GetRootNum())

}
