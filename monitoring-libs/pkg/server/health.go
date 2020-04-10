package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	health_pb "google.golang.org/grpc/health/grpc_health_v1"
)

func attachHealthServer(grpcServer *grpc.Server) {
	health := health.NewServer()
	health_pb.RegisterHealthServer(grpcServer, health)
}
