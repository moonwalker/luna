package support

import (
	"os/exec"
)

type service struct {
	Chdir string
	Build string
	Start string
	Clean string
	Watch bool

	name string
	cmd  *exec.Cmd
}

type Config struct {
	Services map[string]*service
}
