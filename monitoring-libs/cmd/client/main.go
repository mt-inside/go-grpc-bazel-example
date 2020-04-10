/* Tbh I think fx is unnecessary here and makes this overly complicated */

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	cli "github.com/jawher/mow.cli" // this alias isn't necessary, but `go fmt` removes this import if it's not present

	"github.com/mt-inside/go-grpc-bazel-example/pkg/client"
	"github.com/mt-inside/go-grpc-bazel-example/pkg/common"
)

func main() {
	log := common.NewLogger()
	log.Debug(runtime.Version())

	app := cli.App("client", "Get Greeted.")
	app.Spec = "ADDRESS [NAME]"
	app.Version("v version", fmt.Sprintf("client %v / %v", common.Version, runtime.Version()))

	// arg names must be supplied all-uppercase
	address := app.StringArg("ADDRESS", "", "Address of the Greeter server; host:port")
	name := app.StringArg("NAME", "world", "Name to Greet")

	app.Action = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) // will expire after 2 seconds. We pass this to the gRPC library, so it will give up if the request takes more than that.
		defer cancel()                                                          // manually cancel the context (and this gRPC lib's operation) if we panic or otherwise return from this function

		c := client.NewClient(ctx, log, *address)
		defer c.Close()

		fmt.Println(fmt.Sprintf("Greeting: %s", c.GetGreeting(*name)))
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("could not start app: %v", err)
	}
}
