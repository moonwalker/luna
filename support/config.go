package support

import (
	"os/exec"
	"strings"
)

type service struct {
	Chdir string
	Build string
	Start string
	Clean string
	Pack  string
	Watch bool

	name string
	cmd  *exec.Cmd

	Changed bool
}

type Config struct {
	Services map[string]*service
}

func (c *Config) PopulateChanges(compareRange string) {
	gitDiff := getGitDiff(compareRange)
	for _, diff := range gitDiff {
		for _, svc := range c.Services {
			if strings.Contains(diff, svc.Chdir) {
				svc.Changed = true
			}
		}
	}
}
