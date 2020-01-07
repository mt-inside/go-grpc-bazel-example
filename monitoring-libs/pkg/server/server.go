package server

import (
	"context"
	"net"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	pb "github.com/mt-inside/go-grpc-bazel-example/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerConfig struct {
	Port     string
	PromPort string
}

// Used to implement helloworld.GreeterServer
type Server struct {
	pb.UnimplementedGreeterServer // defines "unimplemented" methods for all RPCs so that this code is forwards-compatible

	log       *zap.SugaredLogger
	port      string
	name_lens prometheus.Histogram
}

func NewServer(log *zap.SugaredLogger, config *ServerConfig) *Server {
	log.Debugf("NewServer")

	return &Server{
		log:  log.With(zap.Namespace("server"), zap.String("port", config.Port)),
		port: config.Port,
		name_lens: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "name_lengths",
			Help:    "Lenghts of the names that have asked to be greeted",
			Buckets: prometheus.LinearBuckets(0, 1, 10),
		}),
	}
}

func (s Server) Listen() {
	sock, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		s.log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterGreeterServer(srv, s)
	reflection.Register(srv)

	s.log.Infof("Listening...")
	if err := srv.Serve(sock); err != nil {
		s.log.Fatalf("failed to serve: %v", err)
	}
}

func (s Server) SayHello(ctxt context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	name := in.GetName()

	s.log.Infof("Received: %v", in)
	s.name_lens.Observe(float64(len(name)))
	return &pb.HelloReply{Message: generateReply(name)}, nil
}

func generateReply(name string) string {
	return "Hello " + name
}
