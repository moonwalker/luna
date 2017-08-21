package main

import (
	"github.com/moonwalker/luna/commands"
)

var (
	GitCommit = "HEAD"
)

func main() {
	commands.Execute(GitCommit)
}
