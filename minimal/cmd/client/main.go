package main

import (
	"log"
	"os"

	"github.com/mt-inside/go-grpc-bazel-example/pkg/client"
)

const (
	defaultName = "world"
)

func main() {
	var address string
	var name string

	switch len(os.Args) {
	case 1:
		log.Fatalf("Usage: %v address [name]", os.Args[0])
	case 2:
		address = os.Args[1]
		name = defaultName
	case 3:
		address = os.Args[1]
		name = os.Args[2]
	default:
		log.Fatalf("Usage: %v address [name]", os.Args[0])
	}

	c := client.NewClient(address)
	defer c.Close()

	log.Printf("Greeting: %s", c.GetGreeting(name))
}
