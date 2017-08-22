package main

import (
	"github.com/moonwalker/luna/commands"
)

var (
	version = "dev"
	commit  = "HEAD"
	date    = "n/a"
)

func main() {
	commands.Execute(version, commit)
}
