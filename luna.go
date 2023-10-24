package main

import (
	"fmt"
	"os"

	"github.com/moonwalker/luna/internal/cli"
)

var (
	version = "dev"
	commit  = "HEAD"
	date    = "n/a"
)

func main() {
	if err := cli.Run(version, commit, date); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
