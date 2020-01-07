package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/mt-inside/go-grpc-bazel-example/pkg/common"
	"github.com/mt-inside/go-grpc-bazel-example/pkg/server"
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
		common.NewCommonModule(),
		fx.Provide(
			NewConfig,
			server.NewServer,
		),
		fx.Invoke(func(s *server.Server) {
			go s.Listen()
		}),
		fx.Invoke(func(log *zap.SugaredLogger, config *server.ServerConfig) {
			port := config.PromPort
			log = log.With(zap.Namespace("prom"), zap.String("port", port))
			http.Handle("/metrics", promhttp.Handler())
			log.Infof("Listening for Prometheus scrapes")
			go http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
		}),
	)

	app.Run()
}
