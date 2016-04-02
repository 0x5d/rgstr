package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/castillobg/rgstr/registries"
	// A blank import so that the consul AdapterFactory registers itself.
	_ "github.com/castillobg/rgstr/registries/consul"
	"github.com/castillobg/rgstr/runtimes"
	// A blank import so that the rkt AdapterFactory registers itself.
	_ "github.com/castillobg/rgstr/runtimes/rkt"
	"github.com/coreos/rkt/api/v1alpha"
)

var delay *int
var pods = make(map[string]*v1alpha.Pod)

func main() {
	address := flag.String("a", "localhost:15441", "The `address` where rkt's API service is listening.")
	consulAddress := flag.String("ra", "localhost:8500", "The `registry address`.")
	flag.Parse()

	registryName := "consul"
	registryFactory, ok := registries.LookUp(registryName)
	if !ok {
		fmt.Printf("No registry with name \"%s\" found.\n", registryName)
		os.Exit(1)
	}
	registry, err := registryFactory.New(*consulAddress)
	if err != nil {
		fmt.Printf("Error initializing registry client for \"%s\": %s\n", registryName, err.Error())
		os.Exit(1)
	}

	runtimeName := "rkt"
	runtimeFactory, ok := runtimes.LookUp(runtimeName)
	if !ok {
		fmt.Printf("No runtime with name \"%s\" found.\n", registryName)
		os.Exit(1)
	}
	runtime, err := runtimeFactory.New(*address, &registry)
	if err != nil {
		fmt.Printf("Error initializing runtime client for \"%s\": %s", runtime, err.Error())
		os.Exit(1)
	}

	errs := make(chan error)
	go runtime.Listen(errs)
	fmt.Printf("rgstr is listening for changes in %s...\n", runtimeName)
	for err = range errs {
		fmt.Println(err)
	}
}
