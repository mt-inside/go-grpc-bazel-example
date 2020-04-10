package server

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/mt-inside/go-grpc-bazel-example/api"
)

// Used to implement helloworld.GreeterServer
type greeterServer struct {
	pb.UnimplementedGreeterServer // defines "unimplemented" methods for all RPCs so that this code is forwards-compatible

	log      *zap.SugaredLogger
	nameLens prometheus.Histogram
}

func attachGreeterServer(log *zap.SugaredLogger, grpcServer *grpc.Server) {
	greeterServer := newGreeterServer(log)
	pb.RegisterGreeterServer(grpcServer, greeterServer)
}

func newGreeterServer(log *zap.SugaredLogger) *greeterServer {
	log.Debug("newGreeterServer")

	return &greeterServer{
		log: log.With(zap.Namespace("greeter server")),
		nameLens: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "name_lengths",
			Help:    "Lenghts of the names that have asked to be greeted",
			Buckets: prometheus.LinearBuckets(0, 1, 10),
		}),
	}
}

// Needs to be public, even though the type doesn't
// Also, not a compile time error if this isn't present; must be found with reflection
func (s greeterServer) SayHello(ctxt context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	name := in.GetName()

	s.log.Infof("Received: %v", in)
	s.nameLens.Observe(float64(len(name)))
	return &pb.HelloReply{Message: generateReply(name)}, nil
}

func generateReply(name string) string {
	return "Hello " + name
}
