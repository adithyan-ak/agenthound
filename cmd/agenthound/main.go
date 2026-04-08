package main

import (
	"fmt"
	"os"

	"github.com/adithyan-ak/agenthound/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	cli.SetVersion(version, commit)
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
