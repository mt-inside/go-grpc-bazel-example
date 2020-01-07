package main

import (
	"os"

	"github.com/mt-inside/go-grpc-bazel-example/pkg/common"
	"github.com/mt-inside/go-grpc-bazel-example/pkg/server"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewConfig(log *zap.SugaredLogger) *server.ServerConfig {
	log.Debugf("NewConfig")

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

	var c server.ServerConfig
	if err := provider.Get("").Populate(&c); err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
	}

	return &c
}

func main() {
	app := fx.New(
		fx.Provide(
			NewConfig,
			common.NewLogger,
			server.NewServer,
		),
		fx.Invoke(func(s *server.Server) {
			s.Listen()
		}),
	)

	app.Run()
}
