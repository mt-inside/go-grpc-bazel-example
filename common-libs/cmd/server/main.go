package main

import (
	"os"
	"runtime"

	"github.com/mt-inside/go-grpc-bazel-example/pkg/common"
	"github.com/mt-inside/go-grpc-bazel-example/pkg/server"
	"go.uber.org/config"
)

func main() {
	log := common.NewLogger()
	log.Debug(runtime.Version())

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

	c := new(server.ServerConfig)
	if err := provider.Get("").Populate(c); err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
	}

	s := server.NewServer(log, c)
	s.Listen()
}
