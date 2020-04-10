package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	httpServer *http.Server
}

func myDecider(_ context.Context, _ string, _ interface{}) bool {
	return true
}
func myLevels(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK:
		return zap.InfoLevel // This is acutally the default. We don't change it here becuase this only affects the "call finished" log, which has a code. The acutal body printing is hardwired to Info
	default:
		return grpc_zap.DefaultCodeToLevel(code)
	}
}

func makeGrpcServer(log *zap.SugaredLogger) *grpc.Server {
	log.Debug("makeGrpcServer")

	srv := grpc.NewServer(
		// TODO: make this configurable, call and payload logging especially
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				// TODO extra interceptors: opentracing, prom (register with same place as hist)
				// TODO: is this how we do validate?
				grpc_zap.UnaryServerInterceptor(log.Desugar(), grpc_zap.WithLevels(myLevels)),
				grpc_zap.PayloadUnaryServerInterceptor(log.Desugar(), myDecider),
			),
		),
	)
	reflection.Register(srv)
	return srv
}

func makeHttpServer(log *zap.SugaredLogger, config *ServerConfig) *http.Server {
	log.Debug("makeHttpServer")

	return &http.Server{
		Addr: fmt.Sprintf(":%s", config.PromPort),
	}
}

func attachPromMetrics(httpServer *http.Server) {
	// TODO: make own mux. A server tht doesn't spec a handler will call http.DefaultServeMux though
	http.Handle("/metrics", promhttp.Handler())
}

func NewServer(log *zap.SugaredLogger, config *ServerConfig) *Server {
	grpcServer := makeGrpcServer(log)
	attachGreeterServer(log, grpcServer)
	attachHealthServer(grpcServer)

	httpServer := makeHttpServer(log, config)
	attachPromMetrics(httpServer)

	return &Server{
		// If we want to split these loggers, http and grpc serving need wrapping in a type each; another level of abstraction
		log:        log.With(zap.Namespace("server"), zap.String("port", config.Port), zap.String("prom_port", config.PromPort)),
		config:     config,
		grpcServer: grpcServer,
		httpServer: httpServer,
	}
}

func (s Server) Listen() {
	ctx := context.Background()
	ctx, cancelfn := context.WithCancel(ctx)

	// TODO feels like this should be in main(), but that would be really difficult with the current plumbing. Split this class into two - http and grpc, New and Listen both from main()
	// based on https://gist.github.com/akhenakh/38dbfea70dc36964e23acc19777f3869
	sigs := make(chan os.Signal, 1) // 1 is the recommended buffer size, guess you don't wanna block os.Signal at all (I guess that in a runtime with one system thread, the libc signal handler preëmpts all go routines so this would deadlock)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigs)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		s.log.Info("Listening for Prometheus scrapes")
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		return nil
	})

	g.Go(func() error {
		sock, err := net.Listen("tcp", ":"+s.config.Port)
		if err != nil {
			s.log.Fatalf("failed to listen: %v", err)
		}

		s.log.Info("Listening...")
		return s.grpcServer.Serve(sock)
	})

	select {
	case sig := <-sigs:
		s.log.Debugf("Signal: %s", sig)
		break
	case <-ctx.Done(): // happens if one of the goroutines returns an error
		s.log.Debug("Context cancelled; hit error")
		break
	}

	s.log.Info("Shutting down")

	cancelfn()
	// TODO: what's the format of this string? TODO needs refactoring to get at the healthserver itself
	//healthServer.SetServingStatus("grpc.health.v1.helloworld", healthpb.HealthCheckResponse_NOT_SERVING)

	shutdownCtx, shutdownCancelfn := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancelfn()

	s.grpcServer.GracefulStop()
	s.httpServer.Shutdown(shutdownCtx)

	if err := g.Wait(); err != nil {
		s.log.Fatalf("Shutdown caused by error: %v", err) // prints backtrace here, which is a bit annoying
	}
}
