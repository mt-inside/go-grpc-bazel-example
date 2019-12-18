package client

import (
	"context"
	"log"
	"time"

	pb "github.com/mt-inside/go-grpc-bazel-example/api"
	"google.golang.org/grpc"
)

type Client struct {
	address string
	conn    *grpc.ClientConn
	client  pb.GreeterClient
}

func NewClient(address string) *Client {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	client := pb.NewGreeterClient(conn)

	return &Client{address, conn, client}
}

func (c Client) GetGreeting(name string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // will expire after 1 second. We pass this to the gRPC library, so it will give up if the request takes more than that.
	defer cancel()                                                        // manually cancel the context (and this gRPC lib's operation) if we panic or otherwise return from this function

	r, err := c.client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not get greeting: %v", err)
	}

	return r.GetMessage()
}

func (c Client) Close() {
	c.conn.Close()
}

// write up: optional args:
// - nils
// - variadic (if they're all the same type)
// - params struct (if you don't specify some when making the struct they get the zero value. If that's not OK, have a new() func for the struct that supplies other defaults)
// - functional options
