package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerConfig struct {
	Port     string
	PromPort string
}

type Server struct {
	log        *zap.SugaredLogger
	config     *ServerConfig
	grpcServer *grpc.Server
}

func myDecider(_ context.Context, _ string, _ interface{}) bool {
	return true
}

func makeGrpcServer(log *zap.SugaredLogger) *grpc.Server {
	log.Debugf("makeGrpcServer")

	srv := grpc.NewServer(
		// TODO: make this configurable, call and payload logging especially
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// TODO extra interceptors: opentracing, prom (register with same place as hist)
				// TODO: is this how we do validate?
				grpc_zap.UnaryServerInterceptor(log.Desugar()), // always logs at Info, nothing you can do about it
				grpc_zap.PayloadUnaryServerInterceptor(log.Desugar(), myDecider),
			),
		),
	)
	reflection.Register(srv)
	return srv
}

/* Very basic for now, just uses the implicit mux etc
* TODO make serveHttp, httpSrv.attachProm, httpSrv.listen - using explicit http mux and server
 */
func serveProm(log *zap.SugaredLogger, c *ServerConfig) {
	port := c.PromPort
	log = log.With(zap.Namespace("prom"), zap.String("port", port))
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Listening for Prometheus scrapes")
	go http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func NewServer(log *zap.SugaredLogger, config *ServerConfig) *Server {
	serveProm(log, config)

	grpcServer := makeGrpcServer(log)
	attachGreeterServer(log, grpcServer)
	attachHealthServer(grpcServer)

	return &Server{
		log:        log.With(zap.Namespace("server"), zap.String("port", config.Port)),
		config:     config,
		grpcServer: grpcServer,
	}
}

func (s Server) Listen() {
	sock, err := net.Listen("tcp", ":"+s.config.Port)
	if err != nil {
		s.log.Fatalf("failed to listen: %v", err)
	}
	s.log.Info("Listening...")
	go func() {
		if err := s.grpcServer.Serve(sock); err != nil {
			s.log.Fatalf("failed to serve: %v", err)
		}
	}()

	// TODO move to main, impliment https://gist.github.com/akhenakh/38dbfea70dc36964e23acc19777f3869
	sigs := make(chan os.Signal, 1) // 1 is the recommended buffer size, guess you don't wanna block os.Signal at all (I guess that in a runtime with one system thread, the libc signal handler preëmpts all go routines so this would deadlock)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigs
	s.log.Debugf("Signal: %s", sig)

	s.log.Info("Shutting down")
	s.grpcServer.GracefulStop()
}
