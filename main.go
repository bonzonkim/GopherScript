package main

import (
	"fmt"
	"os"

	"github.com/bonzonkim/gopher-script/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
