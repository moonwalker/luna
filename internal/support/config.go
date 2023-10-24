package support

import (
	"os/exec"
)

type build struct {}

type service struct {
	Dir   string
	Run   string
	Dep   []string
	Watch bool

	Build *build

	name string
	cmd  *exec.Cmd

	Changed bool
}

type Config struct {
	BuildTags []string
	Services map[string]*service
}
