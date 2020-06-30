package server

import (
	"go.uber.org/fx"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewHealthServerModule() fx.Option {
	return fx.Options(
		fx.Provide(health.NewServer), // Just defers to the built-in health server
		fx.Provide(func(hs *health.Server) grpc_health_v1.HealthServer { return hs }), // See GreeterServerModule
		fx.Invoke(grpc_health_v1.RegisterHealthServer),
	)
}
