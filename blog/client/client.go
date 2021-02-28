package main

import (
	"context"
	"fmt"
	"grpc-course/blog/blogpb"
	"log"

	"google.golang.org/grpc"
)

func main() {

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	cc, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatalf("Unable to create a client connection: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	createCall(c)
}

func createCall(c blogpb.BlogServiceClient) {

	blog := &blogpb.Blog{
		AuthorId: "Khaleed",
		Content:  "Story of an unlikely friendship",
		Title:    "The Kite Runner",
	}
	req := &blogpb.CreateBlogRequest{
		Blog: blog,
	}

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("Server returned following error: %v", err)
	}

	fmt.Printf("Id of new blog: %s \n", res.GetId())
}
