package client

import (
	"context"
	"go.uber.org/zap"
	"time"

	pb "github.com/mt-inside/go-grpc-bazel-example/api"
	"google.golang.org/grpc"
)

type Client struct {
	log *zap.SugaredLogger
	address string
	conn    *grpc.ClientConn
	client  pb.GreeterClient
}

func NewClient(log *zap.SugaredLogger, address string) *Client {
	log.Debugf("Connecting to to %v...", address)

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	client := pb.NewGreeterClient(conn)

	log.Debugf("Connected to %v", address)

	return &Client{log, address, conn, client}
}

func (c Client) GetGreeting(name string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // will expire after 1 second. We pass this to the gRPC library, so it will give up if the request takes more than that.
	defer cancel()                                                        // manually cancel the context (and this gRPC lib's operation) if we panic or otherwise return from this function

	r, err := c.client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		c.log.Fatalf("could not get greeting: %v", err)
	}

	c.log.Debugf("Made SayHello call for name: %v", name)

	return r.GetMessage()
}

func (c Client) Close() {
	c.conn.Close()
}