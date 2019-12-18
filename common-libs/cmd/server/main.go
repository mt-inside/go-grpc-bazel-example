package main

import (
	"github.com/mt-inside/go-grpc-bazel-example/pkg/common"
	"github.com/mt-inside/go-grpc-bazel-example/pkg/server"
	"go.uber.org/config"
	"go.uber.org/zap"
	"os"
)

type cfg struct {
	Port string
}

func main() {
	log := common.NewLogger()

	var cfgPath string
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	} else {
		cfgPath = "config.yaml"
	}

	cfgF, err := os.Open(cfgPath)
	if err != nil {
		log.Fatalf("cannot open config file: %v", err)
	}
	defer cfgF.Close()

	provider, err := config.NewYAML(config.Source(cfgF))
	if err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	var c cfg
	if err := provider.Get("").Populate(&c);err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
	}

	log = log.With(zap.Namespace("server"), zap.String("port", c.Port))

	s := server.NewServer(log, c.Port)
	s.Listen()
}
