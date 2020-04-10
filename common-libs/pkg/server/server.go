package server

import (
	"context"
	"net"

	pb "github.com/mt-inside/go-grpc-bazel-example/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerConfig struct {
	Port string
}

// Used to implement helloworld.GreeterServer
type Server struct {
	pb.UnimplementedGreeterServer // defines "unimplemented" methods for all RPCs so that this code is forwards-compatible

	log  *zap.SugaredLogger
	port string
}

func NewServer(log *zap.SugaredLogger, config *ServerConfig) *Server {
	log.Debugf("NewServer")

	return &Server{
		log:  log.With(zap.Namespace("server"), zap.String("port", config.Port)),
		port: config.Port,
	}
}

func (s Server) Listen() {
	sock, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		s.log.Fatalf("failed to listen: %v", err)
	}

	// TODO: split class into GreeterServer and GrpcServer, then fix up New*
	srv := grpc.NewServer()
	pb.RegisterGreeterServer(srv, s)
	reflection.Register(srv)

	s.log.Infof("Listening...")
	if err := srv.Serve(sock); err != nil {
		s.log.Fatalf("failed to serve: %v", err)
	}
}

func (s Server) SayHello(ctxt context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	s.log.Infof("Received: %v", in)
	return &pb.HelloReply{Message: generateReply(in.GetName())}, nil
}

func generateReply(name string) string {
	return "Hello " + name
}
