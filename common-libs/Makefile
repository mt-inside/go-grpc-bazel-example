.PHONY: gazelle-force build run test
.DEFAULT_GOAL := build

%:
	@true

build:
	bazel run //:gazelle -- update-repos -prune=true -from_file=go.mod -to_macro=go_repos.bzl%go_repositories
	bazel run //:gazelle
	bazel build //...

test:
	bazel test --test_output=errors //...

run-server: build
	bazel run //cmd/server "cmd/server/config/config.yaml"
run-client: build
	bazel run //cmd/client -- $(filter-out $@, $(MAKECMDGOALS))
