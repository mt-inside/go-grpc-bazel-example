.PHONY: build test run-server run-client
.DEFAULT_GOAL := build

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
	bazel run //cmd/client -- localhost:50001 $(filter-out $@, $(MAKECMDGOALS))
