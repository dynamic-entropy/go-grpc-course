package main

import (
	"context"
	"fmt"
	"grpc-course/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBidirectionalStreaming(c)

	doUnaryWithDeadline(c, 5)
	doUnaryWithDeadline(c, 1)
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
		log.Fatalf("Unable to create a GreetManyTimesClient: %v", err)
	}

	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error occured while receiving stream response: %v ", err)
		}
		fmt.Println("GreetManyTimes Message Received: ", res.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {

	reqStream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Unable to create a LongGreetClient: %v", err)
	}

	ins := []*greetpb.LongGreetRequest{
		{Greeting: &greetpb.Greeting{
			FirstName: "India",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "China",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Russia",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "United Kingdom",
		}},
	}

	for _, in := range ins {
		if err := reqStream.Send(in); err != nil {
			log.Fatalf("Error occured while sending requests: %v", err)
		}
		time.Sleep(time.Millisecond * 1000)
	}

	res, err := reqStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to recieve response: %v", err)
	}

	fmt.Println("Response from server: ", res.GetResult())

}

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	reqStream, err := c.GreetEveryone(context.TODO())
	if err != nil {
		log.Fatalf("Unable to create a GreetEveryone Client: %v", err)
	}

	ins := []*greetpb.GreetEveryoneRequest{
		{Greeting: &greetpb.Greeting{
			FirstName: "India",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "China",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "Russia",
		}},
		{Greeting: &greetpb.Greeting{
			FirstName: "United Kingdom",
		}},
	}

	wc := make(chan struct{})
	go func() {

		for _, in := range ins {
			if err := reqStream.Send(in); err != nil {
				log.Fatalf("Error occured while sending requests: %v", err)
			}
			fmt.Println("Sent Message: ", in)
			time.Sleep(time.Millisecond * 1000)
		}
		reqStream.CloseSend()
	}()

	go func() {
		for {
			res, err := reqStream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalf("Error while receiving response: %v", err)
				break
			}

			fmt.Println(res.GetResult())
		}
		close(wc)
	}()

	<-wc
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {
	in := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Rahul",
			LastName:  "Chauhan",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, in)
	if err != nil {
		refErr, ok := status.FromError(err)
		if ok {
			if refErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit: ", refErr.Code())
			} else {
				fmt.Println("Unknown Error: ", refErr.Code())
			}
		} else {
			log.Fatalf("Greeting Request Failed: %v", err)
		}
		return
	}
	fmt.Printf("Greeting Response: %s \n", res.Result)
}
