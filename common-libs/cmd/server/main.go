package main

import (
	"github.com/mt-inside/go-grpc-bazel-example/pkg/common"
	"github.com/mt-inside/go-grpc-bazel-example/pkg/server"
	"go.uber.org/zap"
)

const (
	port string = "50051"
)

func main() {
	log := common.NewLogger().With(zap.Namespace("server"), zap.String("port", port))

	s := server.NewServer(log, port)
	s.Listen()
}
