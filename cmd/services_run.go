package cmd

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/support"
)

var (
	detach bool

	svcRunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run services",

		Run: func(cmd *cobra.Command, args []string) {
			pm := support.NewPM(cfg, args)
			pm.Run()
		},
	}
)

func init() {
	servicesCmd.AddCommand(svcRunCmd)
}
