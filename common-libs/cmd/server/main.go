package main

import "github.com/mt-inside/go-grpc-bazel-example/pkg/server"

func main() {
	s := server.NewServer("50051")
	s.Listen()
}
