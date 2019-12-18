package client

import (
	"context"
	"go.uber.org/zap"

	pb "github.com/mt-inside/go-grpc-bazel-example/api"
	"google.golang.org/grpc"
)

type Client struct {
	ctx context.Context
	log *zap.SugaredLogger
	address string
	conn    *grpc.ClientConn
	client  pb.GreeterClient
}

func NewClient(ctx context.Context, log *zap.SugaredLogger, address string) *Client {
	log.Debugf("Connecting to to %v...", address)

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	client := pb.NewGreeterClient(conn)

	log.Debugf("Connected to %v", address)

	return &Client{ctx, log, address, conn, client}
}

func (c Client) GetGreeting(name string) string {
	r, err := c.client.SayHello(c.ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		c.log.Fatalf("could not get greeting: %v", err)
	}

	c.log.Debugf("Made SayHello call for name: %v", name)

	return r.GetMessage()
}

func (c Client) Close() {
	c.conn.Close()
}