package main

import (
	"context"
	"fmt"
	"grpc-course/greet/greetpb"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	log.Print("Hello World, I am a Client!!")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to open a connection: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	// doUnary(c)
	doServerStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	in := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rahul",
			LastName:  "Chauhan",
		},
	}
	res, err := c.Greet(context.Background(), in)
	if err != nil {
		log.Fatalf("Greeting Request Failed: %v", err)
	}
	fmt.Printf("Greeting Response: %s", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	in := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rahul",
			LastName:  "Chauhan",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), in)
	if err != nil {
		log.Fatalf("Error occured while calling server: %v", err)
	}

	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error occured while server streaming: %v ", err)
		}
		fmt.Println("GreetManyTimes Message Received: ", res.GetResult())
	}
}
