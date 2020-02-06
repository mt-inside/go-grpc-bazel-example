package server

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	pb "github.com/mt-inside/go-grpc-bazel-example/api"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Used to implement helloworld.GreeterServer
type GreeterServer struct {
	pb.UnimplementedGreeterServer // defines "unimplemented" methods for all RPCs so that this code is forwards-compatible

	log       *zap.SugaredLogger
	name_lens prometheus.Histogram
}

func NewGreeterServer(log *zap.SugaredLogger, config *ServerConfig) *GreeterServer {
	log.Debugf("NewGreeterServer")

	return &GreeterServer{
		log: log.With(zap.Namespace("greeter server")),
		name_lens: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "name_lengths",
			Help:    "Lenghts of the names that have asked to be greeted",
			Buckets: prometheus.LinearBuckets(0, 1, 10),
		}),
	}
}

func (s GreeterServer) SayHello(ctxt context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	name := in.GetName()

	s.log.Infof("Received: %v", in)
	s.name_lens.Observe(float64(len(name)))
	return &pb.HelloReply{Message: generateReply(name)}, nil
}

func generateReply(name string) string {
	return "Hello " + name
}

func NewGreeterServerModule() fx.Option {
	return fx.Options(
		fx.Provide(NewGreeterServer),
		/* Problem is that we want to use a struct as an iface, and because of duck-typing we have to be explicit about the tie-up.
		* Is there a built-in way to do this? ie to provide a ctor for a named iface?
		* If not, it could be a bit neater to pass a func to Invoke which takes *grpc.Server and the struct type, and just call Register*Server() on them */
		fx.Provide(func(gs *GreeterServer) pb.GreeterServer { return gs }),
		fx.Invoke(pb.RegisterGreeterServer),
	)
}
