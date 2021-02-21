package main

import (
	"context"
	"fmt"
	"grpc-course/greet/greetpb"
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

	doUnary(c)
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
