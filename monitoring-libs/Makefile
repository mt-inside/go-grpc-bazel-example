.PHONY: gazelle-force build run test
.DEFAULT_GOAL := build

%:
	@true

gazelle-force:
	bazel run //:gazelle -- update-repos -from_file=go.mod

build:
	bazel run //:gazelle
	bazel build //...

test:
	bazel test --test_output=errors //...

run-server: build
	bazel run //cmd/server "cmd/server/config/config.yaml"
run-client: build
	bazel run //cmd/client -- $(filter-out $@, $(MAKECMDGOALS))
