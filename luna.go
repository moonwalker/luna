package main

import (
	"fmt"
	"os"

	"github.com/moonwalker/luna/internal/cli"
)

var (
	version = "dev"
	commit  = "HEAD"
)

func main() {
	if err := cli.Run(version, commit); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
