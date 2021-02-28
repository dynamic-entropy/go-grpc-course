package main

import (
	"context"
	"fmt"
	"grpc-course/blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
)

var coll *mongo.Collection

type server struct {
	blogpb.UnimplementedBlogServiceServer
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id,omitempty"`
	Title    string             `bson:"title,omitempty"`
	Content  string             `bson:"content,omitempty"`
}

func (s *server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetTitle(),
	}

	result, err := coll.InsertOne(ctx, data)
	if err != nil {
		log.Fatalf("Unable to insert document: %v", err)
	}

	res := &blogpb.CreateBlogResponse{
		Id: result.InsertedID.(primitive.ObjectID).Hex(),
	}

	return res, nil
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile) //get File name and line number on failure

	fmt.Println("--->Connecting to MongoDB Server")
	defer fmt.Println("--->Server Stopped")
	client, errMongo := mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if errMongo != nil {
		log.Fatalf("Failed to create new mongo client: %v", errMongo)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to initialise the client and start background monitoring: %v", err)
	}
	defer client.Disconnect(ctx)
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	db := client.Database("web_db")
	coll = db.Collection("blogs")

	// r, e := coll.InsertOne(ctx, bson.D{
	// 	{Key: "title", Value: "The White Tiger"},
	// 	{Key: "author_id", Value: "Aravind Adiga"},
	// 	{Key: "content", Value: "An excellent book"},
	// })
	// if e != nil {
	// 	log.Fatal(e)
	// }
	// fmt.Println(r)

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create a listener: %v ", err)
	}
	defer lis.Close()

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	defer s.Stop()
	blogpb.RegisterBlogServiceServer(s, &server{})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		fmt.Println("--->Server Running ")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ch

	fmt.Println(" ::User Interrupt: Stopping Server")
}
