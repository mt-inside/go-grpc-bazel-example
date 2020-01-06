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

local-image:
	bazel build --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64 //cmd/server:server_image.tar
	docker load -i bazel-bin/cmd/server/server_image.tar

local-image-run: local-image
	docker run bazel/cmd/server:server_image
