package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/mt-inside/go-grpc-bazel-example/pkg/common"
	"github.com/mt-inside/go-grpc-bazel-example/pkg/client"
	"github.com/jawher/mow.cli"
)

const (
	defaultName = "world"
)

func main() {
	app := cli.App("client", "Get Greeted.")
	app.Spec = "ADDRESS [NAME]"
	app.Version("v version", fmt.Sprintf("client %v / %v", common.Version, runtime.Version()))

	var (
		// arg names must be supplied all-uppercase
		address = app.StringArg("ADDRESS", "", "Address of the Greeter server; host:port")
		name = app.StringArg("NAME", "world", "Name to Greet")
	)

	log := common.NewLogger()

	app.Action = func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second) // will expire after 1 second. We pass this to the gRPC library, so it will give up if the request takes more than that.
		defer cancel()                                                        // manually cancel the context (and this gRPC lib's operation) if we panic or otherwise return from this function

		c := client.NewClient(ctx, log, *address)
		defer c.Close()

		fmt.Println(fmt.Sprintf("Greeting: %s", c.GetGreeting(*name)))
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("could not start app: %v", err)
	}
}
