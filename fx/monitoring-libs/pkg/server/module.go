package server

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ServerConfig struct {
	Port     string
	PromPort string
}

func NewGrpcServer(log *zap.SugaredLogger) (*grpc.Server, error) {
	log.Debugf("NewGrpcServer")

	srv := grpc.NewServer(
		// TODO: make this configurable
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(log.Desugar()),
			),
		),
	)
	reflection.Register(srv)
	return srv, nil
}

func NewServerModule() fx.Option {
	return fx.Options(
		fx.Provide(NewGrpcServer),

		NewHealthServerModule(),
		NewGreeterServerModule(),

		fx.Invoke(func(lifecycle fx.Lifecycle, srv *grpc.Server, log *zap.SugaredLogger, cfg *ServerConfig) {
			log = log.With(zap.Namespace("grpc server"), zap.String("port", cfg.Port))
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					log.Debugf("OnStart")

					sock, err := net.Listen("tcp", ":"+cfg.Port)
					if err != nil {
						log.Fatalf("failed to listen: %v", err)
					}

					go func() {
						log.Infof("Listening...")
						if err := srv.Serve(sock); err != nil {
							log.Fatalf("failed to serve: %v", err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Debugf("OnStop")

					srv.Stop()
					return nil
				},
			})
		}),
	)
}
