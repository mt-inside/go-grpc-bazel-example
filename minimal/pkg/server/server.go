package server

import (
	"context"
	"log"
	"net"

	pb "github.com/mt-inside/go-grpc-bazel-example/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Used to implement helloworld.GreeterServer
type Server struct {
	pb.UnimplementedGreeterServer // defines "unimplemented" methods for all RPCs so that this code is forwards-compatible

	port string
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s Server) Listen() {
	sock, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Fatalf("failed to lsiten: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterGreeterServer(srv, s)
	reflection.Register(srv)

	log.Printf("Listening on %v", s.port)
	if err := srv.Serve(sock); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s Server) SayHello(ctxt context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in)
	return &pb.HelloReply{Message: generateReply(in.GetName())}, nil
}

func generateReply(name string) string {
	return "Hello " + name
}
