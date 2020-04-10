package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func attachHealthServer(grpcServer *grpc.Server) {
	health := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, health)
}
