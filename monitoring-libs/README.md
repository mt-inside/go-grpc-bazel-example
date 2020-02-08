Example combining
* [Go](https://golang.org/) (with go 1.11+ modules)
* [gRPC](https://grpc.io/) ([This example](https://github.com/grpc/grpc-go/tree/master/examples/helloworld))
* [Bazel](https://bazel.build/) ([These rules](https://github.com/bazelbuild/rules_go))
* [Gazelle](https://github.com/bazelbuild/bazel-gazelle)

And

* Cmd-line parsing with [mow.cli](https://github.com/jawher/mow.cli)
* Config-file parsing with [uber config](https://github.com/uber-go/config)
* Logging with [zap](https://github.com/uber-go/zap)
* Dependancy Injection with [fx](https://github.com/uber-go/fx)

And

* Metrics with [Prometheus](https://prometheus.io/) ([This example](https://godoc.org/github.com/prometheus/client_golang/prometheus/promauto))

Notes
Container image contains the gRPC health check probe, run with: `docker exec -ti <container id> /bin/grpc_health_probe-linux-amd64 -addr localhost:50051`
