package commands

import (
	"github.com/spf13/cobra"

	"github.com/moonwalker/luna/support"
)

var svcRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run services",

	Run: func(cmd *cobra.Command, args []string) {
		pm := support.NewPM(cfg)
		pm.Run()
	},
}

func init() {
	servicesCmd.AddCommand(svcRunCmd)
}
