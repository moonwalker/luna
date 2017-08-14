package commands

import (
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/support"
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

type config struct {
	Services map[string]*service
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run services specified in config",

	Run: func(cmd *cobra.Command, args []string) {
		pm := support.NewPM()
		pm.Run()
	},
}

func init() {
	servicesCmd.AddCommand(runCmd)
}
