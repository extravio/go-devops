package main

import (
	"fmt"
	"os"
)

// import "rsc.io/quote"

func main() {
	// var args []string
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Usage: hello-devops <argument>\n")
		os.Exit(1)
	}

	fmt.Printf("Hello world\nArguments: %v\n", args[1:])
}
