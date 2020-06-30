.PHONY: build test run-server run-client local-image local-image-run local-stack-run
.DEFAULT_GOAL := build

LINUX := --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64

gazelle:
	bazel run //:gazelle -- update-repos -prune=true -from_file=go.mod -to_macro=go_repos.bzl%go_repositories
	bazel run //:gazelle

build-all: gazelle
	bazel build //...

test: gazelle
	bazel test --test_output=errors //...

run-server: gazelle
	bazel run //cmd/server "cmd/server/config/config.yaml"
run-client: gazelle
	bazel run //cmd/client -- localhost:50051 $(filter-out $@, $(MAKECMDGOALS))

image:
	bazel build $(LINUX) //cmd/server:image.tar

image-push:
	bazel run $(LINUX) //cmd/server:push

image-local-load: image
	docker load -i bazel-bin/cmd/server/image.tar

image-local-run: image-local-load
	docker run -p50051:50051 bazel/cmd/server:image

local-stack-run: local-image
	docker-compose -f deploy/local/docker-compose.yaml rm
	docker-compose -f deploy/local/docker-compose.yaml up
