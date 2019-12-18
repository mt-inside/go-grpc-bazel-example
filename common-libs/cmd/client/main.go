package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

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

	app.Action = func() {
		c := client.NewClient(*address)
		defer c.Close()

		log.Printf("Greeting: %s", c.GetGreeting(*name))
	}

	app.Run(os.Args)
}
